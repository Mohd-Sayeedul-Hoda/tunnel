package encoding

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/api/request"
)

var ErrInvalidData = errors.New("invalid data")
var ErrInvalidRequest = errors.New("invalid request")

func EncodeJson[T any](w http.ResponseWriter, r *http.Request, status int, data T) error {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err := encoder.Encode(data)
	if err != nil {
		return err
	}

	return nil

}

func Decode[T any](w http.ResponseWriter, r *http.Request, data *T) error {

	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)

	dec.DisallowUnknownFields()

	err := dec.Decode(data)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("%w: body contain badly-formated JSON (at character %d)", ErrInvalidRequest, syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return fmt.Errorf("%w: body contains badly-formated JSON", ErrInvalidRequest)

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("%w: body contains incorrect JSON type for field %q", ErrInvalidRequest, unmarshalTypeError.Field)
			}
			return fmt.Errorf("%w: body contians incorrect JSON type (at character %d)", ErrInvalidRequest, unmarshalTypeError.Offset)

		// happen when body is empty
		case errors.Is(err, io.EOF):
			return fmt.Errorf("%w: body must not be empty", ErrInvalidRequest)

		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field")
			return fmt.Errorf("%w: body contains unknown key %s", ErrInvalidRequest, fieldName)

		// happen when body size bigger then specified size
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("%w: body must not be larger than %d bytes", ErrInvalidRequest, maxBytesError.Limit)

			// invalid argument when pass to decode function
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	return nil
}

func Validated[T request.Validator](w http.ResponseWriter, r *http.Request, v *request.Valid, data T) error {
	if err := Decode(w, r, &data); err != nil {
		return err
	}

	problem := data.Valid(r.Context(), v)
	if !problem.Valid() {
		return fmt.Errorf("%w: invalid %T: %d problems", ErrInvalidData, data, len(problem.Errors))
	}

	return nil
}
