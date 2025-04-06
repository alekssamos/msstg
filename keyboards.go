package main

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot/models"
)

var CallbackCancelButton []models.InlineKeyboardButton
var CallbackOKButton []models.InlineKeyboardButton

func init() {
	CallbackCancelButton = []models.InlineKeyboardButton{{Text: "Отмена", CallbackData: "btn_cancel"}}
	CallbackOKButton = []models.InlineKeyboardButton{{Text: "OK", CallbackData: "btn_ok"}}
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
				{Text: fmt.Sprintf("Выбрать голос \t 🗣 %s", vn), CallbackData: "selectvoice"},
			}, {
				{Text: fmt.Sprintf("Изменить скорость речи \t 🎯 %d", u.VoiceRate), CallbackData: "selectrate"},
			}, {
				{Text: fmt.Sprintf("Изменить высоту голоса \t 🎯 %d", u.VoicePitch), CallbackData: "selectpitch"},
			},
			CallbackCancelButton,
		},
	}
}

type adjustmentType int

const (
	KeyboardRate adjustmentType = iota + 1
	KeyboardPitch
)

func BuildAdjustmentKeyboard(ctx context.Context, what adjustmentType) *models.InlineKeyboardMarkup {
	u := USER(ctx)
	v := 0
	btn := ""
	switch what {
	case KeyboardRate:
		v = u.VoiceRate
		btn = "rate"
	case KeyboardPitch:
		v = u.VoicePitch
		btn = "pitch"
	default:
		panic("unknown keyboard type")
	}
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "➖ -10", CallbackData: fmt.Sprintf("%s:-10", btn)},
				{Text: fmt.Sprintf("🎯 %d", v), CallbackData: "currentvalue"},
				{Text: "➕ +10", CallbackData: fmt.Sprintf("%s:+10", btn)},
			},
			CallbackOKButton,
		},
	}
}
