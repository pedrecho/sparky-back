package recorder

import (
	"net/http"
)

type ResponseRecorder struct {
	http.ResponseWriter
	Status int
	Body   []byte
}

func NewResponseRecorder(w http.ResponseWriter) *ResponseRecorder {
	return &ResponseRecorder{
		ResponseWriter: w,
	}
}

func (r *ResponseRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *ResponseRecorder) Write(data []byte) (int, error) {
	r.Body = append(r.Body, data...)
	return r.ResponseWriter.Write(data)
}

func (r *ResponseRecorder) IsDefault() bool {
	return r.Status == 0
}

func (r *ResponseRecorder) GetBody() string {
	return string(r.Body)
}
