package usecases

import (
	"dungeons-dragon-service/internal/domain/model"
	"dungeons-dragon-service/internal/domain/repository"
	"dungeons-dragon-service/internal/dto"
	"dungeons-dragon-service/internal/helper"
	"dungeons-dragon-service/internal/http/custom"
	"encoding/json"
)

type CharacterUseCase interface {
	ListPublic() ([]dto.CharacterResponse, error)
	ListForUser(authenticated bool) ([]dto.CharacterResponse, error)
	Create(userID string, in *dto.CreateCharacterInput) (*dto.CharacterResponse, error)
	Update(userID string, id string, in *dto.UpdateCharacterInput) error
	Delete(userID string, id string) error
}

type characterUseCase struct {
	characters repository.CharacterRepository
	classes    repository.ClassRepository
	races      repository.RaceRepository
}

func NewCharacterUsecase(c repository.CharacterRepository, cl repository.ClassRepository, r repository.RaceRepository) CharacterUseCase {
	return &characterUseCase{characters: c, classes: cl, races: r}
}

func ResponseCharacters(c []model.Character) []dto.CharacterResponse {
	res := make([]dto.CharacterResponse, len(c))
	for i, char := range c {
		//unmarshal images
		images := []string{}
		urls := []string{}
		if err := json.Unmarshal(char.ImagePath, &images); err == nil {
			for _, img := range images {
				url := helper.GetImageURL(img)
				urls = append(urls, url)
			}
		}
		res[i] = dto.CharacterResponse{
			ID:          char.ID.String(),
			Title:       char.Title,
			Description: char.Description,
			ClassID:     char.ClassID.String(),
			RaceID:      char.RaceID.String(),
			UserID:      char.UserID.String(),
			Privacy:     char.Privacy,
			Status:      string(char.Status),
			Images:      urls,
		}
	}
	return res
}

func (u *characterUseCase) ListPublic() ([]dto.CharacterResponse, error) {
	list, err := u.characters.ListPublic()
	if err != nil {
		return nil, custom.NewUnexpectedError("failed to list characters")
	}
	return ResponseCharacters(list), nil
}

func (u *characterUseCase) ListForUser(authenticated bool) ([]dto.CharacterResponse, error) {
	if authenticated {
		list, err := u.characters.ListAll()
		if err != nil {
			return nil, custom.NewUnexpectedError("failed to list characters")
		}
		return ResponseCharacters(list), nil
	}
	list, err := u.characters.ListPublic()
	if err != nil {
		return nil, custom.NewUnexpectedError("failed to list characters")
	}
	return ResponseCharacters(list), nil
}

func (u *characterUseCase) Create(userID string, in *dto.CreateCharacterInput) (*dto.CharacterResponse, error) {
	// Validate description and images
	if err := helper.ValidateDescription(in.Description); err != nil {
		return nil, custom.NewBadRequestError("invalid description")
	}
	// Validate class & race existence
	if _, err := u.classes.FindByID(in.ClassID); err != nil {
		return nil, custom.NewNotFoundError("class not found")
	}
	if _, err := u.races.FindByID(in.RaceID); err != nil {
		return nil, custom.NewNotFoundError("race not found")
	}

	// imgJSON, _ := json.Marshal(in.Images)
	m := &model.Character{
		UserID:      helper.ParseUUIDOrNil(userID),
		Title:       in.Title,
		Description: in.Description,
		ClassID:     helper.ParseUUIDOrNil(in.ClassID),
		RaceID:      helper.ParseUUIDOrNil(in.RaceID),
		Privacy:     in.Privacy,
		// Images:      []byte("[]"),
		Status: model.ItemStatusActive,
	}
	if _, err := u.characters.Create(m); err != nil {
		return nil, custom.NewUnexpectedError("failed to create character")
	}
	response := ResponseCharacters([]model.Character{*m})[0]
	return &response, nil
}

func (u *characterUseCase) Update(userID string, id string, in *dto.UpdateCharacterInput) error {
	m, err := u.characters.FindByID(id)
	if err != nil {
		return custom.NewNotFoundError("character not found")
	}
	if m.UserID != helper.ParseUUIDOrNil(userID) {
		return custom.NewForbiddenError("forbidden")
	}
	if m.Status == model.ItemStatusArchived {
		return custom.NewBadRequestError("cannot modify archived")
	}

	if in.Title != nil {
		m.Title = *in.Title
	}
	if in.Description != nil {
		if err := helper.ValidateDescription(*in.Description); err != nil {
			return custom.NewBadRequestError("invalid description")
		}
		m.Description = *in.Description
	}
	if in.ClassID != nil {
		if _, err := u.classes.FindByID(*in.ClassID); err != nil {
			return custom.NewNotFoundError("class not found")
		}
		m.ClassID = helper.ParseUUIDOrNil(*in.ClassID)
	}
	if in.RaceID != nil {
		if _, err := u.races.FindByID(*in.RaceID); err != nil {
			return custom.NewNotFoundError("race not found")
		}
		m.RaceID = helper.ParseUUIDOrNil(*in.RaceID)
	}
	if in.Privacy != nil {
		m.Privacy = *in.Privacy
	}
	if _, err := u.characters.Update(m); err != nil {
		return custom.NewUnexpectedError("failed to update character")
	}
	return nil
}

func (u *characterUseCase) Delete(userID string, id string) error {
	m, err := u.characters.FindByID(id)
	if err != nil {
		return custom.NewNotFoundError("character not found")
	}
	if m.UserID != helper.ParseUUIDOrNil(userID) {
		return custom.NewForbiddenError("forbidden")
	}
	return u.characters.Delete(id)
}
