package core

type HTTPResponseBodyConfig struct {
}

type HTTPResponseConfig struct {
	Body    HTTPResponseBodyConfig
	Headers map[string]string
}
