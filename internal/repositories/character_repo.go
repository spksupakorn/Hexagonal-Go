package repositories

import (
	"dungeons-dragon-service/internal/domain/model"
	"dungeons-dragon-service/internal/domain/repository"

	"gorm.io/gorm"
)

type characterRepo struct{ db *gorm.DB }

func NewCharacterRepo(db *gorm.DB) repository.CharacterRepository {
	return &characterRepo{db: db}
}

func (r *characterRepo) Create(m *model.Character) (*model.Character, error) {
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return m, nil
}
func (r *characterRepo) Update(m *model.Character) (*model.Character, error) {
	if err := r.db.Save(m).Error; err != nil {
		return nil, err
	}
	return m, nil
}
func (r *characterRepo) Delete(id string) error { return r.db.Delete(&model.Character{}, id).Error }
func (r *characterRepo) FindByID(id string) (*model.Character, error) {
	var m model.Character
	if err := r.db.Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}
func (r *characterRepo) ListAll() ([]model.Character, error) {
	var list []model.Character
	err := r.db.Where("status = ?", model.ItemStatusActive).Order("created_at desc").Find(&list).Error
	return list, err
}
func (r *characterRepo) ListPublic() ([]model.Character, error) {
	var list []model.Character
	err := r.db.Where("privacy = ? AND status = ?", model.PrivacyPublic, model.ItemStatusActive).Order("created_at desc").Find(&list).Error
	return list, err
}
func (r *characterRepo) ListByUser(userID string) ([]model.Character, error) {
	var list []model.Character
	err := r.db.Where("user_id = ? AND status = ?", userID, model.ItemStatusActive).Order("created_at desc").Find(&list).Error
	return list, err
}
func (r *characterRepo) ArchiveByClassID(classID string) error {
	return r.db.Model(&model.Character{}).
		Where("class_id = ? AND status = ?", classID, model.ItemStatusActive).
		Update("status", model.ItemStatusArchived).Error
}
func (r *characterRepo) ArchiveByRaceID(raceID string) error {
	return r.db.Model(&model.Character{}).
		Where("race_id = ? AND status = ?", raceID, model.ItemStatusActive).
		Update("status", model.ItemStatusArchived).Error
}
