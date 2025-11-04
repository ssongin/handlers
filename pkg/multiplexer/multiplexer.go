package multiplexer

import (
	"net/http"

	"github.com/ssongin/core"
	"github.com/ssongin/handlers/pkg/middleware"
	"github.com/ssongin/handlers/pkg/model"
)

func LoadMultiplexers(gh model.GeneratedHandlers, config *model.Configuration) model.GeneratedHandlers {
	for _, muxDef := range config.Multiplexers {
		mux := http.NewServeMux()
		for _, handlerName := range muxDef.Handlers {
			generatedHandler, exists := gh[handlerName]
			if !exists {
				core.GetLogger().Error("Multiplexer population: handler not found", "handlerName", handlerName)
				continue
			}
			mux.Handle(generatedHandler.GetPrefix()+"/", http.StripPrefix(generatedHandler.GetPrefix(), generatedHandler))
		}
		namedMux := &model.NamedMux{
			Name: muxDef.Name,
			Mux:  mux,
		}

		namedMux.Mux = middleware.ApplyChain(namedMux.Mux, muxDef.Middleware).(*http.ServeMux)

		gh[muxDef.Name] = namedMux
		core.GetLogger().Info("Created multiplexer", "name", muxDef.Name, "numHandlers", len(muxDef.Handlers))
	}
	return gh
}
