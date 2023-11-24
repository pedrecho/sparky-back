package middlewares

import (
	"github.com/uptrace/bunrouter"
	"go.uber.org/zap"
	"net/http"
)

func Log(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		zap.S().With("path", req.RequestURI).With("remote_addr", req.RemoteAddr).Info("start request")
		//rec := recorder.NewResponseRecorder(w)
		err := next(w, req)
		//TODO
		if err != nil {
			zap.S().With("path", req.RequestURI).With("remote_addr", req.RemoteAddr).Error(err)
			w.WriteHeader(http.StatusBadRequest)
		} else {
			zap.S().With("path", req.RequestURI).With("remote_addr", req.RemoteAddr).Info("end request")
		}
		return err
	}
}
