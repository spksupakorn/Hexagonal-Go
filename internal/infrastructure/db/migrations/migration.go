package main

import (
	"dungeons-dragon-service/internal/config"
	"dungeons-dragon-service/internal/domain/model"
	"dungeons-dragon-service/internal/helper"
	database "dungeons-dragon-service/internal/infrastructure/db"

	"github.com/labstack/gommon/log"

	"gorm.io/gorm"
)

func main() {
	config.LoadConfig()
	db := database.NewPostgresDatabase()

	tx := db.ConnectDB().Begin()

	TaskMigration(tx)

	if err := tx.Commit().Error; err != nil {
		log.Fatalf("Migration failed: %v", err)
	} else {
		log.Info("Migration completed successfully")
	}
}

func TaskMigration(tx *gorm.DB) {
	if err := tx.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		log.Errorf("Error enabling uuid-ossp extension: %v", err)
	}

	enumStatements := []struct {
		name string
		sql  string
	}{
		{
			"user_role",
			`
			DO $$
			BEGIN
				IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN
					CREATE TYPE user_role AS ENUM ('user', 'admin');
				END IF;
			END$$;
			`,
		},
		{
			"privacy",
			`
			DO $$
			BEGIN
				IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'privacy') THEN
					CREATE TYPE privacy AS ENUM ('public', 'private');
				END IF;
			END$$;
			`,
		},
		{
			"item_status",
			`
			DO $$
			BEGIN
				IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'item_status') THEN
					CREATE TYPE item_status AS ENUM ('active', 'archived');
				END IF;
			END$$;
			`,
		},
	}

	for _, enum := range enumStatements {
		if err := tx.Exec(enum.sql).Error; err != nil {
			log.Fatalf("Error creating %s ENUM: %v", enum.name, err)
		}
	}

	log.Info("ENUM types created (if not exist)!")

	tx.Debug().AutoMigrate(
		&model.User{},
		&model.Class{},
		&model.Race{},
		&model.Character{},
		&model.QuestLevel{},
		&model.Quest{},
		&model.CharacterImage{},
		&model.QuestImage{},
	)

	// Insert pre data for Class
	classes := []model.Class{
		{Name: "Warrior"},
		{Name: "Mage"},
		{Name: "Archer"},
	}
	for _, c := range classes {
		if err := tx.FirstOrCreate(&c, model.Class{Name: c.Name}).Error; err != nil {
			log.Errorf("Error inserting class %s: %v", c.Name, err)
		}
	}

	// Insert pre data for Race
	races := []model.Race{
		{Name: "Human"},
		{Name: "Elf"},
		{Name: "Orc"},
	}
	for _, r := range races {
		if err := tx.FirstOrCreate(&r, model.Race{Name: r.Name}).Error; err != nil {
			log.Errorf("Error inserting race %s: %v", r.Name, err)
		}
	}

	// Insert pre data for QuestLevel
	questLevels := []model.QuestLevel{
		{Name: "Easy"},
		{Name: "Medium"},
		{Name: "Hard"},
	}
	for _, q := range questLevels {
		if err := tx.FirstOrCreate(&q, model.QuestLevel{Name: q.Name}).Error; err != nil {
			log.Errorf("Error inserting quest level %s: %v", q.Name, err)
		}
	}

	// Insert pre data for User
	salt, err := helper.GenerateSalt(16)
	if err != nil {
		log.Fatalf("Error generating salt: %v", err)
	}
	hashedPassword := helper.HashPasswordArgon2("password", salt)
	users := []model.User{
		{Username: "admin", Email: "admin@example.com", PasswordHash: hashedPassword, Role: model.RoleAdmin},
		{Username: "user", Email: "user@example.com", PasswordHash: hashedPassword, Role: model.RoleUser},
	}
	for _, u := range users {
		if err := tx.FirstOrCreate(&u, model.User{Username: u.Username}).Error; err != nil {
			log.Errorf("Error inserting user %s: %v", u.Username, err)
		}
	}

	// Insert pre data for Character
	characters := []model.Character{
		{Title: "Arthas", Privacy: "public"},
		{Title: "Sylvanas", Privacy: "private"},
		{Title: "Jaina", Privacy: "public"},
	}

	//get user for set user id
	var user model.User
	if err := tx.First(&user, "username = ?", "user").Error; err != nil {
		log.Fatalf("Error fetching user: %v", err)
	}
	characters[0].UserID = user.ID
	characters[1].UserID = user.ID
	characters[2].UserID = user.ID

	var class []model.Class
	if err := tx.Find(&class).Error; err != nil {
		log.Fatalf("Error fetching class: %v", err)
	}
	for i := range class {
		switch class[i].Name {
		case "Warrior":
			characters[0].ClassID = class[i].ID
		case "Archer":
			characters[1].ClassID = class[i].ID
		case "Mage":
			characters[2].ClassID = class[i].ID
		}
	}

	var race []model.Race
	if err := tx.Find(&race).Error; err != nil {
		log.Fatalf("Error fetching race: %v", err)
	}
	for i := range race {
		switch race[i].Name {
		case "Human":
			characters[0].RaceID = race[i].ID
		case "Elf":
			characters[1].RaceID = race[i].ID
		case "Orc":
			characters[2].RaceID = race[i].ID
		}
	}

	for _, ch := range characters {
		if err := tx.FirstOrCreate(&ch, model.Character{Title: ch.Title, ClassID: ch.ClassID, RaceID: ch.RaceID, UserID: ch.UserID}).Error; err != nil {
			log.Errorf("Error inserting character %s: %v", ch.Title, err)
			return // Stop further execution if error occurs
		}
	}

	// Insert pre data for Quest
	quests := []model.Quest{
		{Title: "Defeat the Dragon", Privacy: "public"},
		{Title: "Rescue the Princess", Privacy: "private"},
		{Title: "Find the Lost Sword", Privacy: "public"},
	}

	//get quest level for set quest level id
	var questLevel []model.QuestLevel
	if err := tx.Find(&questLevel).Error; err != nil {
		log.Fatalf("Error fetching quest level: %v", err)
	}
	for i := range questLevel {
		switch questLevel[i].Name {
		case "Hard":
			quests[0].QuestLevelID = questLevel[i].ID
		case "Medium":
			quests[1].QuestLevelID = questLevel[i].ID
		case "Easy":
			quests[2].QuestLevelID = questLevel[i].ID
		}
	}

	//set user id
	quests[0].UserID = user.ID
	quests[1].UserID = user.ID
	quests[2].UserID = user.ID
	for _, q := range quests {
		if err := tx.FirstOrCreate(&q, model.Quest{Title: q.Title, UserID: q.UserID, QuestLevelID: q.QuestLevelID, Privacy: q.Privacy}).Error; err != nil {
			log.Errorf("Error inserting quest %s: %v", q.Title, err)
			return // Stop further execution if error occurs
		}
	}
}
