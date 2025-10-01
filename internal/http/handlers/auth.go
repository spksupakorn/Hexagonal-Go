package handlers

import (
	"dungeons-dragon-service/internal/dto"
	"dungeons-dragon-service/internal/http/custom"
	usecase "dungeons-dragon-service/internal/usecases"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	uc usecase.AuthUseCase
	v  *validator.Validate
}

func NewAuthHandler(uc usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{uc: uc, v: validator.New()}
}

// Login godoc
// @Summary      Login get token
// @Description  Authenticates a user and returns a JWT token.
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Param        loginRequest  body      dto.LoginRequest  true  "Login Request"
// @Success      200           {object}  dto.APIObjectResponse{data=dto.LoginResponse}  "Successful login"
// @Failure      400           {object}  dto.APIErrorResponse{data=interface{}}  "Invalid request"
// @Failure      401           {object}  dto.APIErrorResponse{data=interface{}}  "Unauthorized"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	defer custom.PanicController(c)
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		e := custom.NewBadRequestError("invalid request body")
		custom.PanicException(e)
	}
	if err := h.v.Struct(req); err != nil {
		e := custom.NewValidationError("required fields are missing or invalid")
		custom.PanicException(e)
	}
	token, err := h.uc.Login(req.Username, req.Password)
	if err != nil {
		custom.PanicException(err)
	}
	return c.JSON(http.StatusOK, custom.BuildResponse(custom.Success, token))
}
