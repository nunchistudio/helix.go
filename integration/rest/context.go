package rest

import (
	"context"

	"github.com/uptrace/bunrouter"
)

/*
ParamsFromContext returns request's params found in the context passed, if any.
Returns true if some params are present, false otherwise.
*/
func ParamsFromContext(ctx context.Context) (map[string]string, bool) {
	params := bunrouter.ParamsFromContext(ctx).Map()

	var found bool
	if len(params) > 0 {
		found = true
	}

	return params, found
}
