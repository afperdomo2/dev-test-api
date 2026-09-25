package ai

import (
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestEvaluateAiQuota(t *testing.T) {
	tests := []struct {
		name        string
		isAdmin     bool
		limit       int
		used        int64
		wantErr     bool
		wantLimitIn int
	}{
		{
			name:    "admin exento con cero uso",
			isAdmin: true,
			limit:   20,
			used:    0,
			wantErr: false,
		},
		{
			name:    "admin exento aunque exceda el cupo",
			isAdmin: true,
			limit:   20,
			used:    999,
			wantErr: false,
		},
		{
			name:    "admin exento con limite cero",
			isAdmin: true,
			limit:   0,
			used:    0,
			wantErr: false,
		},
		{
			name:    "usuario normal por debajo del limite",
			isAdmin: false,
			limit:   20,
			used:    19,
			wantErr: false,
		},
		{
			name:    "usuario normal en el limite exacto",
			isAdmin: false,
			limit:   20,
			used:    20,
			wantErr: true,
		},
		{
			name:    "usuario normal excedido",
			isAdmin: false,
			limit:   20,
			used:    25,
			wantErr: true,
		},
		{
			name:    "usuario normal sin uso",
			isAdmin: false,
			limit:   20,
			used:    0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := evaluateAiQuota(tt.isAdmin, tt.limit, tt.used)
			if !tt.wantErr {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			var limitErr *DailyLimitError
			require.ErrorAs(t, err, &limitErr)
			assert.Equal(t, tt.limit, limitErr.Limit)
		})
	}
}

// TestCheckDailyAiLimit_FailClosedOnDBError verifica que un error de base de datos
// NO se interprete como "sin límite": previously devolvía nil y permitía generar IA
// sin control cuando la consulta del usuario fallaba.
func TestCheckDailyAiLimit_FailClosedOnDBError(t *testing.T) {
	sqlDB, err := sql.Open("pgx", "postgres://user:pass@127.0.0.1:1/db")
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())

	db, err := gorm.Open(
		postgres.New(postgres.Config{Conn: sqlDB, PreferSimpleProtocol: true}),
		&gorm.Config{DisableAutomaticPing: true},
	)
	require.NoError(t, err)

	g := &Generator{db: db}

	err = g.checkDailyAiLimit(uuid.New())
	require.Error(t, err, "un error de BD debe propagarse, no asumirse como sin límite")
	assert.NotErrorIs(t, err, &DailyLimitError{}, "no debe confundirse con cupo agotado")
}
