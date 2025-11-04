package middleware

import (
	"net/http"
	"sort"
)

type Middleware func(http.Handler) http.Handler

type Entry struct {
	Name       string
	Middleware Middleware
	Weight     int
}

var registry = map[string]Entry{}

func Register(name string, mw Middleware, weight int) {
	registry[name] = Entry{Name: name, Middleware: mw, Weight: weight}
}

func ApplyChain(h http.Handler, names []string) http.Handler {
	var entries []Entry
	for _, name := range names {
		if e, ok := registry[name]; ok {
			entries = append(entries, e)
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Weight < entries[j].Weight
	})

	for i := len(entries) - 1; i >= 0; i-- {
		h = entries[i].Middleware(h)
	}
	return h
}
