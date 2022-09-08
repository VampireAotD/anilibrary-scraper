package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	writer http.ResponseWriter
}

func New(writer http.ResponseWriter) *Response {
	return &Response{writer: writer}
}

func (r Response) JSON(code int, data any) error {
	r.writer.Header().Set("Content-Type", "application/json")
	r.writer.WriteHeader(code)

	return json.NewEncoder(r.writer).Encode(data)
}
