package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ChatID     int64     `gorm:"primaryKey;not null;unique"`
	CreatedAt  time.Time // Время создания записи (заполняется автоматически)
	UpdatedAt  time.Time // Время последнего обновления записи (заполняется автоматически)
	Lang       string    `gorm:"size:3;default:ru"`
	VoiceName  string    `gorm:"size:64;default:ru-RU-SvetlanaNeural"`
	VoiceRate  int       `gorm:"default:0"`
	VoicePitch int       `gorm:"default:0"`
	Books      []Book    `gorm:"foreignKey:UserID"`
}

type Book struct {
	ID                uint       `gorm:"primaryKey"`
	UserID            int64      // Внешний ключ
	CreatedAt         time.Time  // Время создания записи (заполняется автоматически)
	UpdatedAt         time.Time  // Время последнего обновления записи (заполняется автоматически)
	Filename          string     `gorm:"size:256;not null;unique"`
	ConvertedFilename string     `gorm:"size:256;unique"`
	BookParts         []BookPart `gorm:"foreignKey:BookID"`
}

type BookPart struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time // Время создания записи (заполняется автоматически)
	UpdatedAt time.Time // Время последнего обновления записи (заполняется автоматически)
	Chapter   uint
	Text      string `gorm:"size:60000"`
	BookID    uint   // Внешний ключ
}

// helpers

func checkFile(filename string) error {
	fi, er := os.Stat(filename)
	if os.IsNotExist(er) {
		return fmt.Errorf("the file %q dus not exists", filename)
	}
	if er != nil {
		return er
	}
	if fi.IsDir() {
		return fmt.Errorf("%q - is a directory", filename)
	}
	return nil
}

// hooks:

func (book *Book) BeforeCreate(tx *gorm.DB) (err error) {
	var result error
	result = checkFile(book.Filename)
	if result != nil {
		return result
	}
	if book.ConvertedFilename != "" {
		result = checkFile(book.ConvertedFilename)
		if result != nil {
			return result
		}
	}
	return result
}

func (book *Book) BeforeDelete(tx *gorm.DB) (err error) {
	os.Remove(book.Filename)
	if book.ConvertedFilename != "" {
		os.Remove(book.ConvertedFilename)
	}
	return nil
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	if user.ChatID == 0 {
		return fmt.Errorf("ChatID is required")
	}
	if user.VoiceRate < -100 || user.VoiceRate > 100 || user.VoicePitch < -100 || user.VoicePitch > 100 {
		return fmt.Errorf("voice rate and voice pitch must be value between from -100 to +100")
	}
	return nil
}

// ctx getters

func DB(ctx context.Context) *gorm.DB {
	v := ctx.Value(ctxDbKey{})
	switch x := v.(type) {
	case *gorm.DB:
		return x
	default:
		panic("unknown var db type in ctx")
	}
}

func USER(ctx context.Context) *User {
	v := ctx.Value(ctxUserKey{})
	switch x := v.(type) {
	case *User:
		return x
	default:
		panic("unknown var user type in ctx")
	}
}
