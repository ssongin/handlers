package model

import "net/http"

type GeneratedHandlers map[string]DescribedHandler

type Configuration struct {
	Pages        Pages         `yaml:"pages"`
	Multiplexers []Multiplexer `yaml:"multiplexer"`
}

type Multiplexer struct {
	Name       string   `yaml:"name"`
	Handlers   []string `yaml:"handlers"`
	Middleware []string `yaml:"middleware"`
}

type Pages struct {
	InMemoryPages []InMemoryMux `yaml:"inmemory"`
}

type InMemoryMux struct {
	Name          string             `yaml:"name"`
	Prefix        string             `yaml:"prefix"`
	Endpoints     []InMemoryEndpoint `yaml:"endpoints"`
	CommonHeaders []Header           `yaml:"commonHeaders"`
	Middleware    []string           `yaml:"middleware"`
	Mux           *http.ServeMux     // not part of YAML, value comes from source file
}

type InMemoryEndpoint struct {
	Source     string   `yaml:"source"`
	Path       string   `yaml:"path"`
	Headers    []Header `yaml:"headers"`
	Middleware []string `yaml:"middleware"`
	Content    []byte   // not part of YAML, value comes from source file
}

type Header struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type DescribedHandler interface {
	http.Handler
	GenName() string
	GetPrefix() string
}

func (im *InMemoryMux) GenName() string                                  { return im.Name }
func (im *InMemoryMux) GetPrefix() string                                { return im.Prefix }
func (im *InMemoryMux) ServeHTTP(w http.ResponseWriter, r *http.Request) { im.Mux.ServeHTTP(w, r) }

type NamedMux struct {
	Name   string
	Prefix string
	Mux    *http.ServeMux
}

func (n *NamedMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	n.Mux.ServeHTTP(w, r)
}

func (n *NamedMux) GenName() string   { return n.Name }
func (n *NamedMux) GetPrefix() string { return n.Prefix }
