package request

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func ReadIDParam(r *http.Request) (int, error) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func ReadInt(r *http.Request, v *Valid, key string, defaultValue int) int {
	valueStr := r.URL.Query().Get(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		v.AddError(key, fmt.Sprintf("query parameter '%s' must be an integer", key))
		return defaultValue
	}
	return value
}
