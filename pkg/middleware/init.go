package middleware

func init() {
	Register("secureHttp", BlockHttp, 200)
	Register("secureHttps", BlockHttps, 200)
	Register("redirectToHttp", RedirectToHttp, 300)
	Register("redirectToHttps", RedirectToHttps, 300)
	Register("certAuth", CertAuthMiddleware, 400)
	Register("logger", LogRequestMiddleware, 800)
}
