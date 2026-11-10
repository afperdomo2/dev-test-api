package sessions

import (
	"hash/fnv"
	"sort"

	"github.com/felipe/dev-test-api/internal/models"
	"github.com/google/uuid"
)

// shuffleOptions reordena las opciones de una pregunta de forma determinista.
//
// No usa math/rand a propósito: el mismo (sesión, pregunta) debe devolver
// siempre el mismo orden aunque Postgres devuelva las filas en distinto orden
// físico (los preloads no llevan ORDER BY). La clave es un hash estable de
// sesión + pregunta + opción, no la posición de entrada: así el resultado es
// independiente del orden en que lleguen las filas desde la DB.
//
// El objetivo es romper el sesgo del banco (la IA y el seed tienden a dejar la
// respuesta correcta primero, y el modelo repite el patrón por lote). Al
// presentar, cada sesión ve las opciones en un orden distinto.
func shuffleOptions(opts []models.QuestionOption, sessionID, questionID uuid.UUID) []models.QuestionOption {
	if len(opts) < 2 {
		return opts
	}

	seed := sessionID.String() + "|" + questionID.String()
	out := make([]models.QuestionOption, len(opts))
	copy(out, opts)
	sort.Slice(out, func(i, j int) bool {
		ki := optionOrderKey(seed, out[i].ID)
		kj := optionOrderKey(seed, out[j].ID)
		if ki != kj {
			return ki < kj
		}
		return out[i].ID.String() < out[j].ID.String()
	})
	return out
}

func optionOrderKey(seed string, optionID uuid.UUID) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(seed))
	_, _ = h.Write([]byte("|"))
	_, _ = h.Write([]byte(optionID.String()))
	return h.Sum64()
}

// optionsForPresentation aplica el orden que verá el usuario. `true_false` se
// deja intacto: "Verdadero"/"Falso" es un orden convencional y su posición ya
// está equilibrada en el banco.
func optionsForPresentation(qType string, opts []models.QuestionOption, sessionID, questionID uuid.UUID) []models.QuestionOption {
	if qType == models.QuestionTypeTrueFalse {
		return opts
	}
	return shuffleOptions(opts, sessionID, questionID)
}
