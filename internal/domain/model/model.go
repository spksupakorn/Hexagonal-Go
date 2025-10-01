package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Privacy string
type ItemStatus string
type Role string

const (
	PrivacyPublic  Privacy = "public"
	PrivacyPrivate Privacy = "private"

	ItemStatusActive   ItemStatus = "active"
	ItemStatusArchived ItemStatus = "archived"

	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type Base struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt time.Time      `gorm:"type:timestamptz;not null;autoCreateTime;"`
	UpdatedAt time.Time      `gorm:"type:timestamptz;not null;autoUpdateTime;"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamptz;->:false;column:deleted_at" json:"-"`
}

// Users table
type User struct {
	Base
	Username     string `gorm:"type:varchar(64);unique;not null"`
	Email        string `gorm:"type:varchar(128);unique;not null"`
	PasswordHash string `gorm:"type:varchar(255);not null"`
	Role         Role   `gorm:"type:user_role;default:'user';not null"`
}

// Characters table
type Character struct {
	Base
	UserID      uuid.UUID        `gorm:"type:uuid;not null"`
	User        *User            `gorm:"foreignKey:UserID"`
	Title       string           `gorm:"type:varchar(128);not null"`
	Description string           `gorm:"type:text;not null;size:5000"`
	ClassID     uuid.UUID        `gorm:"type:uuid;not null"`
	Class       *Class           `gorm:"foreignKey:ClassID"`
	RaceID      uuid.UUID        `gorm:"type:uuid;not null"`
	ImagePath   datatypes.JSON   `gorm:"type:jsonb;default:'[]'::jsonb"`
	Race        *Race            `gorm:"foreignKey:RaceID"`
	Privacy     Privacy          `gorm:"type:privacy;default:'public';not null"`
	Status      ItemStatus       `gorm:"type:item_status;default:'active';not null"`
	Images      []CharacterImage `gorm:"foreignKey:CharacterID"`
}

// CharacterClasses table
type Class struct {
	Base
	Name      string `gorm:"type:varchar(128);unique;not null"`
	IsDeleted bool   `gorm:"not null;default:false"`
}

// CharacterRaces table
type Race struct {
	Base
	Name      string `gorm:"type:varchar(128);unique;not null"`
	IsDeleted bool   `gorm:"not null;default:false"`
}

// Quests table
type Quest struct {
	Base
	UserID       uuid.UUID      `gorm:"type:uuid;not null"`
	User         *User          `gorm:"foreignKey:UserID"`
	Title        string         `gorm:"type:varchar(128);not null"`
	Description  string         `gorm:"type:text;not null;size:5000"`
	QuestLevelID uuid.UUID      `gorm:"type:uuid;not null"`
	QuestLevel   *QuestLevel    `gorm:"foreignKey:QuestLevelID"`
	ImagePath    datatypes.JSON `gorm:"type:jsonb;default:'[]'::jsonb"`
	Privacy      Privacy        `gorm:"type:privacy;default:'public';not null"`
	Status       ItemStatus     `gorm:"type:item_status;default:'active';not null"`
	Images       []QuestImage   `gorm:"foreignKey:QuestID"`
}

// QuestDifficulties table
type QuestLevel struct {
	Base
	Name      string `gorm:"type:varchar(128);unique;not null"`
	IsDeleted bool   `gorm:"not null;default:false"`
}

type CharacterImage struct {
	Base
	CharacterID uuid.UUID `gorm:"type:uuid;not null"`
	Path        string    `gorm:"type:text;not null"`
}

type QuestImage struct {
	Base
	QuestID uuid.UUID `gorm:"type:uuid;not null"`
	Path    string    `gorm:"type:text;not null"`
}
