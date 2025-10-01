package dto

import (
	"dungeons-dragon-service/internal/domain/model"
)

type CharacterResponse struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	UserID      string        `json:"user_id"`
	ClassID     string        `json:"class_id"`
	RaceID      string        `json:"race_id"`
	Privacy     model.Privacy `json:"privacy"`
	Status      string        `json:"status"`
	Images      []string      `json:"images"`
}

type CreateCharacterInput struct {
	Title       string        `json:"title"`
	Description string        `json:"description"`
	ClassID     string        `json:"class_id"`
	RaceID      string        `json:"race_id"`
	Privacy     model.Privacy `json:"privacy"`
}

type UpdateCharacterInput struct {
	Title       *string
	Description *string
	ClassID     *string
	RaceID      *string
	Privacy     *model.Privacy
}

type CharacterCreateRequest struct {
	Title       string        `json:"title" validate:"required,max=200"`
	Description string        `json:"description" validate:"required"`
	ClassID     string        `json:"class_id" validate:"required"`
	RaceID      string        `json:"race_id" validate:"required"`
	Privacy     model.Privacy `json:"privacy" validate:"oneof=public private"`
}
type CharacterUpdateRequest struct {
	Title       *string        `json:"title" validate:"omitempty,max=200"`
	Description *string        `json:"description"`
	ClassID     *string        `json:"class_id"`
	RaceID      *string        `json:"race_id"`
	Privacy     *model.Privacy `json:"privacy" validate:"omitempty,oneof=public private"`
}
