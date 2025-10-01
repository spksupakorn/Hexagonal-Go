package usecases

import (
	"dungeons-dragon-service/internal/domain/model"
	"dungeons-dragon-service/internal/dto"
	"dungeons-dragon-service/internal/helper"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

// Mock implementation for QuestLevelRepository
type mockQuestLevelRepo struct {
	levels map[string]*model.QuestLevel
}

func (m *mockQuestLevelRepo) Create(q *model.QuestLevel) (*model.QuestLevel, error) {
	m.levels[q.ID.String()] = q
	return q, nil
}

func (m *mockQuestLevelRepo) Update(q *model.QuestLevel) (*model.QuestLevel, error) {
	if _, exists := m.levels[q.ID.String()]; !exists {
		return nil, errors.New("not found")
	}
	m.levels[q.ID.String()] = q
	return q, nil
}

func (m *mockQuestLevelRepo) Delete(id string) error {
	if _, exists := m.levels[id]; !exists {
		return errors.New("not found")
	}
	delete(m.levels, id)
	return nil
}

func (m *mockQuestLevelRepo) FindByID(id string) (*model.QuestLevel, error) {
	if q, exists := m.levels[id]; exists {
		return q, nil
	}
	return nil, errors.New("not found")
}

func (m *mockQuestLevelRepo) List() ([]model.QuestLevel, error) {
	var res []model.QuestLevel
	for _, v := range m.levels {
		res = append(res, *v)
	}
	return res, nil
}

// Mock implementation for QuestRepository
type mockQuestRepo struct {
	quests   map[string]*model.Quest
	archived []string
}

func (m *mockQuestRepo) Create(q *model.Quest) (*model.Quest, error) {
	m.quests[q.ID.String()] = q
	return q, nil
}

func (m *mockQuestRepo) Update(q *model.Quest) (*model.Quest, error) {
	if _, exists := m.quests[q.ID.String()]; !exists {
		return nil, errors.New("not found")
	}
	m.quests[q.ID.String()] = q
	return q, nil
}

func (m *mockQuestRepo) Delete(id string) error {
	if _, exists := m.quests[id]; !exists {
		return errors.New("not found")
	}
	delete(m.quests, id)
	return nil
}

func (m *mockQuestRepo) FindByID(id string) (*model.Quest, error) {
	if q, exists := m.quests[id]; exists {
		return q, nil
	}
	return nil, errors.New("not found")
}

func (m *mockQuestRepo) ListAll() ([]model.Quest, error) {
	var res []model.Quest
	for _, v := range m.quests {
		res = append(res, *v)
	}
	return res, nil
}

func (m *mockQuestRepo) ListPublic() ([]model.Quest, error) {
	var res []model.Quest
	for _, v := range m.quests {
		if v.Privacy == model.PrivacyPublic {
			res = append(res, *v)
		}
	}
	return res, nil
}

func (m *mockQuestRepo) ListByUser(userID string) ([]model.Quest, error) {
	var res []model.Quest
	for _, v := range m.quests {
		if v.UserID == helper.ParseUUIDOrNil(userID) {
			res = append(res, *v)
		}
	}
	return res, nil
}

func (m *mockQuestRepo) ArchiveByQuestLevelID(questsLevelID string) error {
	m.archived = append(m.archived, questsLevelID)
	return nil
}

func TestOptionDeleteArchives(t *testing.T) {
	charRepo := newMockCharRepo()
	classRepo := mockClassRepo{
		m: make(map[string]*model.Class),
	}
	classRepo.m["3c75ef02-b390-423b-86fc-99c590921f29"] = &model.Class{Name: "Warrior"}
	classRepo.m["c3e2b2e0-53f1-44b3-8f2e-8eec1d5a6f72"] = &model.Class{Name: "Mage"}
	raceRepo := mockRaceRepo{
		m: make(map[string]*model.Race),
	}
	raceRepo.m["66e9e0cf-8b74-4e73-8c90-7a2d4351f2e6"] = &model.Race{Name: "Human"}
	raceRepo.m["a96e8f7c-4a47-4a5c-9a6c-4b2a24e7a3c9"] = &model.Race{Name: "Elf"}
	questLevelRepo := mockQuestLevelRepo{
		levels: make(map[string]*model.QuestLevel),
	}
	questLevelRepo.levels["d1e8e1de-5f6a-4ff0-84d6-d13d7b9c0e9b"] = &model.QuestLevel{Name: "Easy"}
	questLevelRepo.levels["b6e3f5d4-3b8f-4eaf-bd77-cb4a2f11e5c1"] = &model.QuestLevel{Name: "Hard"}
	questRepo := mockQuestRepo{}

	uc := NewOptionUseCase(&classRepo, &raceRepo, &questLevelRepo, charRepo, &questRepo)

	// Create a character using class and race
	_, _ = NewCharacterUsecase(charRepo, &classRepo, &raceRepo).Create("f6d28968-b689-4c50-b4cc-03ab84b47039", &dto.CreateCharacterInput{
		Title:       "Hero",
		Description: "ok",
		ClassID:     "3c75ef02-b390-423b-86fc-99c590921f29",
		RaceID:      "66e9e0cf-8b74-4e73-8c90-7a2d4351f2e6",
		Privacy:     "public",
	})

	// Delete class
	err := uc.DeleteClass("3c75ef02-b390-423b-86fc-99c590921f29")
	require.NoError(t, err)

	//check character is archived
	all, _ := charRepo.ListAll()
	require.Equal(t, model.ItemStatusArchived, all[0].Status)
}
