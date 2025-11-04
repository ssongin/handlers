package inmemorypages

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ssongin/core"
	"github.com/ssongin/handlers/pkg/helper"
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

		inMemory.Endpoints[i].Headers = ensureContentType(inMemory.Endpoints[i].Headers, inMemory.Endpoints[i].Source)
	}

	for _, endpoint := range inMemory.Endpoints {
		e := endpoint

		handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			helper.SetHeaders(w, e.Headers, inMemory.CommonHeaders)
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

func detectMimeType(path string) string {
	ext := filepath.Ext(path)
	if ext == "" {
		return "application/octet-stream"
	}
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		return "application/octet-stream"
	}
	return mimeType
}

func ensureContentType(headers []model.Header, source string) []model.Header {
	if helper.HasHeader(headers, "Content-Type") {
		return headers
	}
	ct := detectMimeType(source)
	headers = append(headers, model.Header{
		Name:  "Content-Type",
		Value: ct,
	})
	return headers
}
