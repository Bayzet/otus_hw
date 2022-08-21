package internalhttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/server/http/models"
)

// BadRequest ...
type BadRequest struct {
	StatusCode int
	Msg        string
}

// nolint:cyclop,varnamelen
func decodeJSONBody(r *http.Request, dst *models.Event) *BadRequest {
	ctype := strings.TrimSpace(r.Header.Get("Content-Type"))
	if ctype != "" && ctype != "application/json" {
		return &BadRequest{
			StatusCode: http.StatusUnsupportedMediaType,
			Msg:        "Content-Type header is not application/json",
		}
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var (
			syntaxError        *json.SyntaxError
			unmarshalTypeError *json.UnmarshalTypeError
		)

		switch {
		case errors.As(err, &syntaxError):
			return &BadRequest{
				StatusCode: http.StatusBadRequest,
				Msg:        fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset),
			}

		case errors.Is(err, io.ErrUnexpectedEOF):
			return &BadRequest{StatusCode: http.StatusBadRequest, Msg: "Request body contains badly-formed JSON"}

		case errors.As(err, &unmarshalTypeError):
			return &BadRequest{
				StatusCode: http.StatusBadRequest,
				Msg:        fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset),
			}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			return &BadRequest{
				StatusCode: http.StatusBadRequest,
				Msg:        fmt.Sprintf("Request body contains unknown field %s", strings.TrimPrefix(err.Error(), "json: unknown field ")),
			}

		case errors.Is(err, io.EOF):
			return &BadRequest{StatusCode: http.StatusBadRequest, Msg: "Request body must not be empty"}

		case err.Error() == "http: request body too large":
			return &BadRequest{StatusCode: http.StatusRequestEntityTooLarge, Msg: "Request body must not be larger than 1MB"}

		default:
			return &BadRequest{StatusCode: http.StatusInternalServerError, Msg: fmt.Sprint(err)}
		}
	}

	return nil
}

func goodResponse(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(v)
}

func errorResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte(err.Error()))
}

func errorsBadRequest(w http.ResponseWriter, mr *BadRequest) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(mr.StatusCode)
	_, _ = w.Write([]byte(mr.Msg))
}
