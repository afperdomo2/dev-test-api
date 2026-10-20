package push

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/felipe/dev-test-api/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func okTickets(count int) expoResponse {
	tickets := make([]expoTicket, count)
	for i := range tickets {
		tickets[i] = expoTicket{Status: "ok", Message: "ok"}
	}
	return expoResponse{Data: tickets}
}

func TestExpoSenderBatches(t *testing.T) {
	var mu sync.Mutex
	var batches [][]expoMessage

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var messages []expoMessage
		require.NoError(t, json.NewDecoder(r.Body).Decode(&messages))

		mu.Lock()
		batches = append(batches, messages)
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(okTickets(len(messages))))
	}))
	defer server.Close()

	sender := NewExpoSender(config.PushConfig{ExpoURL: server.URL}, newPushStore())
	messages := make([]expoMessage, 101)
	for i := range messages {
		messages[i] = expoMessage{To: "tok", Title: "t", Body: "b"}
	}

	require.NoError(t, sender.Send(context.Background(), messages))
	require.Len(t, batches, 2)
	assert.Len(t, batches[0], 100)
	assert.Len(t, batches[1], 1)
}

func TestExpoSenderPayloadAndAuth(t *testing.T) {
	var got []expoMessage
	var authorization string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization = r.Header.Get("Authorization")
		require.NoError(t, json.NewDecoder(r.Body).Decode(&got))
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(okTickets(len(got))))
	}))
	defer server.Close()

	sender := NewExpoSender(config.PushConfig{
		ExpoURL:     server.URL,
		AccessToken: "secret",
	}, newPushStore())

	require.NoError(t, sender.Send(context.Background(), []expoMessage{{
		To:    "ExponentPushToken[abc]",
		Title: "Repaso de hoy",
		Body:  "Tienes 3 preguntas para repasar",
		Data:  map[string]string{"notificationId": "n1", "type": "daily_review"},
	}}))

	require.Len(t, got, 1)
	assert.Equal(t, "ExponentPushToken[abc]", got[0].To)
	assert.Equal(t, "Repaso de hoy", got[0].Title)
	assert.Equal(t, "Tienes 3 preguntas para repasar", got[0].Body)
	assert.Equal(t, "n1", got[0].Data["notificationId"])
	assert.Equal(t, "daily_review", got[0].Data["type"])
	assert.Equal(t, "Bearer secret", authorization)
}

func TestExpoSenderDeletesDeviceNotRegistered(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var messages []expoMessage
		require.NoError(t, json.NewDecoder(r.Body).Decode(&messages))

		tickets := make([]expoTicket, len(messages))
		tickets[0] = expoTicket{
			Status:  "error",
			Message: "not registered",
			Details: expoTicketDetails{Error: "DeviceNotRegistered"},
		}
		for i := 1; i < len(tickets); i++ {
			tickets[i] = expoTicket{Status: "ok"}
		}

		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(expoResponse{Data: tickets}))
	}))
	defer server.Close()

	store := newPushStore()
	var deleted []string
	store.deleteByTokenFn = func(token string) error {
		deleted = append(deleted, token)
		return nil
	}

	sender := NewExpoSender(config.PushConfig{ExpoURL: server.URL}, store)
	require.NoError(t, sender.Send(context.Background(), []expoMessage{
		{To: "tok-muerto"},
		{To: "tok-vivo"},
	}))

	assert.Equal(t, []string{"tok-muerto"}, deleted)
}

func TestExpoSenderHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	sender := NewExpoSender(config.PushConfig{ExpoURL: server.URL}, newPushStore())
	err := sender.Send(context.Background(), []expoMessage{{To: "tok"}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "500")
}

func TestExpoSenderTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(300 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(okTickets(1))
	}))
	defer server.Close()

	sender := NewExpoSender(config.PushConfig{ExpoURL: server.URL}, newPushStore())
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	require.Error(t, sender.Send(ctx, []expoMessage{{To: "tok"}}))
}

func TestExpoSenderEmpty(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
	}))
	defer server.Close()

	sender := NewExpoSender(config.PushConfig{ExpoURL: server.URL}, newPushStore())
	require.NoError(t, sender.Send(context.Background(), nil))
	assert.Zero(t, calls)
}
