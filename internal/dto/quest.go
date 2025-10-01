package dto

import (
	"dungeons-dragon-service/internal/domain/model"
)

type QuestResponse struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	QuestLevel  string        `json:"quest_level"`
	Privacy     model.Privacy `json:"privacy"`
	Status      string        `json:"status"`
	Images      []string      `json:"images"`
}

type CreateQuestInput struct {
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	QuestLevelID string        `json:"quest_level_id"`
	Privacy      model.Privacy `json:"privacy"`
}

type UpdateQuestInput struct {
	Title        *string        `json:"title"`
	Description  *string        `json:"description"`
	QuestLevelID *string        `json:"quest_level_id"`
	Privacy      *model.Privacy `json:"privacy"`
}

type QuestCreateRequest struct {
	Title        string        `json:"title" validate:"required,max=200"`
	Description  string        `json:"description" validate:"required"`
	QuestLevelID string        `json:"quest_level_id" validate:"required"`
	Privacy      model.Privacy `json:"privacy" validate:"oneof=public private"`
}
type QuestUpdateRequest struct {
	Title        *string        `json:"title" validate:"omitempty,max=200"`
	Description  *string        `json:"description"`
	QuestLevelID *string        `json:"quest_level_id"`
	Privacy      *model.Privacy `json:"privacy" validate:"omitempty,oneof=public private"`
}
