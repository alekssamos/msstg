package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
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
	Не забудьте выбрать нужный голос.
	`
	fmt.Println("ctx:", ctx)
	kb := BuildSettingsKeyboard(ctx)
	_, err := b.SendChatAction(ctx, &bot.SendChatActionParams{ChatID: update.Message.Chat.ID, Action: models.ChatActionTyping})
	LogError(err)
	p := &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: text, ReplyMarkup: kb}
	_, err = b.SendMessage(ctx, p)
	LogError(err)
}

const UseButtons = "Используйте кнопки"

func commandSettingsHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	text := UseButtons
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

func commandFindVoiceHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// Команда /findvoice <имя_голоса>
	// Ищем голос по имени через findvoice
	text := "Укажите имя голоса"
	command := strings.SplitN(update.Message.Text, " ", 2)
	if len(command) != 2 {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   text,
		})
		LogError(err)
		return
	}
	voiceName := strings.TrimSpace(command[1])
	voiceName = voiceTranslit(voiceName)
	// ищем голос
	voices, err := FindVoices(voiceName)
	if err != nil {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Голос не найден",
		})
		LogError(err)
		return
	}
	// Если найдено несколько, то выводим но всё равно выбираем первый
	switch len(voices) {
	case 0:
		text = "Голос не найден"
	case 1:
		text = fmt.Sprintf("Найден голос: %s (%s)\n", voices[0].FriendlyName, voices[0].ShortName)
	default:
		text = "Найдено несколько голосов:\n"
		for _, v := range voices {
			text += fmt.Sprintf("%s (%s)\n", v.FriendlyName, v.ShortName)
		}
	}
	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   text,
	})
	LogError(err)
	// Выбираем первый голос
	voice, ok := voices[0], len(voices) > 0
	if !ok {
		return
	}
	// Сохраняем голос в БД
	u := USER(ctx)
	DB(ctx).Model(&u).Update("VoiceName", voice.ShortName)
	// Отправляем сообщение об успешном изменении голоса
	text = fmt.Sprintf("Голос успешно изменён: %s", voice.FriendlyName)
	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   text,
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

func selectRateCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	LogError(err)
	kb := BuildAdjustmentKeyboard(ctx, KeyboardRate)
	_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		MessageID:   update.CallbackQuery.Message.Message.ID,
		Text:        "Выбор скорости речи",
		ReplyMarkup: kb,
	})
	LogError(err)
}

func selectPitchCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	LogError(err)
	kb := BuildAdjustmentKeyboard(ctx, KeyboardPitch)
	_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		MessageID:   update.CallbackQuery.Message.Message.ID,
		Text:        "Выбор высоты голоса",
		ReplyMarkup: kb,
	})
	LogError(err)
}

func selectedRatePitchCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	showText := ""
	showAlert := false
	defer (func() {
		_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			Text:            showText,
			ShowAlert:       showAlert,
		})
		LogError(err)
	})()
	values := strings.Split(update.CallbackQuery.Data, ":")
	if len(values) != 2 {
		log.Printf("no rate or pitch in update data %q\n", update.CallbackQuery.Data)
		return
	}
	v, err := strconv.Atoi(values[1])
	LogError(err)
	u := USER(ctx)
	var field string
	var sum int
	var kbt adjustmentType
	switch values[0] {
	case "rate":
		field = "VoiceRate"
		sum = u.VoiceRate + v
		kbt = KeyboardRate
	case "pitch":
		field = "VoicePitch"
		sum = u.VoicePitch + v
		kbt = KeyboardPitch
	default:
		panic("unknown handler type")
	}
	if sum > 100 || sum < -100 {
		showAlert = true
		showText = "Значение может быть в диапазоне от  -100 до +100"
		return
	}
	DB(ctx).Model(&u).Update(field, sum)
	kb := BuildAdjustmentKeyboard(ctx, kbt)
	_, err = b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		MessageID:   update.CallbackQuery.Message.Message.ID,
		ReplyMarkup: kb,
	})
	LogError(err)
}

func okCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	LogError(err)
	kb := BuildSettingsKeyboard(ctx)
	_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		MessageID:   update.CallbackQuery.Message.Message.ID,
		Text:        UseButtons,
		ReplyMarkup: kb,
	})
	LogError(err)
}

func dummyCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	LogError(err)
}
