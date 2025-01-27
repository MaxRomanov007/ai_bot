package paginator

import (
	models2 "ai-bot/internal/domain/models"
	"fmt"
	"github.com/go-telegram/bot/models"
)

func (p *Paginator) buildList(data []models2.User) [][]models.InlineKeyboardButton {
	list := make([][]models.InlineKeyboardButton, len(data))
	for i := 0; i < len(data); i++ {
		list[i] = make([]models.InlineKeyboardButton, 1)
		list[i][0] = models.InlineKeyboardButton{
			Text:         fmt.Sprintf("%d. %s", data[i].UserID, data[i].Username),
			CallbackData: fmt.Sprintf("%s%s.%X", p.prefix, cmdUser, data[i].UserID),
		}
	}
	return list
}
