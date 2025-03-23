package main

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot/models"
)

var CallbackCancelButton []models.InlineKeyboardButton

func init() {
	CallbackCancelButton = []models.InlineKeyboardButton{{Text: "Отмена", CallbackData: "btn_cancel"}}
}

/*
func BuildStartKeyboard() models.ReplyKeyboardMarkup {
	return models.ReplyKeyboardMarkup{Keyboard: [][]models.KeyboardButton{
		{
			{Text: "/settings"},
//			{Text: "два"},
//			{Text: "три"},
		},
//		{
//			{Text: "четыре"},
//			{Text: "пять"},
//			{Text: "шесть"},
//		},
	}}
}
*/

func BuildSettingsKeyboard(ctx context.Context) *models.InlineKeyboardMarkup {
	u := USER(ctx)
	vn := "_"
	v, err := FindVoices(u.VoiceName)
	LogError(err)
	if err == nil {
		vn = v[0].OnlyName()
	}
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: fmt.Sprintf("Выбрать голос \t %s", vn), CallbackData: "selectvoice"},
			},
			CallbackCancelButton,
		},
	}
}
