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
		var call_first_handler bool
		result := db.First(&user, ChatID)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				fmt.Println("Пользователь в базе не найден, сейчас создадим")
				user = User{ChatID: ChatID}
				result = db.Create(&user)
				if result.Error != nil {
					fmt.Println("Ошибка создания записи:", result.Error)
					return
				}
				if update.Message != nil {
					call_first_handler = true
				}
			} else {
				fmt.Println("Ошибка поиска записи:", result.Error)
			}
		}
		// fmt.Printf("Найден пользователь: %+v\n", user)
		vctx = context.WithValue(ctx, ctxUserKey{}, &user)
		if call_first_handler {
			firstHandler(vctx, b, update)
			return
		}
		next(vctx, b, update)
	}
}
