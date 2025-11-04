package helper

import (
	"net/http"

	"github.com/ssongin/handlers/pkg/model"
)

func SetHeaders(w http.ResponseWriter, headerGroups ...[]model.Header) {
	for _, group := range headerGroups {
		for _, h := range group {
			w.Header().Set(h.Name, h.Value)
		}
	}
}

func HasHeader(headers []model.Header, name string) bool {
	for _, h := range headers {
		if http.CanonicalHeaderKey(h.Name) == http.CanonicalHeaderKey(name) {
			return true
		}
	}
	return false
}
