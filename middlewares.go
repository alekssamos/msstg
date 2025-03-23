package main

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"gorm.io/gorm"
)

func dbUser(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		var user User
		var vctx context.Context
		var ChatID int64
		db := DB(ctx)
		if update.CallbackQuery != nil {
			ChatID = update.CallbackQuery.Message.Message.Chat.ID
		}
		if update.Message != nil {
			ChatID = update.Message.Chat.ID
		}
		result := db.First(&user, ChatID)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				fmt.Println("Пользователь в базе не найден, сейчас создадим")
				if update.Message != nil {
					firstHandler(ctx, b, update)
				}
				user = User{ChatID: ChatID}
				result = db.Create(&user)
				if result.Error != nil {
					fmt.Println("Ошибка создания записи:", result.Error)
					return
				}
			} else {
				fmt.Println("Ошибка поиска записи:", result.Error)
			}
			return
		}
		// fmt.Printf("Найден пользователь: %+v\n", user)
		vctx = context.WithValue(ctx, ctxUserKey{}, &user)
		next(vctx, b, update)
	}
}
