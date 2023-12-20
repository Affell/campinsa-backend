package exceptions

import "github.com/kataras/iris/v12"

func StatusCode(code int) (int, map[string]interface{}) {
	data := make(map[string]interface{})
	switch code {
	case iris.StatusNotImplemented:
		data["message"] = "Ce contenu n'est pas encore implémenté"
	case iris.StatusUnauthorized:
		data["message"] = "Votre session à expirée"
	case iris.StatusForbidden:
		data["message"] = "Vous n'êtes pas autorisé(e) à consulter ce contenu"
	case iris.StatusNotFound:
		data["message"] = "L'objet est introuvable"
	case iris.StatusBadRequest:
		data["message"] = "Requête invalide"
	case iris.StatusInternalServerError:
		data["message"] = "Erreur interne du serveur"
	case iris.StatusOK:
		data["message"] = "Action réalisée avec succès"
	case iris.StatusNoContent:
		data = nil
	default:
		data["message"] = "Erreur inconnue"
	}

	return code, data
}
