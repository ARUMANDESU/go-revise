package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ARUMANDESU/go-revise/pkg/errs"
)

type envelope map[string]any

// maxRequestBodySize is the maximum size of a request body.
// NOTE: on changing this value, don't forget to update the error message in readJSON function.
const maxRequestBodySize = 10 << 20 // 10MB

func readJSON(w http.ResponseWriter, r *http.Request, v any) error {
	op := errs.Op("handler.readJSON")

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(v)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return errs.
				NewIncorrectInputError(op, err, "badly-formed JSON").
				WithMessages([]errs.Message{{
					Key:   "message",
					Value: fmt.Sprintf("badly-formed JSON (at character %d)", syntaxError.Offset),
				}}).
				WithContext("body", r.Body).
				WithContext("offset", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errs.
				NewIncorrectInputError(op, err, "body contains badly-formed JSON").
				WithMessages([]errs.Message{{Key: "message", Value: "body contains badly-formed JSON"}}).
				WithContext("body", r.Body)
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return errs.
					NewIncorrectInputError(op, err, "body contains invalid JSON").
					WithMessages([]errs.Message{{
						Key: "message",
						Value: fmt.Sprintf(
							"body contains invalid JSON (at character %d)",
							unmarshalTypeError.Offset,
						),
					}}).
					WithContext("body", r.Body).
					WithContext("offset", unmarshalTypeError.Offset).
					WithContext("field", unmarshalTypeError.Field)
			}
			return errs.
				NewIncorrectInputError(op, err, "body contains invalid JSON").
				WithMessages([]errs.Message{{
					Key: "message",
					Value: fmt.Sprintf(
						"body contains invalid JSON (at character %d)",
						unmarshalTypeError.Offset,
					),
				}}).
				WithContext("body", r.Body).
				WithContext("offset", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errs.
				NewIncorrectInputError(op, err, "body must not be empty").
				WithMessages([]errs.Message{{Key: "message", Value: "body must not be empty"}}).
				WithContext("body", r.Body)
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return errs.
				NewIncorrectInputError(op, err, "body contains unknown field").
				WithMessages([]errs.Message{{
					Key:   "message",
					Value: fmt.Sprintf("body contains unknown field %s", fieldName),
				}}).
				WithContext("body", r.Body).
				WithContext("field", fieldName)
		case errors.As(err, &maxBytesError):
			return errs.
				NewIncorrectInputError(op, err, "body must not be larger than 10MB").
				WithMessages([]errs.Message{{Key: "message", Value: "body must not be larger than 10MB"}}).
				WithContext("body", r.Body)
		case errors.As(err, &invalidUnmarshalError):
			return errs.
				NewUnknownError(op, err, "invalid unmarshal error").
				WithContext("body", r.Body)
		default:
			return errs.
				NewUnknownError(op, err, "failed to read JSON").
				WithContext("body", r.Body)
		}

	}

	// This is to ensure that the body contains only a single JSON value.
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errs.
			NewIncorrectInputError(op, err, "body must only contain a single JSON value").
			WithMessages([]errs.Message{{Key: "message", Value: "body must only contain a single JSON value"}}).
			WithContext("body", r.Body)
	}

	return nil
}

func writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	op := errs.Op("handler.writeJSON")
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return errs.
			NewUnknownError(op, err, "failed to marshal JSON").
			WithContext("data", data)
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	write, err := w.Write(js)
	if err != nil {
		return errs.
			NewUnknownError(op, err, "failed to write JSON").
			WithContext("data", data).
			WithContext("write", write)
	}
	return nil
}
