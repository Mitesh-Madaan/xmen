package router

type RequestIDKey string

type RespDetail struct {
	Statuscode int
	Message    []byte `default:""`
	Type       string `default:"text/plain"`
}
