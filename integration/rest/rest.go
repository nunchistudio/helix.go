package rest

import (
	"net/http"

	"go.nunchi.studio/helix/errorstack"
	"go.nunchi.studio/helix/service"

	"github.com/getkin/kin-openapi/routers"
	"github.com/uptrace/bunrouter"
)

/*
REST exposes the HTTP REST API functions.
*/
type REST interface {
	GET(path string, handler http.HandlerFunc)
	HEAD(path string, handler http.HandlerFunc)
	DELETE(path string, handler http.HandlerFunc)
	OPTIONS(path string, handler http.HandlerFunc)
	PATCH(path string, handler http.HandlerFunc)
	POST(path string, handler http.HandlerFunc)
	PUT(path string, handler http.HandlerFunc)
}

/*
rest represents the rest integration. It respects the integration.Integration
and REST interfaces.
*/
type rest struct {

	// config holds the Config initially passed when creating a new REST API.
	config *Config

	// bun is the underlying router. This package has been designed to easily
	// switch from one underlying router to another if necessary, in case one goes
	// unmaintained or doesn't meet our requirements anymore.
	bun *bunrouter.CompatRouter

	// server is the standard http.Server used to serve HTTP requests.
	server *http.Server

	// oapirouter is the OpenAPI router used to validate requests and responses
	// against the OpenAPI description passed in Config.
	oapirouter routers.Router
}

/*
New tries to build a new HTTP API server for Config. Returns an error if Config
or OpenAPI description are not valid.
*/
func New(cfg Config) (REST, error) {

	// No need to continue if Config is not valid.
	err := cfg.sanitize()
	if err != nil {
		return nil, err
	}

	// Start to build an error stack, so we can add validations as we go.
	stack := errorstack.New("Failed to initialize integration", errorstack.WithIntegration(identifier))
	r := &rest{
		config: &cfg,
	}

	var validations []errorstack.Validation
	r.bun, validations = r.buildRouter()
	if validations != nil {
		stack.WithValidations(validations...)
	}

	// Only try to build the OpenAPI router if enabled in Config.
	if cfg.OpenAPI.Enabled {
		r.oapirouter, validations = r.buildRouterOpenAPI()
		if validations != nil {
			stack.WithValidations(validations...)
		}
	}

	// Stop here if error validations were encountered.
	if stack.HasValidations() {
		return nil, stack
	}

	// Otherwise, try to attach the integration to the service.
	if err := service.Attach(r); err != nil {
		return nil, err
	}

	return r.bun, nil
}
