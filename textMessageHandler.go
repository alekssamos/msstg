package main

import (
	"context"
	"fmt"
	"os"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func voiceMessageHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	text := update.Message.Text
	writeMedia := fmt.Sprintf("%d_%d.mp3", update.Message.Chat.ID, update.Message.ID)
	user := USER(ctx)
	_, err := b.SendChatAction(ctx, &bot.SendChatActionParams{ChatID: update.Message.Chat.ID, Action: models.ChatActionRecordVoice})
	LogError(err)
	if !Speak(text, user.VoiceName, writeMedia, user.VoiceRate, user.VoicePitch, 0) {
		os.Remove(writeMedia)
		text = "Произошла ошибка."
		p := &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: text}
		_, err = b.SendMessage(ctx, p)
		LogError(err)
		return
	}
	f, err := os.OpenFile(writeMedia, os.O_RDONLY, 0644)
	LogError(err)
	defer (func() {
		f.Close()
		os.Remove(writeMedia)
	})()
	duration, err := DurationMp3(writeMedia, BITRATE)
	LogError(err)
	_, err = b.SendChatAction(ctx, &bot.SendChatActionParams{ChatID: update.Message.Chat.ID, Action: models.ChatActionUploadDocument})
	LogError(err)
	p := &bot.SendAudioParams{ChatID: update.Message.Chat.ID, Duration: duration, Audio: &models.InputFileUpload{Filename: writeMedia, Data: f}}
	_, err = b.SendAudio(ctx, p)
	LogError(err)
}
