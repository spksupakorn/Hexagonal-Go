package usecases

import (
	"dungeons-dragon-service/internal/domain/model"
	"dungeons-dragon-service/internal/domain/repository"
	"dungeons-dragon-service/internal/dto"
	"dungeons-dragon-service/internal/http/custom"
)

type OptionUseCase interface {
	// Classes
	CreateClass(name string) error
	UpdateClass(id string, name string) error
	DeleteClass(id string) error
	ListClasses() ([]dto.ClassResponse, error)

	// Races
	CreateRace(name string) error
	UpdateRace(id string, name string) error
	DeleteRace(id string) error
	ListRaces() ([]dto.RaceResponse, error)

	// Quest Levels
	CreateQuestLevel(name string) error
	UpdateQuestLevel(id string, name string) error
	DeleteQuestLevel(id string) error
	ListQuestLevels() ([]dto.QuestLevelResponse, error)
}

type optionUseCase struct {
	classes     repository.ClassRepository
	races       repository.RaceRepository
	questLevels repository.QuestLevelRepository

	chars  repository.CharacterRepository
	quests repository.QuestRepository
}

func NewOptionUseCase(c repository.ClassRepository, r repository.RaceRepository, d repository.QuestLevelRepository,
	char repository.CharacterRepository, q repository.QuestRepository) OptionUseCase {
	return &optionUseCase{classes: c, races: r, questLevels: d, chars: char, quests: q}
}

func ResponseClasses(c []model.Class) []dto.ClassResponse {
	res := make([]dto.ClassResponse, len(c))
	for i, class := range c {
		res[i] = dto.ClassResponse{
			ID:   class.ID.String(),
			Name: class.Name,
		}
	}
	return res
}

func ResponseRaces(r []model.Race) []dto.RaceResponse {
	res := make([]dto.RaceResponse, len(r))
	for i, race := range r {
		res[i] = dto.RaceResponse{
			ID:   race.ID.String(),
			Name: race.Name,
		}
	}
	return res
}

func ResponseQuestLevels(d []model.QuestLevel) []dto.QuestLevelResponse {
	res := make([]dto.QuestLevelResponse, len(d))
	for i, questLevel := range d {
		res[i] = dto.QuestLevelResponse{
			ID:   questLevel.ID.String(),
			Name: questLevel.Name,
		}
	}
	return res
}

// Classes
func (u *optionUseCase) CreateClass(name string) error {
	if name == "" {
		return custom.NewBadRequestError("name required")
	}
	m := &model.Class{Name: name}
	_, err := u.classes.Create(m)
	if err != nil {
		return custom.NewUnexpectedError("failed to create class")
	}
	return nil
}
func (u *optionUseCase) UpdateClass(id string, name string) error {
	m, err := u.classes.FindByID(id)
	if err != nil {
		return custom.NewNotFoundError("class not found")
	}
	m.Name = name
	_, err = u.classes.Update(m)
	if err != nil {
		return custom.NewUnexpectedError("failed to update class")
	}
	return nil
}
func (u *optionUseCase) DeleteClass(id string) error {
	// Archive related characters, then delete class
	if err := u.chars.ArchiveByClassID(id); err != nil {
		return err
	}
	return u.classes.Delete(id)
}
func (u *optionUseCase) ListClasses() ([]dto.ClassResponse, error) {
	list, err := u.classes.List()
	if err != nil {
		return nil, custom.NewUnexpectedError("failed to list classes")
	}
	return ResponseClasses(list), nil
}

// Races
func (u *optionUseCase) CreateRace(name string) error {
	if name == "" {
		return custom.NewBadRequestError("name required")
	}
	m := &model.Race{Name: name}
	_, err := u.races.Create(m)
	return err
}
func (u *optionUseCase) UpdateRace(id string, name string) error {
	m, err := u.races.FindByID(id)
	if err != nil {
		return custom.NewNotFoundError("race not found")
	}
	m.Name = name
	_, err = u.races.Update(m)
	if err != nil {
		return custom.NewUnexpectedError("failed to update race")
	}
	return nil
}
func (u *optionUseCase) DeleteRace(id string) error {
	if err := u.chars.ArchiveByRaceID(id); err != nil {
		return err
	}
	return u.races.Delete(id)
}
func (u *optionUseCase) ListRaces() ([]dto.RaceResponse, error) {
	list, err := u.races.List()
	if err != nil {
		return nil, custom.NewUnexpectedError("failed to list races")
	}
	return ResponseRaces(list), nil
}

// Difficulties
func (u *optionUseCase) CreateQuestLevel(name string) error {
	if name == "" {
		return custom.NewBadRequestError("name required")
	}
	m := &model.QuestLevel{Name: name}
	_, err := u.questLevels.Create(m)
	if err != nil {
		return custom.NewUnexpectedError("failed to create quest level")
	}
	return nil
}
func (u *optionUseCase) UpdateQuestLevel(id string, name string) error {
	m, err := u.questLevels.FindByID(id)
	if err != nil {
		return custom.NewNotFoundError("quest level not found")
	}
	m.Name = name
	_, err = u.questLevels.Update(m)
	if err != nil {
		return custom.NewUnexpectedError("failed to update quest level")
	}
	return nil
}
func (u *optionUseCase) DeleteQuestLevel(id string) error {
	if err := u.quests.ArchiveByQuestLevelID(id); err != nil {
		return err
	}
	return u.questLevels.Delete(id)
}
func (u *optionUseCase) ListQuestLevels() ([]dto.QuestLevelResponse, error) {
	list, err := u.questLevels.List()
	if err != nil {
		return nil, custom.NewUnexpectedError("failed to list quest levels")
	}
	return ResponseQuestLevels(list), nil
}
