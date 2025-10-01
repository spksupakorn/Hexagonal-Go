package handlers

import (
	"dungeons-dragon-service/internal/domain/model"
	"dungeons-dragon-service/internal/dto"
	"dungeons-dragon-service/internal/http/custom"
	middleware "dungeons-dragon-service/internal/http/middlewares"
	usecase "dungeons-dragon-service/internal/usecases"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CharacterHandler struct {
	uc usecase.CharacterUseCase
	v  *validator.Validate
}

func NewCharacterHandler(uc usecase.CharacterUseCase) *CharacterHandler {
	return &CharacterHandler{uc: uc, v: validator.New()}
}

// ListCharacters godoc
// @Summary      List characters
// @Description  Retrieves a list of characters for the authenticated user.
// @Tags         characters
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.APIObjectResponse{data=[]dto.CharacterResponse}  "List of characters"
// @Failure      401  {object}  dto.APIErrorResponse{data=interface{}}  "Unauthorized"
// @Router       /characters [get]
func (h *CharacterHandler) List(c echo.Context) error {
	defer custom.PanicController(c)
	auth := middleware.IsAuthenticated(c)
	list, err := h.uc.ListForUser(auth)
	if err != nil {
		custom.PanicException(err)
	}

	return c.JSON(http.StatusOK, custom.BuildResponse(custom.Success, list))
}

// CreateCharacter godoc
// @Summary      Create character
// @Description  Creates a new character for the authenticated user.
// @Tags         characters
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        characterCreateRequest  body      dto.CharacterCreateRequest  true  "Character Create Request"
// @Success      201  {object}  dto.APIObjectResponse{data=string}  "Character created successfully"
// @Failure      400  {object}  dto.APIErrorResponse{data=interface{}}  "Invalid request"
// @Failure      401  {object}  dto.APIErrorResponse{data=interface{}}  "Unauthorized"
// @Router       /characters [post]
func (h *CharacterHandler) Create(c echo.Context) error {
	var req dto.CharacterCreateRequest
	if err := c.Bind(&req); err != nil {
		e := custom.NewBadRequestError("invalid payload")
		custom.PanicException(e)
	}
	if req.Privacy == "" {
		req.Privacy = model.PrivacyPublic
	}
	if err := h.v.Struct(req); err != nil {
		e := custom.NewValidationError("required fields are missing or invalid")
		custom.PanicException(e)
	}
	uid, _ := middleware.GetUserID(c)
	_, err := h.uc.Create(uid, &dto.CreateCharacterInput{
		Title: req.Title, Description: req.Description, ClassID: req.ClassID,
		RaceID: req.RaceID, Privacy: req.Privacy,
	})
	if err != nil {
		custom.PanicException(err)
	}
	return c.JSON(http.StatusCreated, custom.BuildResponse(custom.Success, "character created"))
}

// UpdateCharacter godoc
// @Summary      Update character
// @Description  Updates an existing character for the authenticated user.
// @Tags         characters
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id                       path      string                      true  "Character ID"
// @Param        characterUpdateRequest   body      dto.CharacterUpdateRequest  true  "Character Update Request"
// @Success      200  {object}  dto.APIObjectResponse{data=string}  "Character updated successfully"
// @Failure      400  {object}  dto.APIErrorResponse{data=interface{}}  "Invalid request"
// @Failure      401  {object}  dto.APIErrorResponse{data=interface{}}  "Unauthorized"
// @Failure      404  {object}  dto.APIErrorResponse{data=interface{}}  "Character not found"
// @Router       /characters/{id} [put]
func (h *CharacterHandler) Update(c echo.Context) error {
	defer custom.PanicController(c)
	id := c.Param("id")
	var req dto.CharacterUpdateRequest
	if err := c.Bind(&req); err != nil {
		e := custom.NewBadRequestError("invalid payload")
		custom.PanicException(e)
	}
	if err := h.v.Struct(req); err != nil {
		e := custom.NewValidationError("required fields are missing or invalid")
		custom.PanicException(e)
	}
	uid, _ := middleware.GetUserID(c)
	err := h.uc.Update(uid, id, &dto.UpdateCharacterInput{
		Title: req.Title, Description: req.Description, ClassID: req.ClassID,
		RaceID: req.RaceID, Privacy: req.Privacy,
	})
	if err != nil {
		custom.PanicException(err)
	}
	return c.JSON(http.StatusOK, custom.BuildResponse(custom.Success, "character updated"))
}

// DeleteCharacter godoc
// @Summary      Delete character
// @Description  Deletes a character for the authenticated user.
// @Tags         characters
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Character ID"
// @Success      200  {object}  dto.APIObjectResponse{data=string}  "Character deleted successfully"
// @Failure      401  {object}  dto.APIErrorResponse{data=interface{}}  "Unauthorized"
// @Failure      404  {object}  dto.APIErrorResponse{data=interface{}}  "Character not found"
// @Router       /characters/{id} [delete]
func (h *CharacterHandler) Delete(c echo.Context) error {
	defer custom.PanicController(c)
	id := c.Param("id")
	uid, _ := middleware.GetUserID(c)
	if err := h.uc.Delete(uid, id); err != nil {
		custom.PanicException(err)
	}
	return c.JSON(http.StatusOK, custom.BuildResponse(custom.Success, "character deleted"))
}
