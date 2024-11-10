package rest

import (
	"net/http"

	"golang.org/x/text/language"
)

/*
addFrenchLocalesForTest adds locales in French used in response writer tests.
*/
func addFrenchLocalesForTest() {
	AddOrEditLanguage(language.French, map[int]string{
		http.StatusBadRequest:            "Échec de la validation de la requête",
		http.StatusUnauthorized:          "Vous n'êtes pas autorisé à effectuer cette action",
		http.StatusPaymentRequired:       "La requête a échoué car un paiement est requis",
		http.StatusForbidden:             "Vous n'avez pas les permissions requises pour effectuer cette action",
		http.StatusNotFound:              "La ressource n'existe pas",
		http.StatusMethodNotAllowed:      "La ressource ne supporte pas cette méthode",
		http.StatusConflict:              "Échec du traitement de la requête en raison d'un conflit",
		http.StatusRequestEntityTooLarge: "Impossible de traiter une requête avec un payload trop large",
		http.StatusTooManyRequests:       "La limite du taux de requêtes a été atteinte",
		http.StatusInternalServerError:   "Nous avons été informés de cette erreur interne inattendue",
		http.StatusServiceUnavailable:    "Veuillez réessayer dans quelques instants",
	})
}
