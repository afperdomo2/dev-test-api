package generation

import "time"

type JobResponse struct {
	ID          string     `json:"id"`
	RequestedBy string     `json:"requestedBy"`
	TopicIDs    []string   `json:"topicIds"`
	Topics      []string   `json:"topics"`
	Difficulty  string     `json:"difficulty"`
	Quantity    int        `json:"quantity"`
	Status      string     `json:"status"`
	Generated   int        `json:"generated"`
	Error       string     `json:"error,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	StartedAt   *time.Time `json:"startedAt,omitempty"`
	FinishedAt  *time.Time `json:"finishedAt,omitempty"`
}

type StatsResponse struct {
	// GeneratedToday es lo realmente generado y guardado hoy.
	GeneratedToday int `json:"generatedToday"`
	// ReservedToday es lo comprometido hoy (pending + running + completed).
	// Es lo que se cobra contra el tope diario.
	ReservedToday   int  `json:"reservedToday"`
	DailyCap        int  `json:"dailyCap"`
	PerJobCap       int  `json:"perJobCap"`
	DefaultQuantity int  `json:"defaultQuantity"`
	Pending         int  `json:"pending"`
	Running         int  `json:"running"`
	Failed          int  `json:"failed"`
	AiConfigured    bool `json:"aiConfigured"`
}
