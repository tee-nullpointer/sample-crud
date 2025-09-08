package middleware

import (
	"errors"
	"net/http"
	customerrors "sample-crud/pkg/errors"
	"sample-crud/pkg/response"

	"github.com/gin-gonic/gin"
)

func ErrorRecover() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		errs := c.Errors
		if len(errs) > 0 {
			if err := errs.Last().Err; err != nil {
				var customErr *customerrors.CustomError
				switch {
				case errors.As(err, &customErr):
					c.JSON(customErr.HttpStatus, response.Error(customErr.Code, customErr.Message))
				default:
					c.JSON(http.StatusInternalServerError, response.InternalServerError())
				}
			}
		}
	}
}
