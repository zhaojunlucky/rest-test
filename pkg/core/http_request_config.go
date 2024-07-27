package core

type HTTPRequestBodyConfig struct {
}
type HTTPRequestConfig struct {
	URLParams map[string]string
	Headers   map[string]string
	Body      HTTPRequestBodyConfig
	URL       string
}
