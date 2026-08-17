package util

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrInvalidID is returned when an ID isn't a valid ObjectID hex string.
var ErrInvalidID = errors.New("invalid id")

// ErrNotFound is returned when a document doesn't exist or doesn't belong to the caller.
var ErrNotFound = errors.New("document not found")

// RespondError writes the appropriate status code and message for err.
func RespondError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalidID):
		JSONwithoutBody(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, ErrNotFound):
		JSONwithoutBody(c, http.StatusNotFound, err.Error())
	default:
		JSONwithoutBody(c, http.StatusInternalServerError, "internal error")
	}
}
