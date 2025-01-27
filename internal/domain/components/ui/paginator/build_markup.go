package paginator

import (
	models2 "ai-bot/internal/domain/models"
	"github.com/go-telegram/bot/models"
)

func (p *Paginator) buildMarkup(data []models2.User) models.InlineKeyboardMarkup {
	markup := append(p.buildList(data), p.buildKeyboard().InlineKeyboard...)
	return models.InlineKeyboardMarkup{InlineKeyboard: markup}
}
