package handlers

import (
	"net/http"
	"pet-clinic/utils"
)

func ErrorResponse(w http.ResponseWriter, message string, code int, err error) {
	if err != nil {
		utils.Log.WithError(err).Error(message)
	} else {
		utils.Log.Warn(message)
	}

	http.Error(w, message, code)
}
