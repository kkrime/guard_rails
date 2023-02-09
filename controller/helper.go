package controller

import (
	"context"
	"encoding/json"
	"guard_rails/errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func getSuccessResponse() *response {
	return &response{
		Status: success,
	}
}

func GetFailResponse(err string) *response {
	return &response{
		Status: fail,
		Error:  err,
	}
}

func AbortAndError(c *gin.Context, err error) {
	c.Abort()

	switch err := err.(type) {
	case *errors.RestError:
		c.JSON(err.Code, GetFailResponse(err.Err))
	case validator.ValidationErrors:
		c.JSON(400, GetFailResponse(Invalid_Input))
	case *json.InvalidUnmarshalError:
		c.JSON(400, GetFailResponse(Invalid_Input))
	default:
		c.JSON(500, GetFailResponse(Internal_Error))
	}
}

func createLogger(gl getLogger, ctx context.Context) *logrus.Entry {
	return gl.getLogger().WithContext(ctx).WithField("requestID", ctx.Value("X-Request-ID"))
}
