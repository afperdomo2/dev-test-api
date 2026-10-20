package gamification

// BadgeID identifica un logro. Los ids son estables: web y mobile los usan
// como clave, así que no se renombran.
type BadgeID string

const (
	BadgeFirstSession    BadgeID = "first_session"
	BadgeStreak3         BadgeID = "streak_3"
	BadgeStreak7         BadgeID = "streak_7"
	BadgeStreak30        BadgeID = "streak_30"
	BadgePerfectSession  BadgeID = "perfect_session"
	BadgeAnswers100      BadgeID = "answers_100"
	BadgeTopicsMastered5 BadgeID = "topics_mastered_5"
	BadgeFastSolver      BadgeID = "fast_solver"
)

// BadgeDefinition es la definición inmutable de un logro (D-16: catálogo en
// código, sin tabla).
type BadgeDefinition struct {
	ID          BadgeID
	Name        string
	Description string
	// Icon es un nombre de MaterialCommunityIcons con prefijo `mdi-`.
	// Vuetify lo usa tal cual; mobile le quita el prefijo.
	Icon   string
	Target int
}

// badgeCatalog es el orden en el que se muestran los logros.
var badgeCatalog = []BadgeDefinition{
	{
		ID:          BadgeFirstSession,
		Name:        "Primera sesión",
		Description: "Completa tu primera sesión",
		Icon:        "mdi-flag-checkered",
		Target:      1,
	},
	{
		ID:          BadgeStreak3,
		Name:        "Racha de 3",
		Description: "3 días seguidos de estudio",
		Icon:        "mdi-fire",
		Target:      3,
	},
	{
		ID:          BadgeStreak7,
		Name:        "Racha de 7",
		Description: "7 días seguidos de estudio",
		Icon:        "mdi-fire",
		Target:      7,
	},
	{
		ID:          BadgeStreak30,
		Name:        "Racha de 30",
		Description: "30 días seguidos de estudio",
		Icon:        "mdi-fire",
		Target:      30,
	},
	{
		ID:          BadgePerfectSession,
		Name:        "Sesión perfecta",
		Description: "Termina una sesión con el 100% de acierto (mín. 5 preguntas)",
		Icon:        "mdi-star",
		Target:      1,
	},
	{
		ID:          BadgeAnswers100,
		Name:        "100 respuestas",
		Description: "Responde 100 preguntas",
		Icon:        "mdi-check-all",
		Target:      100,
	},
	{
		ID:          BadgeTopicsMastered5,
		Name:        "5 temas dominados",
		Description: "Domina 5 temas",
		Icon:        "mdi-school",
		Target:      5,
	},
	{
		ID:          BadgeFastSolver,
		Name:        "Resolvedor veloz",
		Description: "Sesión perfecta con menos de 10 s por pregunta",
		Icon:        "mdi-lightning-bolt",
		Target:      1,
	},
}
