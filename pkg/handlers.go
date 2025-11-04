package handlers

import (
	"os"

	"github.com/ssongin/core"
	"github.com/ssongin/handlers/pkg/inmemorypages"
	"github.com/ssongin/handlers/pkg/model"
	"github.com/ssongin/handlers/pkg/multiplexer"

	"gopkg.in/yaml.v3"
)

func ProcessYAML(configPath string) model.GeneratedHandlers {
	data, err := os.ReadFile(configPath)
	core.CheckError("Failed to read YAML file", err, "file", configPath)

	var root model.Configuration
	err = yaml.Unmarshal(data, &root)
	core.CheckError("Failed to parse YAML", err, "file", configPath)
	generatedHandlers := processInMemoryPages(&root)
	return multiplexer.LoadMultiplexers(generatedHandlers, &root)
}

func processInMemoryPages(config *model.Configuration) model.GeneratedHandlers {
	generatedHandlers := make(model.GeneratedHandlers)

	for _, inMemory := range config.Pages.InMemoryPages {
		name := inmemorypages.LoadInMemoryStaticPages(&inMemory)
		core.GetLogger().Info("Loaded in memory static pages", "name", name)
		generatedHandlers[name] = &inMemory
	}
	return generatedHandlers
}
