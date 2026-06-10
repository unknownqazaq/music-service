package core_http_request

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "music-service/internal/core/errors"
)

// GetIntQueryParam reads query parameter and parses it as int.
// Returns nil if parameter is missing.
func GetIntQueryParam(r *http.Request, key string) (*int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	val, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s' by key='%s' not a valid integer: %v: %w",
			param,
			key,
			err,
			core_errors.ErrBadRequest,
		)
	}

	return &val, nil
}
