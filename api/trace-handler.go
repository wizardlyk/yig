package api

import (
	"github.com/journeymidnight/yig/helper"
	"net/http"
)

type traceHandler struct {
	handler     http.Handler
	objectLayer ObjectLayer
}

func (t traceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	helper.TracerLogger.Println(5, "enter traceHandler")
	t.handler.ServeHTTP(w, r)
	helper.TracerLogger.Println(5, "exit  traceHandler")
}

func SetTraceHandler(handler http.Handler, objectLayer ObjectLayer) http.Handler {
	return traceHandler{
		handler:     handler,
		objectLayer: objectLayer,
	}
}
