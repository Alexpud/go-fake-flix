package openapi

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
)

// Handler converts a Swagger 2.0 (OpenAPI 2.0) spec file to OpenAPI 3 and serves it as JSON.
// This is useful when your generator outputs swagger.json but your tooling/UI prefers OpenAPI 3.
func Handler(swagger2Path string) http.HandlerFunc {
	var (
		once     sync.Once
		cached   []byte
		cachedEr error
	)

	return func(w http.ResponseWriter, r *http.Request) {
		once.Do(func() {
			raw, err := os.ReadFile(swagger2Path)
			if err != nil {
				cachedEr = err
				return
			}

			var doc2 openapi2.T
			if err := json.Unmarshal(raw, &doc2); err != nil {
				cachedEr = err
				return
			}

			// ToV3 expects a pointer.
			doc3, err := openapi2conv.ToV3(&doc2)
			if err != nil {
				cachedEr = err
				return
			}

			cached, cachedEr = json.MarshalIndent(doc3, "", "  ")
		})

		if cachedEr != nil {
			http.Error(w, cachedEr.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write(cached)
	}
}
