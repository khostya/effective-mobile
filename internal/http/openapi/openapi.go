package openapi

import (
	"context"
	"github.com/getkin/kin-openapi/openapi3"
)

func GetOpenapiV3(ctx context.Context) (*openapi3.T, error) {
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	return loader.LoadFromFile("./api/openapi.yaml")
}
