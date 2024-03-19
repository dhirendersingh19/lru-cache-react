package handlers

import (
	"lru-cache/controllers"
	"lru-cache/restapi/operations"

	"github.com/go-openapi/runtime/middleware"
)

func LRUCacheGETHandler(params operations.LrucacheGetParams) middleware.Responder {
	var model interface{}
	res, err, code := controllers.GET(params.ID)
	if err != nil {
		model = err
	} else {
		model = res
	}
	var restResponder = &responder{
		ResponseCode: code,
		Payload:      &model,
	}
	return restResponder
}
