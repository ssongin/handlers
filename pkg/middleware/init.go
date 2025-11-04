package middleware

func init() {
	Register("certAuth", CertAuthMiddleware, 100)
	Register("logger", LogRequestMiddleware, 900)
}
