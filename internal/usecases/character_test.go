package usecases

import (
	"dungeons-dragon-service/internal/domain/model"
	"dungeons-dragon-service/internal/dto"
	"dungeons-dragon-service/internal/helper"
	"errors"
	"mime/multipart"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockCharRepo struct {
	m map[string]*model.Character
}

func newMockCharRepo() *mockCharRepo {
	return &mockCharRepo{m: map[string]*model.Character{}}
}

func (m *mockCharRepo) Create(c *model.Character) (*model.Character, error) {
	m.m[c.ID.String()] = c
	return c, nil
}

func (m *mockCharRepo) Update(c *model.Character) (*model.Character, error) {
	if _, ok := m.m[c.ID.String()]; !ok {
		return nil, errors.New("not found")
	}
	m.m[c.ID.String()] = c
	return c, nil
}

func (m *mockCharRepo) Delete(id string) error {
	if _, ok := m.m[id]; !ok {
		return errors.New("not found")
	}
	delete(m.m, id)
	return nil
}

func (m *mockCharRepo) FindByID(id string) (*model.Character, error) {
	c, ok := m.m[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return c, nil
}

func (m *mockCharRepo) ListAll() ([]model.Character, error) {
	var chars []model.Character
	for _, c := range m.m {
		chars = append(chars, *c)
	}
	return chars, nil
}

func (m *mockCharRepo) ListPublic() ([]model.Character, error) {
	var chars []model.Character
	for _, c := range m.m {
		if c.Privacy == model.PrivacyPublic {
			chars = append(chars, *c)
		}
	}
	return chars, nil
}

func (m *mockCharRepo) ListByUser(userID string) ([]model.Character, error) {
	var chars []model.Character
	for _, c := range m.m {
		if c.UserID == helper.ParseUUIDOrNil(userID) {
			chars = append(chars, *c)
		}
	}
	return chars, nil
}

func (m *mockCharRepo) ArchiveByClassID(classID string) error {
	found := false
	for _, c := range m.m {
		if c.ClassID == helper.ParseUUIDOrNil(classID) {
			c.Status = model.ItemStatusArchived
			found = true
		}
	}
	if !found {
		return errors.New("not found")
	}
	return nil
}

func (m *mockCharRepo) ArchiveByRaceID(raceID string) error {
	found := false
	for _, c := range m.m {
		if c.RaceID == helper.ParseUUIDOrNil(raceID) {
			c.Status = model.ItemStatusArchived
			found = true
		}
	}
	if !found {
		return errors.New("not found")
	}
	return nil
}

type mockClassRepo struct {
	m map[string]*model.Class
}

func (m *mockClassRepo) Create(c *model.Class) (*model.Class, error) {
	m.m[c.ID.String()] = c
	return c, nil
}

func (m *mockClassRepo) Update(c *model.Class) (*model.Class, error) {
	if _, ok := m.m[c.ID.String()]; !ok {
		return nil, errors.New("not found")
	}
	m.m[c.ID.String()] = c
	return c, nil
}

func (m *mockClassRepo) Delete(id string) error {
	if _, ok := m.m[id]; !ok {
		return errors.New("not found")
	}
	delete(m.m, id)
	return nil
}

func (m *mockClassRepo) FindByID(id string) (*model.Class, error) {
	c, ok := m.m[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return c, nil
}

func (m *mockClassRepo) List() ([]model.Class, error) {
	var classes []model.Class
	for _, c := range m.m {
		classes = append(classes, *c)
	}
	return classes, nil
}

type mockRaceRepo struct {
	m map[string]*model.Race
}

func (m *mockRaceRepo) Create(r *model.Race) (*model.Race, error) {
	m.m[r.ID.String()] = r
	return r, nil
}

func (m *mockRaceRepo) Update(r *model.Race) (*model.Race, error) {
	if _, ok := m.m[r.ID.String()]; !ok {
		return nil, errors.New("not found")
	}
	m.m[r.ID.String()] = r
	return r, nil
}

func (m *mockRaceRepo) Delete(id string) error {
	if _, ok := m.m[id]; !ok {
		return errors.New("not found")
	}
	delete(m.m, id)
	return nil
}

func (m *mockRaceRepo) FindByID(id string) (*model.Race, error) {
	r, ok := m.m[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return r, nil
}

func (m *mockRaceRepo) List() ([]model.Race, error) {
	var races []model.Race
	for _, r := range m.m {
		races = append(races, *r)
	}
	return races, nil
}

func TestCharacterCreateValidation(t *testing.T) {
	charRepo := newMockCharRepo()
	classRepo := mockClassRepo{m: map[string]*model.Class{"f6d28968-b689-4c50-b4cc-03ab84b47039": {Name: "Warrior"}}}
	raceRepo := mockRaceRepo{m: map[string]*model.Race{"4fa768c3-79a2-4362-845b-5b869784d7c7": {Name: "Elf"}}}
	uc := NewCharacterUsecase(charRepo, &classRepo, &raceRepo)

	imageUc := NewImageUsecase(nil, charRepo, nil)
	//test image upload
	var img []*multipart.FileHeader
	//set image to 11
	for i := 0; i < 11; i++ {
		img = append(img, &multipart.FileHeader{Filename: "a.jpg"})
	}
	err := imageUc.UploadCharacterImage("1680b136-8862-4ea4-9d80-b2a6a7e71988", "72aa7e28-47d7-4625-bc2d-9790e78c9025", img)
	require.Error(t, err)

	// Too long description
	long := make([]rune, 5001)
	_, err = uc.Create("1680b136-8862-4ea4-9d80-b2a6a7e71988", &dto.CreateCharacterInput{
		Title:       "Hero",
		Description: string(long),
		ClassID:     "f6d28968-b689-4c50-b4cc-03ab84b47039",
		RaceID:      "4fa768c3-79a2-4362-845b-5b869784d7c7",
	})
	require.Error(t, err)

	// Missing class
	_, err = uc.Create("1680b136-8862-4ea4-9d80-b2a6a7e71988", &dto.CreateCharacterInput{
		Title:       "Hero",
		Description: "ok",
		RaceID:      "4fa768c3-79a2-4362-845b-5b869784d7c7",
	})
	require.Error(t, err)

	// Success
	char, err := uc.Create("1680b136-8862-4ea4-9d80-b2a6a7e71988", &dto.CreateCharacterInput{
		Title:       "Hero",
		Description: "ok",
		ClassID:     "f6d28968-b689-4c50-b4cc-03ab84b47039",
		RaceID:      "4fa768c3-79a2-4362-845b-5b869784d7c7",
	})
	require.NoError(t, err)
	require.Equal(t, "1680b136-8862-4ea4-9d80-b2a6a7e71988", char.UserID)
}

func TestCharacterUpdateOwnershipAndArchive(t *testing.T) {
	charRepo := newMockCharRepo()
	classRepo := mockClassRepo{m: map[string]*model.Class{"f6d28968-b689-4c50-b4cc-03ab84b47039": {Name: "Warrior"}}}
	raceRepo := mockRaceRepo{m: map[string]*model.Race{"4fa768c3-79a2-4362-845b-5b869784d7c7": {Name: "Elf"}}}
	uc := NewCharacterUsecase(charRepo, &classRepo, &raceRepo)

	// Create a character
	char, _ := uc.Create("00ec53c1-276b-4d9f-944c-637e75475650", &dto.CreateCharacterInput{
		Title:       "Hero",
		Description: "ok",
		ClassID:     "f6d28968-b689-4c50-b4cc-03ab84b47039",
		RaceID:      "4fa768c3-79a2-4362-845b-5b869784d7c7",
	})

	// Forbidden update
	err := uc.Update("1680b136-8862-4ea4-9d80-b2a6a7e71988", char.ID, &dto.UpdateCharacterInput{
		Title: strPtr("X"),
	})
	require.Error(t, err)

	// Archive via option delete then attempt update
	_ = charRepo.ArchiveByClassID("f6d28968-b689-4c50-b4cc-03ab84b47039")
	err = uc.Update("00ec53c1-276b-4d9f-944c-637e75475650", char.ID, &dto.UpdateCharacterInput{
		Title: strPtr("X"),
	})
	require.Error(t, err)

}

func strPtr(s string) *string { return &s }
