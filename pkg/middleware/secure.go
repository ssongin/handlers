package middleware

import (
	"fmt"
	"net/http"
)

func BlockHttp(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.TLS == nil {
			http.Error(w, "HTTP access is not allowed. Use HTTPS.", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func BlockHttps(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.TLS != nil {
			http.Error(w, "HTTPS access is not allowed. Use HTTP.", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RedirectToHttps(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.TLS == nil {
			target := fmt.Sprintf("https://%s%s", r.Host, r.URL.RequestURI())
			http.Redirect(w, r, target, http.StatusMovedPermanently)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RedirectToHttp(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.TLS != nil {
			target := fmt.Sprintf("http://%s%s", r.Host, r.URL.RequestURI())
			http.Redirect(w, r, target, http.StatusMovedPermanently)
			return
		}
		next.ServeHTTP(w, r)
	})
}
