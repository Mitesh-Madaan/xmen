package router

type RequestIDKey string

type RespDetail struct {
	statusCode int
	message    []byte
}
