package inmemorypages

import (
	"net/http"
	"os"

	"github.com/ssongin/core"
	"github.com/ssongin/handlers/pkg/handlerutil"
	"github.com/ssongin/handlers/pkg/middleware"
	"github.com/ssongin/handlers/pkg/model"
)

func LoadInMemoryStaticPages(inMemory *model.InMemoryMux) string {
	if inMemory.Mux == nil {
		inMemory.Mux = http.NewServeMux()
	}
	for i := range inMemory.Endpoints {
		content, err := os.ReadFile(inMemory.Endpoints[i].Source)
		core.CheckError("In memory pages: failed to read source file", err, "source", inMemory.Endpoints[i].Source)
		inMemory.Endpoints[i].Content = content
	}

	for _, endpoint := range inMemory.Endpoints {
		e := endpoint

		handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerutil.SetHeaders(w, e.Headers, inMemory.CommonHeaders)
			w.Write(e.Content)
		})

		allMiddleware := append(e.Middleware, inMemory.Middleware...)
		handler := middleware.ApplyChain(handlerFunc, allMiddleware)

		inMemory.Mux.Handle(e.Path, handler)
		core.GetLogger().Info("Loaded route into subrouter for in memory static file", "route", e.Path, "source", e.Source)
	}

	core.GetLogger().Info("Created Serve Multiplexer", "name", inMemory.Name, "numHandlers", len(inMemory.Endpoints))
	return inMemory.Name
}
