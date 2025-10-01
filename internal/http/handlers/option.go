package handlers

import (
	"dungeons-dragon-service/internal/dto"
	"dungeons-dragon-service/internal/http/custom"
	usecase "dungeons-dragon-service/internal/usecases"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type OptionHandler struct {
	uc usecase.OptionUseCase
	v  *validator.Validate
}

func NewOptionHandler(optUseCase usecase.OptionUseCase) *OptionHandler {
	return &OptionHandler{
		uc: optUseCase,
		v:  validator.New(),
	}
}

// ListClasses godoc
// @Summary      List all classes
// @Description  Retrieves a list of all available classes.
// @Tags         options
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.APIObjectResponse{data=[]dto.ClassResponse}  "List of classes"
// @Failure      400  {object}  dto.APIErrorResponse{data=interface{}}  "Bad Request"
// @Failure      401  {object}  dto.APIErrorResponse{data=interface{}}  "Unauthorized"
// @Router       /options/classes [get]
func (h *OptionHandler) ListClasses(c echo.Context) error {
	defer custom.PanicController(c)
	res, err := h.uc.ListClasses()
	if err != nil {
		custom.PanicException(err)
	}
	return c.JSON(http.StatusOK, custom.BuildResponse(custom.Success, res))
}

// ListRaces godoc
// @Summary	  List all races
// @Description  Retrieves a list of all available races.
// @Tags         options
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.APIObjectResponse{data=[]dto.RaceResponse}  "List of races"
// @Failure      400  {object}  dto.APIErrorResponse{data=interface{}}  "Bad Request"
// @Failure      401  {object}  dto.APIErrorResponse{data=interface{}}  "Unauthorized"
// @Router       /options/races [get]
func (h *OptionHandler) ListRaces(c echo.Context) error {
	defer custom.PanicController(c)
	res, err := h.uc.ListRaces()
	if err != nil {
		custom.PanicException(err)
	}
	return c.JSON(http.StatusOK, custom.BuildResponse(custom.Success, res))
}

// ListQuestLevels godoc
// @Summary	  List all quest levels
// @Description  Retrieves a list of all available quest levels.
// @Tags         options
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.APIObjectResponse{data=[]dto.QuestLevelResponse}  "List of quest levels"
// @Failure      400  {object}  dto.APIErrorResponse{data=interface{}}  "Bad Request"
// @Failure      401  {object}  dto.APIErrorResponse{data=interface{}}  "Unauthorized"
// @Router       /options/quest-levels [get]
func (h *OptionHandler) ListQuestLevels(c echo.Context) error {
	defer custom.PanicController(c)
	res, err := h.uc.ListQuestLevels()
	if err != nil {
		custom.PanicException(err)
	}
	return c.JSON(http.StatusOK, custom.BuildResponse(custom.Success, res))
}

// CreateClass godoc
// @Summary      Create a new class
// @Description  Creates a new class with the provided name.
// @Tags         options
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        class  body      dto.OptionReq  true  "Class to create"
// @Success      201    {object}  dto.APIObjectResponse{data=string}  "Class created successfully"
// @Failure      400    {object}  dto.APIErrorResponse{data=interface{}}  "Bad Request"
// @Failure      401    {object}  dto.APIErrorResponse{data=interface{}}  "Unauthorized"
// @Router       /options/classes [post]
func (h *OptionHandler) CreateClass(c echo.Context) error {
	defer custom.PanicController(c)
	var req dto.OptionReq
	if err := c.Bind(&req); err != nil {
		e := custom.NewBadRequestError("invalid payload")
		custom.PanicException(e)
	}
	if err := h.v.Struct(req); err != nil {
		e := custom.NewValidationError("required fields are missing or invalid")
		custom.PanicException(e)
	}
	err := h.uc.CreateClass(req.Name)
	if err != nil {
		custom.PanicException(err)
	}
	return c.JSON(http.StatusCreated, custom.BuildResponse(custom.Success, "class created"))
}

// UpdateClass godoc
// @Summary      Update an existing class
// @Description  Updates the name of an existing class identified by its ID.
// @Tags         options
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id     path      string         true  "Class ID"
// @Param        class  body      dto.OptionReq  true  "Updated class data"
// @Success      200    {object}  dto.APIObjectResponse{data=string}  "Class updated successfully"
// @Failure      400    {object}  dto.APIErrorResponse{data=interface{}}  "Bad Request"
// @Failure      401    {object}  dto.APIErrorResponse{data=interface{}}  "Unauthorized"
// @Failure      404    {object}  dto.APIErrorResponse{data=interface{}}  "Class not found"
// @Router       /options/classes/{id} [put]
func (h *OptionHandler) UpdateClass(c echo.Context) error {
	defer custom.PanicController(c)
	id := c.Param("id")
	var req dto.OptionReq
	if err := c.Bind(&req); err != nil {
		e := custom.NewBadRequestError("invalid payload")
		custom.PanicException(e)
	}
	if err := h.v.Struct(req); err != nil {
		e := custom.NewValidationError("required fields are missing or invalid")
		custom.PanicException(e)
	}
	err := h.uc.UpdateClass(id, req.Name)
	if err != nil {
		custom.PanicException(err)
	}
	return c.JSON(http.StatusOK, custom.BuildResponse(custom.Success, "class updated"))
}

func (h *OptionHandler) DeleteClass(c echo.Context) error {
	defer custom.PanicController(c)
	id := c.Param("id")
	if err := h.uc.DeleteClass(id); err != nil {
		custom.PanicException(err)
	}
	return c.JSON(http.StatusOK, custom.BuildResponse(custom.Success, "class deleted"))
}

// CreateRace godoc
// @Summary      Create a new race
// @Description  Creates a new race with the provided name.
// @Tags         options
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        race  body      dto.OptionReq  true  "Race to create"
// @Success      201    {object}  dto.APIObjectResponse{data=string}  "Race created successfully"
// @Failure      400    {object}  dto.APIErrorResponse{data=interface{}}  "Bad Request"
// @Failure      401    {object}  dto.APIErrorResponse{data=interface{}}  "Unauthorized"
// @Router       /options/races [post]
func (h *OptionHandler) CreateRace(c echo.Context) error {
	defer custom.PanicController(c)
	var req dto.OptionReq
	if err := c.Bind(&req); err != nil {
		e := custom.NewBadRequestError("invalid payload")
		custom.PanicException(e)
	}
	if err := h.v.Struct(req); err != nil {
		e := custom.NewValidationError("required fields are missing or invalid")
		custom.PanicException(e)
	}
	err := h.uc.CreateRace(req.Name)
	if err != nil {
		custom.PanicException(err)
	}
	return c.JSON(http.StatusCreated, custom.BuildResponse(custom.Success, "race created"))
}

// UpdateRace godoc
// @Summary      Update an existing race
// @Description  Updates the name of an existing race identified by its ID.
// @Tags         options
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id     path      string         true  "Race ID"
// @Param        race   body      dto.OptionReq  true  "Updated race data"
// @Success      200    {object}  dto.APIObjectResponse{data=string}  "Race updated successfully"
// @Failure      400    {object}  dto.APIErrorResponse{data=interface{}}  "Bad Request"
// @Failure      401    {object}  dto.APIErrorResponse{data=interface{}}  "Unauthorized"
// @Failure      404    {object}  dto.APIErrorResponse{data=interface{}}  "Not Found"
// @Router       /options/races/{id} [put]
func (h *OptionHandler) UpdateRace(c echo.Context) error {
	defer custom.PanicController(c)
	id := c.Param("id")
	var req dto.OptionReq
	if err := c.Bind(&req); err != nil {
		e := custom.NewBadRequestError("invalid payload")
		custom.PanicException(e)
	}
	if err := h.v.Struct(req); err != nil {
		e := custom.NewValidationError("required fields are missing or invalid")
		custom.PanicException(e)
	}
	err := h.uc.UpdateRace(id, req.Name)
	if err != nil {
		custom.PanicException(err)
	}
	return c.JSON(http.StatusOK, custom.BuildResponse(custom.Success, "race updated"))
}

// DeleteRace godoc
// @Summary      Delete an existing race
// @Description  Deletes an existing race identified by its ID.
// @Tags         options
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id     path      string         true  "Race ID"
// @Success      200    {object}  dto.APIObjectResponse{data=string}  "Race deleted successfully"
// @Failure      400    {object}  dto.APIErrorResponse{data=interface{}}  "Bad Request"
// @Failure      401    {object}  dto.APIErrorResponse{data=interface{}}  "Unauthorized"
// @Failure      404    {object}  dto.APIErrorResponse{data=interface{}}  "Not Found"
// @Router       /options/races/{id} [delete]
func (h *OptionHandler) DeleteRace(c echo.Context) error {
	defer custom.PanicController(c)
	id := c.Param("id")
	if err := h.uc.DeleteRace(id); err != nil {
		custom.PanicException(err)
	}
	return c.JSON(http.StatusOK, custom.BuildResponse(custom.Success, "race deleted"))
}

// CreateQuestLevel godoc
// @Summary      Create a new quest level
// @Description  Creates a new quest level with the provided name.
// @Tags         options
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        questLevel  body      dto.OptionReq  true  "Quest level to create"
// @Success      201    {object}  dto.APIObjectResponse{data=string}  "Quest level created successfully"
// @Failure      400    {object}  dto.APIErrorResponse{data=interface{}}  "Bad Request"
// @Failure      401    {object}  dto.APIErrorResponse{data=interface{}}  "Unauthorized"
// @Router       /options/quest-levels [post]
func (h *OptionHandler) CreateQuestLevel(c echo.Context) error {
	defer custom.PanicController(c)
	var req dto.OptionReq
	if err := c.Bind(&req); err != nil {
		e := custom.NewBadRequestError("invalid payload")
		custom.PanicException(e)
	}
	if err := h.v.Struct(req); err != nil {
		e := custom.NewValidationError("required fields are missing or invalid")
		custom.PanicException(e)
	}
	err := h.uc.CreateQuestLevel(req.Name)
	if err != nil {
		custom.PanicException(err)
	}
	return c.JSON(http.StatusCreated, custom.BuildResponse(custom.Success, "quest level created"))
}

// UpdateQuestLevel godoc
// @Summary      Update an existing quest level
// @Description  Updates the name of an existing quest level identified by its ID.
// @Tags         options
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id          path      string         true  "Quest Level ID"
// @Param        questLevel  body      dto.OptionReq  true  "Updated quest level data"
// @Success      200    {object}  dto.APIObjectResponse{data=string}  "Quest level updated successfully"
// @Failure      400    {object}  dto.APIErrorResponse{data=interface{}}  "Bad Request"
// @Failure      401    {object}  dto.APIErrorResponse{data=interface{}}  "Unauthorized"
// @Failure      404    {object}  dto.APIErrorResponse{data=interface{}}  "Not Found"
// @Router       /options/quest-levels/{id} [put]
func (h *OptionHandler) UpdateQuestLevel(c echo.Context) error {
	defer custom.PanicController(c)
	id := c.Param("id")
	var req dto.OptionReq
	if err := c.Bind(&req); err != nil {
		e := custom.NewBadRequestError("invalid payload")
		custom.PanicException(e)
	}
	if err := h.v.Struct(req); err != nil {
		e := custom.NewValidationError("required fields are missing or invalid")
		custom.PanicException(e)
	}
	err := h.uc.UpdateQuestLevel(id, req.Name)
	if err != nil {
		custom.PanicException(err)
	}
	return c.JSON(http.StatusOK, custom.BuildResponse(custom.Success, "quest level updated"))
}

// DeleteQuestLevel godoc
// @Summary      Delete an existing quest level
// @Description  Deletes an existing quest level identified by its ID.
// @Tags         options
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id          path      string         true  "Quest Level ID"
// @Success      200    {object}  dto.APIObjectResponse{data=string}  "Quest level deleted successfully"
// @Failure      400    {object}  dto.APIErrorResponse{data=interface{}}  "Bad Request"
// @Failure      401    {object}  dto.APIErrorResponse{data=interface{}}  "Unauthorized"
// @Failure      404    {object}  dto.APIErrorResponse{data=interface{}}  "Not Found"
// @Router       /options/quest-levels/{id} [delete]
func (h *OptionHandler) DeleteQuestLevel(c echo.Context) error {
	defer custom.PanicController(c)
	id := c.Param("id")
	if err := h.uc.DeleteQuestLevel(id); err != nil {
		custom.PanicException(err)
	}
	return c.JSON(http.StatusOK, custom.BuildResponse(custom.Success, "quest level deleted"))
}
