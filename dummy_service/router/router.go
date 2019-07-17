package router

import (
	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"

	dummy_service "github.com/ivch/dummy-service"
)

func New(svc dummy_service.Service) *fasthttprouter.Router {
	router := fasthttprouter.New()

	router.POST("/event", func(ctx *fasthttp.RequestCtx) {
		res, err := svc.Create(ctx.PostBody())
		if err != nil {
			log.Println(err)
			ctx.Error(err.Error(), fasthttp.StatusNotFound)

			return
		}

		ctx.SetContentType("application/json; charset=utf8")
		ctx.Response.SetBody(res)
	})

	return router
}
