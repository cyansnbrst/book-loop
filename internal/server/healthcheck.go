package server

import (
	"net/http"

	"bookloop.net/config"
	"bookloop.net/pkg/utils"
)

func healthcheckHandler(s *Server) http.HandlerFunc {
	return (func(w http.ResponseWriter, r *http.Request) {
		env := utils.Envelope{
			"status": "available",
			"system_info": map[string]string{
				"environment": s.Config.Env,
				"version":     config.Version,
			},
		}

		err := utils.WriteJSON(w, http.StatusOK, env, nil)
		if err != nil {
			utils.ServerErrorResponse(w, r, s.Logger, err)
		}
	})
}
