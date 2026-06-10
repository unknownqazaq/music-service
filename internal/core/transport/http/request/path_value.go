package core_http_request

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	core_errors "music-service/internal/core/errors"
)

// GetIntPathValue extracts a path parameter using chi.URLParam and parses it as int64.
func GetIntPathValue(r *http.Request, key string) (int64, error) {
	pathValue := chi.URLParam(r, key)
	if pathValue == "" {
		return 0, fmt.Errorf(
			"no key='%s' in path values: %w",
			key,
			core_errors.ErrBadRequest,
		)
	}

	val, err := strconv.ParseInt(pathValue, 10, 64)
	if err != nil {
		return 0, fmt.Errorf(
			"path value='%s' by key='%s' not a valid integer: %v: %w",
			pathValue,
			key,
			err,
			core_errors.ErrBadRequest,
		)
	}

	return val, nil
}
