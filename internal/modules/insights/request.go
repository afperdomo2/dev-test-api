package insights

// SummaryParams es el query de `GET /insights/summary`.
// `period` es opcional; el service aplica `week` por defecto.
type SummaryParams struct {
	Period string `form:"period" binding:"omitempty,oneof=week month"`
}
