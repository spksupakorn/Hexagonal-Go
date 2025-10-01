package helper

import (
	"crypto/rand"
	"path/filepath"

	"dungeons-dragon-service/internal/config"
	"dungeons-dragon-service/internal/domain/model"
	"dungeons-dragon-service/internal/http/custom"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

const (
	maxImages          = 10
	maxDescriptionSize = 5000
)

func ValidatePrivacy(p model.Privacy) error {
	if p != model.PrivacyPublic && p != model.PrivacyPrivate {
		return custom.NewBadRequestError("invalid privacy")
	}
	return nil
}

func ValidateImages(length int) error {
	if length > maxImages {
		return custom.NewBadRequestError(fmt.Sprintf("images must be at most %d", maxImages))
	}
	return nil
}

func ValidateDescription(desc string) error {
	if len([]rune(desc)) > maxDescriptionSize {
		return custom.NewBadRequestError(fmt.Sprintf("description must be at most %d characters", maxDescriptionSize))
	}
	return nil
}

func GenerateSalt(length int) (string, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(salt), nil
}

func HashPasswordArgon2(password, salt string) string {
	hash := argon2.IDKey([]byte(password), []byte(salt), 1, 64*1024, 4, 32)
	return base64.RawStdEncoding.EncodeToString(hash) + ":" + salt
}
func VerifyPasswordArgon2(password, hashed string) bool {
	parts := strings.Split(hashed, ":")
	if len(parts) != 2 {
		return false
	}
	hash := argon2.IDKey([]byte(password), []byte(parts[1]), 1, 64*1024, 4, 32)
	return base64.RawStdEncoding.EncodeToString(hash) == parts[0]
}

// change string to uuid or return nil uuid
func ParseUUIDOrNil(s string) uuid.UUID {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}
	return id
}

func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func DeleteFileIfExists(path string) error {
	if _, err := os.Stat(path); err == nil {
		return os.Remove(path)
	}
	return nil
}

func GetImageURL(path string) string {
	// Assuming the server serves images from the /images/ directory
	return config.GetConfigString("DOMAIN") + "/api/v1/pictures/" + filepath.Base(path)
}
