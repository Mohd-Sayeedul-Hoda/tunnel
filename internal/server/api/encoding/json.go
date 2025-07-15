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
			return fmt.Errorf("body contain badly-formated JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formated JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contians incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		// happen when body is empty
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		// happen when body size bigger then specified size
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	return nil
}

func Validated[T request.Validator](w http.ResponseWriter, r *http.Request, data T) (*request.Valid, error) {
	if err := Decode(w, r, &data); err != nil {
		return nil, err
	}

	problem := data.Valid(r.Context())
	if !problem.Valid() {
		return problem, fmt.Errorf("%w: invalid %T: %d problems", ErrInvalidData, data, len(problem.Errors))
	}

	return problem, nil
}
