package repositories

import (
	"dungeons-dragon-service/internal/domain/model"
	"dungeons-dragon-service/internal/domain/repository"

	"gorm.io/gorm"
)

type classRepo struct{ db *gorm.DB }
type raceRepo struct{ db *gorm.DB }
type questLevelRepo struct{ db *gorm.DB }

func NewClassRepo(db *gorm.DB) repository.ClassRepository { return &classRepo{db} }
func NewRaceRepo(db *gorm.DB) repository.RaceRepository   { return &raceRepo{db} }
func NewQuestLevelRepo(db *gorm.DB) repository.QuestLevelRepository {
	return &questLevelRepo{db}
}

func (r *classRepo) Create(m *model.Class) (*model.Class, error) {
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return m, nil
}
func (r *classRepo) Update(m *model.Class) (*model.Class, error) {
	if err := r.db.Save(m).Error; err != nil {
		return nil, err
	}
	return m, nil
}
func (r *classRepo) Delete(id string) error { return r.db.Delete(&model.Class{}, id).Error }
func (r *classRepo) FindByID(id string) (*model.Class, error) {
	var m model.Class
	if err := r.db.Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}
func (r *classRepo) List() ([]model.Class, error) {
	var list []model.Class
	return list, r.db.Order("name asc").Find(&list).Error
}

func (r *raceRepo) Create(m *model.Race) (*model.Race, error) {
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return m, nil
}
func (r *raceRepo) Update(m *model.Race) (*model.Race, error) {
	if err := r.db.Save(m).Error; err != nil {
		return nil, err
	}
	return m, nil
}
func (r *raceRepo) Delete(id string) error { return r.db.Delete(&model.Race{}, id).Error }
func (r *raceRepo) FindByID(id string) (*model.Race, error) {
	var m model.Race
	if err := r.db.Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}
func (r *raceRepo) List() ([]model.Race, error) {
	var list []model.Race
	return list, r.db.Order("name asc").Find(&list).Error
}

func (r *questLevelRepo) Create(m *model.QuestLevel) (*model.QuestLevel, error) {
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return m, nil
}
func (r *questLevelRepo) Update(m *model.QuestLevel) (*model.QuestLevel, error) {
	if err := r.db.Save(m).Error; err != nil {
		return nil, err
	}
	return m, nil
}
func (r *questLevelRepo) Delete(id string) error {
	return r.db.Delete(&model.QuestLevel{}, id).Error
}
func (r *questLevelRepo) FindByID(id string) (*model.QuestLevel, error) {
	var m model.QuestLevel
	if err := r.db.Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}
func (r *questLevelRepo) List() ([]model.QuestLevel, error) {
	var list []model.QuestLevel
	return list, r.db.Order("name asc").Find(&list).Error
}
