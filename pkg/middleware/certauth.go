package middleware

import (
	"context"
	"crypto/x509"
	"net/http"
	"time"
)

type CertInfo struct {
	CommonName   string
	SerialNumber string
	Issuer       string
	NotBefore    time.Time
	NotAfter     time.Time
}

// Context key type to avoid collisions
type certInfoKeyType struct{}

var certInfoKey = certInfoKeyType{}

// CertAuthMiddleware validates client certificates and attaches info to context
func CertAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.TLS == nil || len(r.TLS.PeerCertificates) == 0 {
			http.Error(w, "Client certificate required", http.StatusUnauthorized)
			return
		}

		clientCert := r.TLS.PeerCertificates[0]
		if err := validateCertificate(clientCert); err != nil {
			http.Error(w, "Invalid client certificate: "+err.Error(), http.StatusUnauthorized)
			return
		}

		info := CertInfo{
			CommonName:   clientCert.Subject.CommonName,
			SerialNumber: clientCert.SerialNumber.String(),
			Issuer:       clientCert.Issuer.CommonName,
			NotBefore:    clientCert.NotBefore,
			NotAfter:     clientCert.NotAfter,
		}

		ctx := context.WithValue(r.Context(), certInfoKey, info)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func validateCertificate(cert *x509.Certificate) error {
	now := time.Now()
	if now.Before(cert.NotBefore) {
		return http.ErrAbortHandler // Not yet valid
	}
	if now.After(cert.NotAfter) {
		return http.ErrAbortHandler // Expired
	}
	return nil
}

// Helper to retrieve certificate info later in your handlers
func GetCertInfo(r *http.Request) (*CertInfo, bool) {
	info, ok := r.Context().Value(certInfoKey).(CertInfo)
	return &info, ok
}
