package handlers

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

type responder struct {
	Payload      interface{} `json:"body,omitempty"`
	ResponseCode int
}

func (r *responder) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {
	rw.WriteHeader(r.ResponseCode)
	if r.Payload != nil {
		payload := r.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err)
		}
	}
}
