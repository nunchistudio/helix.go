package bucket

import (
	"fmt"

	"go.nunchi.studio/helix/telemetry/trace"
)

/*
setDefaultAttributes sets integration attributes to a trace span.
*/
func setDefaultAttributes(span *trace.Span, cfg *Config) {
	if cfg != nil {
		span.SetStringAttribute(fmt.Sprintf("%s.driver", identifier), cfg.Driver.string())
		span.SetStringAttribute(fmt.Sprintf("%s.bucket", identifier), cfg.Bucket)

		if cfg.Subfolder != "" {
			span.SetStringAttribute(fmt.Sprintf("%s.subfolder", identifier), cfg.Subfolder)
		}
	}
}

/*
setKeyAttributes sets blob's key attributes to a trace span.
*/
func setKeyAttributes(span *trace.Span, key string) {
	span.SetStringAttribute(fmt.Sprintf("%s.key", identifier), key)
}
