package repository

import "dungeons-dragon-service/internal/domain/model"

type UserRepository interface {
	Create(*model.User) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	FindByID(id string) (*model.User, error)
}

type ClassRepository interface {
	Create(*model.Class) (*model.Class, error)
	Update(*model.Class) (*model.Class, error)
	Delete(id string) error
	FindByID(id string) (*model.Class, error)
	List() ([]model.Class, error)
}

type RaceRepository interface {
	Create(*model.Race) (*model.Race, error)
	Update(*model.Race) (*model.Race, error)
	Delete(id string) error
	FindByID(id string) (*model.Race, error)
	List() ([]model.Race, error)
}

type QuestLevelRepository interface {
	Create(*model.QuestLevel) (*model.QuestLevel, error)
	Update(*model.QuestLevel) (*model.QuestLevel, error)
	Delete(id string) error
	FindByID(id string) (*model.QuestLevel, error)
	List() ([]model.QuestLevel, error)
}

type CharacterRepository interface {
	Create(*model.Character) (*model.Character, error)
	Update(*model.Character) (*model.Character, error)
	Delete(id string) error
	FindByID(id string) (*model.Character, error)
	ListAll() ([]model.Character, error)
	ListPublic() ([]model.Character, error)
	ListByUser(userID string) ([]model.Character, error)
	ArchiveByClassID(classID string) error
	ArchiveByRaceID(raceID string) error
}

type QuestRepository interface {
	Create(*model.Quest) (*model.Quest, error)
	Update(*model.Quest) (*model.Quest, error)
	Delete(id string) error
	FindByID(id string) (*model.Quest, error)
	ListAll() ([]model.Quest, error)
	ListPublic() ([]model.Quest, error)
	ListByUser(userID string) ([]model.Quest, error)
	ArchiveByQuestLevelID(questsLevelID string) error
}

type ImageRepository interface {
	GetCharacterImageByID(characterID string) ([]model.CharacterImage, error)
	GetQuestImageByID(questID string) ([]model.QuestImage, error)
	DeleteCharacterImageByID(characterID string) error
	CreateCharacterImage(img *model.CharacterImage) (*model.CharacterImage, error)
	DeleteQuestImageByID(questID string) error
	CreateQuestImage(img *model.QuestImage) (*model.QuestImage, error)
}
