package apierror

import (
	"errors"
	"github.com/go-chi/render"
	"log"
	"net/http"
)

var (
	ErrInternal = &APIError{
		HTTPStatusCode: http.StatusInternalServerError,
		Message:        "Internal Server Error",
	}
	ErrNotFound = &APIError{
		HTTPStatusCode: http.StatusNotFound,
		Message:        "Resource Not found",
	}
	ErrBadRequest = &APIError{
		HTTPStatusCode: http.StatusBadRequest,
		Message:        "Bad Request",
	}
	ErrRender = &APIError{
		HTTPStatusCode: http.StatusUnprocessableEntity,
		Message:        "Error Rendering Response",
	}
)

type APIError struct {
	HTTPStatusCode int    `json:"status"`
	Message        string `json:"message"`
}

func (e *APIError) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

type ServerError struct {
	Err    error
	apiErr *APIError
}

func (se *ServerError) Error() string {
	return se.Err.Error()
}

func Wrap(err error, apiErr *APIError) error {
	return &ServerError{
		Err:    err,
		apiErr: apiErr,
	}
}

func Handle(w http.ResponseWriter, r *http.Request, err error) {
	serverError := &ServerError{}
	if errors.As(err, &serverError) {
		log.Print(err.Error())
		render.Render(w, r, serverError.apiErr)
	} else {
		render.Render(w, r, ErrInternal)
	}
}
