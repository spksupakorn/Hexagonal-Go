package usecases

import (
	"dungeons-dragon-service/internal/domain/model"
	"dungeons-dragon-service/internal/domain/repository"
	"dungeons-dragon-service/internal/dto"
	"dungeons-dragon-service/internal/helper"
	"dungeons-dragon-service/internal/http/custom"
	"encoding/json"
)

type QuestUseCase interface {
	ListPublic() ([]dto.QuestResponse, error)
	ListForUser(authenticated bool) ([]dto.QuestResponse, error)
	Create(userID string, in *dto.CreateQuestInput) error
	Update(userID string, id string, in *dto.UpdateQuestInput) error
	Delete(userID string, id string) error
}

type questUseCase struct {
	quests      repository.QuestRepository
	questLevels repository.QuestLevelRepository
}

func NewQuestUsecase(q repository.QuestRepository, ql repository.QuestLevelRepository) QuestUseCase {
	return &questUseCase{quests: q, questLevels: ql}
}

func ResponseQuests(q []model.Quest) []dto.QuestResponse {
	res := make([]dto.QuestResponse, len(q))
	for i, quest := range q {
		//unmarshal images
		images := []string{}
		urls := []string{}
		if err := json.Unmarshal(quest.ImagePath, &images); err == nil {
			for _, img := range images {
				url := helper.GetImageURL(img)
				urls = append(urls, url)
			}
		}
		res[i] = dto.QuestResponse{
			ID:          quest.ID.String(),
			Title:       quest.Title,
			Description: quest.Description,
			Privacy:     quest.Privacy,
			Status:      string(quest.Status),
			Images:      urls,
		}
	}
	return res
}

func (u *questUseCase) ListPublic() ([]dto.QuestResponse, error) {
	list, err := u.quests.ListPublic()
	if err != nil {
		return nil, custom.NewUnexpectedError("failed to list quests")
	}
	return ResponseQuests(list), nil
}

func (u *questUseCase) ListForUser(authenticated bool) ([]dto.QuestResponse, error) {
	if authenticated {
		list, err := u.quests.ListAll()
		if err != nil {
			return nil, custom.NewUnexpectedError("failed to list quests")
		}
		return ResponseQuests(list), nil
	}
	list, err := u.quests.ListPublic()
	if err != nil {
		return nil, custom.NewUnexpectedError("failed to list quests")
	}
	return ResponseQuests(list), nil
}

func (u *questUseCase) Create(userID string, in *dto.CreateQuestInput) error {
	if err := helper.ValidateDescription(in.Description); err != nil {
		return custom.NewBadRequestError("invalid description")
	}
	// if err := helper.ValidateImages(len(in.Images)); err != nil {
	// 	return custom.NewBadRequestError("invalid images")
	// }
	if _, err := u.quests.FindByID(in.QuestLevelID); err != nil {
		return custom.NewNotFoundError("quest level not found")
	}
	// imgJSON, _ := json.Marshal(in.Images)
	m := &model.Quest{
		UserID:       helper.ParseUUIDOrNil(userID),
		Title:        in.Title,
		Description:  in.Description,
		QuestLevelID: helper.ParseUUIDOrNil(in.QuestLevelID),
		// Images:       imgJSON,
		Privacy: in.Privacy,
		Status:  model.ItemStatusActive,
	}
	if _, err := u.quests.Create(m); err != nil {
		return custom.NewUnexpectedError("failed to create quest")
	}
	return nil
}

func (u *questUseCase) Update(userID string, id string, in *dto.UpdateQuestInput) error {
	m, err := u.quests.FindByID(id)
	if err != nil {
		return custom.NewNotFoundError("quest not found")
	}
	if m.UserID != helper.ParseUUIDOrNil(userID) {
		return custom.NewForbiddenError("forbidden")
	}
	if m.Status == model.ItemStatusArchived {
		return custom.NewForbiddenError("cannot modify archived")
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
	if in.QuestLevelID != nil {
		if _, err := u.questLevels.FindByID(*in.QuestLevelID); err != nil {
			return custom.NewNotFoundError("quest level not found")
		}
		m.QuestLevelID = helper.ParseUUIDOrNil(*in.QuestLevelID)
	}
	if in.Privacy != nil {
		m.Privacy = *in.Privacy
	}
	if _, err := u.quests.Update(m); err != nil {
		return custom.NewUnexpectedError("failed to update quest")
	}
	return nil
}

func (u *questUseCase) Delete(userID string, id string) error {
	m, err := u.quests.FindByID(id)
	if err != nil {
		return custom.NewNotFoundError("quest not found")
	}
	if m.UserID != helper.ParseUUIDOrNil(userID) {
		return custom.NewForbiddenError("forbidden")
	}
	return u.quests.Delete(id)
}
