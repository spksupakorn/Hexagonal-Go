package custom

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func PanicController(c echo.Context) error {
	if err := recover(); err != nil {
		str := fmt.Sprint(err)
		parts := strings.SplitN(str, ":", 2)

		code := http.StatusInternalServerError
		msg := "unexpected error"
		if len(parts) == 2 {
			if parsedCode, parseErr := strconv.Atoi(strings.TrimSpace(parts[0])); parseErr == nil {
				code = parsedCode
				msg = strings.TrimSpace(parts[1])
			} else {
				msg = str
			}
		} else {
			msg = str
		}

		var respStatus interface{}
		switch code {
		case http.StatusNotFound:
			respStatus = DataNotFound.GetResponseStatus()
		case http.StatusBadRequest:
			respStatus = BadRequest.GetResponseStatus()
		case http.StatusUnauthorized:
			respStatus = Unauthorized.GetResponseStatus()
		case http.StatusForbidden:
			respStatus = Forbidden.GetResponseStatus()
		case http.StatusUnprocessableEntity:
			respStatus = UnprocessableEntity.GetResponseStatus()
		case http.StatusConflict:
			respStatus = Conflict.GetResponseStatus()
		case http.StatusNoContent:
			respStatus = NoContent.GetResponseStatus()
		default:
			respStatus = InternalServerError.GetResponseStatus()
		}

		return c.JSON(code, BuildResponse_(respStatus.(bool), msg, Null()))
	}
	return nil
}

func PanicException(err error) {
	switch e := err.(type) {
	case *AppError:
		PanicException_(e.Code, e.Message)
	case validator.ValidationErrors:
		for _, ve := range e {
			fieldErr := fmt.Sprintf("'%s' field is required", ve.Field())
			PanicException_(http.StatusBadRequest, fieldErr)
		}
	default:
		PanicException_(http.StatusInternalServerError, "unexpected error")
	}
}

func PanicException_(statusCode int, message string) {
	panic(fmt.Errorf("%d: %s", statusCode, message))
}
