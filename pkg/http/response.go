package http

import (
	sj "encoding/json"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/encoding/json"
	"github.com/go-kratos/kratos/v2/errors"
	nt "net/http"
	"strings"

	"github.com/go-kratos/kratos/v2/transport/http"
)

// DefaultResponseEncoder copy from http.DefaultResponseEncoder
func DefaultResponseEncoder(w http.ResponseWriter, r *http.Request, v interface{}) error {
	if v == nil {
		return nil
	}
	if rd, ok := v.(http.Redirector); ok {
		url, code := rd.Redirect()
		nt.Redirect(w, r, url, code)
		return nil
	}

	codec := encoding.GetCodec(json.Name) // ignore Accept Header
	data, err := codec.Marshal(v)
	if err != nil {
		return err
	}

	bs, _ := sj.Marshal(NewResponse(data))

	w.Header().Set("Content-Type", ContentType(codec.Name()))
	_, err = w.Write(bs)
	if err != nil {
		return err
	}
	return nil
}

// DefaultErrorEncoder copy from http.DefaultErrorEncoder.
func DefaultErrorEncoder(w http.ResponseWriter, r *http.Request, err error) {
	se := FromError(errors.FromError(err)) // change error to BaseResponse

	codec := encoding.GetCodec(json.Name) // ignore Accept header
	body, err := codec.Marshal(se)
	if err != nil {
		w.WriteHeader(nt.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", ContentType(codec.Name()))
	// w.WriteHeader(int(se.Code)) // ignore http status code
	_, _ = w.Write(body)
}

const (
	baseContentType = "application"
)

// ContentType returns the content-type with base prefix.
func ContentType(subtype string) string {
	return strings.Join([]string{baseContentType, subtype}, "/")
}

func NewResponse(data []byte) BaseResponse {
	return BaseResponse{
		Code:    0,
		Message: "success",
		Data:    data,
	}
}

func FromError(e *errors.Error) *BaseResponse {
	if e == nil {
		return nil
	}
	return &BaseResponse{
		Code:    e.Code,
		Message: e.Message,
	}
}

type BaseResponse struct {
	Code    int32         `json:"code"`
	Message string        `json:"message"`
	Data    sj.RawMessage `json:"data,omitempty"` //
}
