package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func commandStartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	text := `
	Бот озвучит голосовые сообщения синтезатором речи от Microsoft.
	`
	_, err := b.SendChatAction(ctx, &bot.SendChatActionParams{ChatID: update.Message.Chat.ID, Action: models.ChatActionTyping})
	LogError(err)
	p := &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: text}
	_, err = b.SendMessage(ctx, p)
	LogError(err)
}

func firstHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	text := `
	Ты здесь впервые, да? Ща всё настроим...
	`
	_, err := b.SendChatAction(ctx, &bot.SendChatActionParams{ChatID: update.Message.Chat.ID, Action: models.ChatActionTyping})
	LogError(err)
	p := &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: text}
	_, err = b.SendMessage(ctx, p)
	LogError(err)
}

func commandSettingsHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	text := "Используйте кнопки"
	kb := BuildSettingsKeyboard(ctx)
	p := &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: text, ReplyMarkup: kb}
	_, err := b.SendMessage(ctx, p)
	LogError(err)
}

func cancelCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	LogError(err)
	_, err = b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		MessageID: update.CallbackQuery.Message.Message.ID,
	})
	LogError(err)
}

func selectedVoiceCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	responseText := ""
	defer (func() {
		_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			Text:            responseText,
			ShowAlert:       false,
		})
		LogError(err)
	})()
	v, ok := strings.CutPrefix(update.CallbackQuery.Data, "voice_")
	if !ok {
		log.Println("selected voice cut prefix error")
		return
	}
	voices, err := FindVoices(v)
	LogError(err)
	u := USER(ctx)
	DB(ctx).Model(&u).Update("VoiceName", voices[0].ShortName)
	responseText = fmt.Sprintf("Голос успешно изменён: %s", voices[0].FriendlyName)
	kb := BuildSettingsKeyboard(ctx)
	_, err = b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		MessageID:   update.CallbackQuery.Message.Message.ID,
		ReplyMarkup: kb,
	})
	LogError(err)
}

func selectVoiceCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	LogError(err)
	var buttons [][]models.InlineKeyboardButton
	col := []models.InlineKeyboardButton{{}}
	if strings.Contains(update.CallbackQuery.Data, "_") {
		country := strings.Split(update.CallbackQuery.Data, "_")[1]
		fmt.Printf("country: %q, data: %q;\n", country, update.CallbackQuery.Data)
		voices, err := FindVoices(country)
		LogError(err)
		for i, v := range voices {
			if i%3 == 0 && len(col) > 0 {
				buttons = append(buttons, col)
				col = nil
			}
			col = append(col, models.InlineKeyboardButton{Text: v.OnlyName(), CallbackData: fmt.Sprintf("voice_%s", v.ShortName)})
		}
	} else {
		countries := AllCountries()
		for i, the_item := range countries {
			if i%3 == 0 && len(col) > 0 {
				buttons = append(buttons, col)
				col = nil
			}
			col = append(col, models.InlineKeyboardButton{Text: the_item, CallbackData: fmt.Sprintf("selectvoice_%s", the_item)})
		}
	}
	// Добавляем оставшиеся кнопки, если есть при не кратных 3 количестве
	if len(col) > 0 {
		buttons = append(buttons, col)
	}
	buttons = append(buttons, CallbackCancelButton)
	kb := &models.InlineKeyboardMarkup{InlineKeyboard: buttons}
	_, err = b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		MessageID:   update.CallbackQuery.Message.Message.ID,
		ReplyMarkup: kb,
	})
	LogError(err)
}
