package repositories

import (
	"dungeons-dragon-service/internal/domain/model"
	"dungeons-dragon-service/internal/domain/repository"

	"gorm.io/gorm"
)

type questRepo struct{ db *gorm.DB }

func NewQuestRepo(db *gorm.DB) repository.QuestRepository {
	return &questRepo{db: db}
}

func (r *questRepo) Create(m *model.Quest) (*model.Quest, error) {
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return m, nil
}
func (r *questRepo) Update(m *model.Quest) (*model.Quest, error) {
	if err := r.db.Save(m).Error; err != nil {
		return nil, err
	}
	return m, nil
}
func (r *questRepo) Delete(id string) error { return r.db.Delete(&model.Quest{}, id).Error }
func (r *questRepo) FindByID(id string) (*model.Quest, error) {
	var m model.Quest
	if err := r.db.Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}
func (r *questRepo) ListAll() ([]model.Quest, error) {
	var list []model.Quest
	err := r.db.Where("status = ?", model.ItemStatusActive).Order("created_at desc").Find(&list).Error
	return list, err
}
func (r *questRepo) ListPublic() ([]model.Quest, error) {
	var list []model.Quest
	err := r.db.Where("privacy = ? AND status = ?", model.PrivacyPublic, model.ItemStatusActive).Order("created_at desc").Find(&list).Error
	return list, err
}
func (r *questRepo) ListByUser(userID string) ([]model.Quest, error) {
	var list []model.Quest
	err := r.db.Where("user_id = ? AND status = ?", userID, model.ItemStatusActive).Order("created_at desc").Find(&list).Error
	return list, err
}
func (r *questRepo) ArchiveByQuestLevelID(difficultyID string) error {
	return r.db.Model(&model.Quest{}).
		Where("difficulty_id = ? AND status = ?", difficultyID, model.ItemStatusActive).
		Update("status", model.ItemStatusArchived).Error
}
