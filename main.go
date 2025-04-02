package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"

	"github.com/glebarez/sqlite"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"gorm.io/gorm"
)

var mu sync.Mutex

type ctxDbKey struct{}
type ctxUserKey struct{}

func main() {
	// context
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// db
	db, err := gorm.Open(sqlite.Open("dbmsstg.db"), &gorm.Config{})
	if err != nil {
		log.Println("failed to connect database " + err.Error())
		return
	}

	// Migrate the schema
	err = db.AutoMigrate(&User{}, &Book{}, &BookPart{})
	if err != nil {
		log.Println("failed to Migrate database " + err.Error())
		return
	}
	ctx = context.WithValue(ctx, ctxDbKey{}, db.WithContext(ctx))
	// telegram
	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
		bot.WithCallbackQueryDataHandler("selectvoice", bot.MatchTypePrefix, selectVoiceCallbackHandler),
		bot.WithCallbackQueryDataHandler("selectrate", bot.MatchTypePrefix, selectRateCallbackHandler),
		bot.WithCallbackQueryDataHandler("selectpitch", bot.MatchTypePrefix, selectPitchCallbackHandler),
		bot.WithCallbackQueryDataHandler("voice_", bot.MatchTypePrefix, selectedVoiceCallbackHandler),
		bot.WithCallbackQueryDataHandler("rate:", bot.MatchTypePrefix, selectedRatePitchCallbackHandler),
		bot.WithCallbackQueryDataHandler("pitch:", bot.MatchTypePrefix, selectedRatePitchCallbackHandler),
		bot.WithCallbackQueryDataHandler("currentvalue", bot.MatchTypePrefix, dummyCallbackHandler),
		bot.WithCallbackQueryDataHandler("btn_cancel", bot.MatchTypePrefix, cancelCallbackHandler),
		bot.WithCallbackQueryDataHandler("btn_ok", bot.MatchTypePrefix, okCallbackHandler),
		bot.WithMiddlewares(dbUser),
	}

	b, err := bot.New(os.Getenv("BOT_TOKEN"), opts...)
	if err != nil {
		log.Println(err.Error())
		return
	}

	fmt.Println("Starting bot...")
	_, err = b.SetMyCommands(ctx, &bot.SetMyCommandsParams{Commands: []models.BotCommand{{Command: "settings", Description: "Открыть настройки бота"}}})
	LogError(err)
	b.Start(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	text := update.Message.Text
	isCommand := true
	switch {
	case strings.HasPrefix(text, "/start"):
		commandStartHandler(ctx, b, update)
	case strings.HasPrefix(text, "/help"):
		commandStartHandler(ctx, b, update)
	case strings.HasPrefix(text, "/settings"):
		commandSettingsHandler(ctx, b, update)
	default:
		isCommand = false
		if len(text) > 0 {
			voiceMessageHandler(ctx, b, update)
		}
	}
	if isCommand {
		_, err := b.DeleteMessage(ctx, &bot.DeleteMessageParams{
			ChatID:    update.Message.Chat.ID,
			MessageID: update.Message.ID,
		})
		LogError(err)
	}
}
