package handlers

import (
	"dungeons-dragon-service/internal/dto"
	"dungeons-dragon-service/internal/http/custom"
	middleware "dungeons-dragon-service/internal/http/middlewares"
	usecase "dungeons-dragon-service/internal/usecases"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type QuestHandler struct {
	uc usecase.QuestUseCase
	v  *validator.Validate
}

func NewQuestHandler(uc usecase.QuestUseCase) *QuestHandler {
	return &QuestHandler{uc: uc, v: validator.New()}
}

// List godoc
// @Summary      List quests
// @Description  Retrieves a list of quests for the authenticated user.
// @Tags         quests
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.APIObjectResponse{data=[]dto.QuestResponse}  "List of quests"
// @Failure      400  {object}  dto.APIErrorResponse{data=interface{}} "Invalid request"
// @Failure      401  {object}  dto.APIErrorResponse{data=interface{}} "Unauthorized"
// @Router       /quests [get]
func (h *QuestHandler) List(c echo.Context) error {
	defer custom.PanicController(c)
	auth := middleware.IsAuthenticated(c)
	list, err := h.uc.ListForUser(auth)
	if err != nil {
		custom.PanicException(err)
	}

	return c.JSON(http.StatusOK, custom.BuildResponse(custom.Success, list))
}

// Create godoc
// @Summary      Create quest
// @Description  Creates a new quest for the authenticated user.
// @Tags         quests
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        quest  body      dto.QuestCreateRequest  true  "Quest creation payload"
// @Success      201    {object}  dto.APIObjectResponse{data=string}  "Quest created successfully"
// @Failure      400    {object}  dto.APIErrorResponse{data=interface{}}  "Invalid request"
// @Failure      401    {object}  dto.APIErrorResponse{data=interface{}}  "Unauthorized"
// @Router       /quests [post]
func (h *QuestHandler) Create(c echo.Context) error {
	defer custom.PanicController(c)
	var req dto.QuestCreateRequest
	if err := c.Bind(&req); err != nil {
		e := custom.NewBadRequestError("invalid payload")
		custom.PanicException(e)
	}
	if err := h.v.Struct(req); err != nil {
		e := custom.NewValidationError("required fields are missing or invalid")
		custom.PanicException(e)
	}
	uid, _ := middleware.GetUserID(c)
	err := h.uc.Create(uid, &dto.CreateQuestInput{
		Title:        req.Title,
		Description:  req.Description,
		QuestLevelID: req.QuestLevelID,
		Privacy:      req.Privacy,
	})
	if err != nil {
		custom.PanicException(err)
	}
	return c.JSON(http.StatusCreated, custom.BuildResponse(custom.Success, "quest created"))
}

// Update godoc
// @Summary      Update quest
// @Description  Updates an existing quest for the authenticated user.
// @Tags         quests
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id     path      string                  true  "Quest ID"
// @Param        quest  body      dto.QuestUpdateRequest  true  "Quest update payload"
// @Success      200    {object}  dto.APIObjectResponse{data=string}  "Quest updated successfully"
// @Failure      400    {object}  dto.APIErrorResponse{data=interface{}}  "Invalid request"
// @Failure      401    {object}  dto.APIErrorResponse{data=interface{}}  "Unauthorized"
// @Router       /quests/{id} [put]
func (h *QuestHandler) Update(c echo.Context) error {
	defer custom.PanicController(c)
	id := c.Param("id")
	var req dto.QuestUpdateRequest
	if err := c.Bind(&req); err != nil {
		e := custom.NewBadRequestError("invalid payload")
		custom.PanicException(e)
	}
	if err := h.v.Struct(req); err != nil {
		e := custom.NewValidationError("required fields are missing or invalid")
		custom.PanicException(e)
	}
	uid, _ := middleware.GetUserID(c)
	err := h.uc.Update(uid, id, &dto.UpdateQuestInput{
		Title:        req.Title,
		Description:  req.Description,
		QuestLevelID: req.QuestLevelID,
		Privacy:      req.Privacy,
	})
	if err != nil {
		custom.PanicException(err)
	}
	return c.JSON(http.StatusOK, custom.BuildResponse(custom.Success, "quest updated"))
}

// Delete godoc
// @Summary      Delete quest
// @Description  Deletes a quest for the authenticated user.
// @Tags         quests
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Quest ID"
// @Success      200  {object}  dto.APIObjectResponse{data=string}  "Quest deleted successfully"
// @Failure      400  {object}  dto.APIErrorResponse{data=interface{}}  "Invalid request"
// @Failure      401  {object}  dto.APIErrorResponse{data=interface{}}  "Unauthorized"
// @Router       /quests/{id} [delete]
func (h *QuestHandler) Delete(c echo.Context) error {
	defer custom.PanicController(c)
	id := c.Param("id")
	uid, _ := middleware.GetUserID(c)
	if err := h.uc.Delete(uid, id); err != nil {
		custom.PanicException(err)
	}
	return c.JSON(http.StatusOK, custom.BuildResponse(custom.Success, "quest deleted"))
}
