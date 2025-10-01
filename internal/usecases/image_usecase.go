package usecases

import (
	"dungeons-dragon-service/internal/config"
	"dungeons-dragon-service/internal/domain/model"
	"dungeons-dragon-service/internal/domain/repository"
	"dungeons-dragon-service/internal/helper"
	"dungeons-dragon-service/internal/http/custom"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"path/filepath"
	"time"

	"gorm.io/datatypes"
)

type ImageUseCase interface {
	UploadCharacterImage(uid string, characterID string, images []*multipart.FileHeader) error
	UploadQuestImage(uid string, questID string, images []*multipart.FileHeader) error
}

type imageUseCase struct {
	images     repository.ImageRepository
	characters repository.CharacterRepository
	quests     repository.QuestRepository
}

func NewImageUsecase(images repository.ImageRepository, characters repository.CharacterRepository, quests repository.QuestRepository) ImageUseCase {
	return &imageUseCase{images: images, characters: characters, quests: quests}
}

func (u *imageUseCase) UploadCharacterImage(userID string, characterID string, images []*multipart.FileHeader) error {
	character, err := u.characters.FindByID(characterID)
	if err != nil {
		return err
	}
	if character == nil {
		return custom.NewNotFoundError("character not found")
	}
	if character.UserID.String() != userID {
		return custom.NewUnauthorizedError("you do not have permission to upload images for this character")
	}
	if character.Status == model.ItemStatusArchived {
		return custom.NewBadRequestError("cannot upload images to an archived character")
	}

	paths := make([]string, 0, len(images))
	if err := helper.ValidateImages(len(images)); err == nil {
		for _, img := range images {
			imageName := fmt.Sprintf("%d-%s", time.Now().UnixNano()%1_000_000, filepath.Base(img.Filename))
			imagePath := filepath.Join(config.GetConfigString("FILE_STORAGE_PATH"), imageName)

			if err := helper.SaveUploadedFile(img, imagePath); err != nil {
				custom.PanicException(custom.NewUnexpectedError("failed to save image"))
			}
			paths = append(paths, imagePath)
		}
	} else {
		return err
	}

	characterImages, err := u.images.GetCharacterImageByID(characterID)
	if err != nil {
		return err
	}
	if characterImages != nil {
		//loop for remove old images
		// Delete existing image
		for _, img := range characterImages {
			if err := helper.DeleteFileIfExists(img.Path); err != nil {
				log.Println("failed to delete old character image:", err)
			}
		}
		if err := u.images.DeleteCharacterImageByID(characterID); err != nil {
			return custom.NewUnexpectedError("failed to delete existing character image")
		}
	}

	// Create new image
	var imagePaths []string
	for _, path := range paths {
		charImg, err := u.images.CreateCharacterImage(&model.CharacterImage{
			CharacterID: character.ID,
			Path:        path,
		})
		if err != nil {
			return custom.NewUnexpectedError("failed to create character image")
		}
		imagePaths = append(imagePaths, charImg.Path)
	}
	// Update character images as string slice
	imageBytes, err := json.Marshal(imagePaths)
	if err != nil {
		return custom.NewUnexpectedError("failed to marshal character images")
	}
	character.ImagePath = datatypes.JSON(imageBytes)
	if _, err := u.characters.Update(character); err != nil {
		return custom.NewUnexpectedError("failed to update character images")
	}

	return nil
}

func (u *imageUseCase) UploadQuestImage(userID string, questID string, images []*multipart.FileHeader) error {
	quest, err := u.quests.FindByID(questID)
	if err != nil {
		return err
	}
	if quest == nil {
		return custom.NewNotFoundError("quest not found")
	}
	if quest.UserID.String() != userID {
		return custom.NewUnauthorizedError("you do not have permission to upload images for this quest")
	}
	if quest.Status == model.ItemStatusArchived {
		return custom.NewBadRequestError("cannot upload images to an archived quest")
	}
	paths := make([]string, 0, len(images))
	if err := helper.ValidateImages(len(images)); err == nil {
		for _, img := range images {
			imageName := fmt.Sprintf("%d-%s", time.Now().UnixNano()%1_000_000, filepath.Base(img.Filename))
			imagePath := filepath.Join(config.GetConfigString("FILE_STORAGE_PATH"), imageName)

			if err := helper.SaveUploadedFile(img, imagePath); err != nil {
				custom.PanicException(custom.NewUnexpectedError("failed to save image"))
			}
			paths = append(paths, imagePath)
		}
	} else {
		return err
	}

	questImages, err := u.images.GetQuestImageByID(questID)
	if err != nil {
		return err
	}
	if questImages != nil {
		//loop for remove old images
		for _, img := range questImages {
			if err := helper.DeleteFileIfExists(img.Path); err != nil {
				log.Println("failed to delete old quest image:", err)
			}
		}
		// Delete existing image
		if err := u.images.DeleteQuestImageByID(questID); err != nil {
			return custom.NewUnexpectedError("failed to delete existing quest image")
		}
	}
	// Create new image
	var imagePaths []string
	for _, path := range paths {
		questImg, err := u.images.CreateQuestImage(&model.QuestImage{
			QuestID: quest.ID,
			Path:    path,
		})
		if err != nil {
			return custom.NewUnexpectedError("failed to create quest image")
		}
		imagePaths = append(imagePaths, questImg.Path)
	}
	// Update quest images as string slice
	imageBytes, err := json.Marshal(imagePaths)
	if err != nil {
		return custom.NewUnexpectedError("failed to marshal quest images")
	}
	quest.ImagePath = datatypes.JSON(imageBytes)
	if _, err := u.quests.Update(quest); err != nil {
		return custom.NewUnexpectedError("failed to update quest images")
	}

	return nil
}
