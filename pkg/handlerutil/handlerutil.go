package handlerutil

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
