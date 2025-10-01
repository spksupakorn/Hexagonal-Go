package repositories

import (
	"dungeons-dragon-service/internal/domain/model"
	"dungeons-dragon-service/internal/domain/repository"

	"gorm.io/gorm"
)

type imageRepo struct{ db *gorm.DB }

func NewImageRepo(db *gorm.DB) repository.ImageRepository {
	return &imageRepo{db: db}
}

func (r *imageRepo) GetCharacterImageByID(characterID string) ([]model.CharacterImage, error) {
	var imgs []model.CharacterImage
	if err := r.db.Where("character_id = ?", characterID).Find(&imgs).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return imgs, nil
}

func (r *imageRepo) GetQuestImageByID(questID string) ([]model.QuestImage, error) {
	var imgs []model.QuestImage
	if err := r.db.Where("quest_id = ?", questID).Find(&imgs).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return imgs, nil
}

func (r *imageRepo) DeleteCharacterImageByID(characterID string) error {
	return r.db.Where("character_id = ?", characterID).Delete(&model.CharacterImage{}).Error
}

func (r *imageRepo) CreateCharacterImage(img *model.CharacterImage) (*model.CharacterImage, error) {
	if err := r.db.Create(img).Error; err != nil {
		return nil, err
	}
	return img, nil
}

func (r *imageRepo) DeleteQuestImageByID(questID string) error {
	return r.db.Where("quest_id = ?", questID).Delete(&model.QuestImage{}).Error
}

func (r *imageRepo) CreateQuestImage(img *model.QuestImage) (*model.QuestImage, error) {
	if err := r.db.Create(img).Error; err != nil {
		return nil, err
	}
	return img, nil
}
