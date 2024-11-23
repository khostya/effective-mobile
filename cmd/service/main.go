package main

import (
	"context"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	openapi2 "github.com/khostya/effective-mobile/cmd/service/openapi"
	"github.com/khostya/effective-mobile/pkg/api"
	"github.com/khostya/effective-mobile/pkg/httpserver"
	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"
	"log"
	"net/http"
	"time"
)

const port = "8079"

func main() {
	openapi, err := openapi2.GetOpenapiV3(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	options := &nethttpmiddleware.Options{
		SilenceServersWarning: true,
		Options: openapi3filter.Options{
			ExcludeRequestBody: true,
		},
	}

	handler := api.HandlerWithOptions(server{}, api.StdHTTPServerOptions{
		Middlewares: []api.MiddlewareFunc{
			cors.AllowAll().Handler,
			nethttpmiddleware.OapiRequestValidatorWithOptions(openapi, options),
		},
	})

	httpserver := httpserver.New(handler,
		httpserver.Port(port),
	)

	httpserver.Start()
	<-httpserver.Notify()
}

var _ api.ServerInterface = server{}

type server struct {
}

func (s server) GetInfo(w http.ResponseWriter, r *http.Request, params api.GetInfoParams) {
	render.Status(r, http.StatusOK)
	render.JSON(w, r, api.SongDetail{
		ReleaseDate: time.Now().Format(time.RFC3339),
		Text:        uuid.NewString(),
		Link:        uuid.NewString(),
	})
}
