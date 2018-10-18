package server

import (
	"encoding/json"
	"net/http"

	"github.com/skygeario/skygear-server/pkg/core/middleware"
	nextSkyerr "github.com/skygeario/skygear-server/pkg/core/skyerr"
	"github.com/skygeario/skygear-server/pkg/server/skyerr"
)

type Option struct {
	RecoverPanic        bool
	RecoverPanicHandler middleware.RecoverHandler
}

// RecoveredResponse is interface for the default RecoverPanicHandler to write response
type RecoveredResponse struct {
	Err skyerr.Error `json:"error,omitempty"`
}

func DefaultOption() Option {
	return Option{
		RecoverPanic: true,
		RecoverPanicHandler: func(w http.ResponseWriter, r *http.Request, err skyerr.Error) {
			httpStatus := nextSkyerr.ErrorDefaultStatusCode(err)

			// TODO: log

			response := RecoveredResponse{Err: err}
			encoder := json.NewEncoder(w)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(httpStatus)
			encoder.Encode(response)
		},
	}
}