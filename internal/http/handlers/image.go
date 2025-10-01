package handlers

import (
	"dungeons-dragon-service/internal/config"
	"dungeons-dragon-service/internal/http/custom"
	"dungeons-dragon-service/internal/http/middlewares"
	usecase "dungeons-dragon-service/internal/usecases"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

type ImageHandler struct {
	uc usecase.ImageUseCase
}

func NewImageHandler(uc usecase.ImageUseCase) *ImageHandler {
	return &ImageHandler{uc: uc}
}

// UploadCharacterImage godoc
// @Summary      Upload character image
// @Description  Uploads an image for a specific character.
// @Tags         images
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        id      path      string                true  "Character ID"
// @Param        images  formData  file   true  "List of character images (can upload multiple)"
// @Success      200     {object}  dto.APIObjectResponse{data=string}  "Character images uploaded successfully"
// @Failure      400     {object}  dto.APIErrorResponse{data=interface{}}  "Invalid request"
// @Failure      401     {object}  dto.APIErrorResponse{data=interface{}}  "Unauthorized"
// @Router       /characters/{id}/images [post]
func (h *ImageHandler) UploadCharacterImage(c echo.Context) error {
	defer custom.PanicController(c)

	id := c.Param("id")
	if id == "" {
		custom.PanicException(custom.NewBadRequestError("character id is required"))
	}
	uid, _ := middlewares.GetUserID(c)
	// Parse multipart form with a reasonable max memory
	form, err := c.MultipartForm()
	if err != nil {
		custom.PanicException(custom.NewBadRequestError("invalid form data"))
	}

	images := form.File["images"]
	//upload to usecase
	if err := h.uc.UploadCharacterImage(uid, id, images); err != nil {
		custom.PanicException(err)
	}
	return c.JSON(200, custom.BuildResponse(custom.Success, "character images uploaded successfully"))

}

// UploadQuestImage godoc
// @Summary      Upload quest image
// @Description  Uploads an image for a specific quest.
// @Tags         images
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        id      path      string                true  "Quest ID"
// @Param        images  formData  file   true  "List of character images (can upload multiple)"
// @Success      200     {object}  dto.APIObjectResponse{data=string}  "Quest images uploaded successfully"
// @Failure      400     {object}  dto.APIErrorResponse{data=interface{}}  "Invalid request"
// @Failure      401     {object}  dto.APIErrorResponse{data=interface{}}  "Unauthorized"
// @Router       /quests/{id}/images [post]
func (h *ImageHandler) UploadQuestImage(c echo.Context) error {
	defer custom.PanicController(c)
	id := c.Param("id")
	if id == "" {
		custom.PanicException(custom.NewBadRequestError("quest id is required"))
	}
	uid, _ := middlewares.GetUserID(c)

	// Parse multipart form with a reasonable max memory
	form, err := c.MultipartForm()
	if err != nil {
		custom.PanicException(custom.NewBadRequestError("invalid form data"))
	}

	images := form.File["images"]
	//upload to usecase
	if err := h.uc.UploadQuestImage(uid, id, images); err != nil {
		custom.PanicException(err)
	}

	return c.JSON(200, custom.BuildResponse(custom.Success, "quest images uploaded successfully"))
}

func (h *ImageHandler) GetImage(c echo.Context) error {
	defer custom.PanicController(c)
	filename := c.Param("filename")
	if filename == "" {
		custom.PanicException(custom.NewBadRequestError("filename is required"))
	}
	filePath := filepath.Join(config.GetConfigString("FILE_STORAGE_PATH"), filename)
	return c.File(filePath)
}
