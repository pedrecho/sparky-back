package convert

import (
	"fmt"
	"net/url"
	"sparky-back/internal/models"
	"strconv"
	"time"
)

const timeLayout = "2006-01-02 15:04:05"

func FormToMessage(form url.Values) (*models.Message, error) {
	message := new(models.Message)
	var err error

	userID := form.Get("user_id")
	if userID != "" {
		*message.UserID, err = strconv.ParseInt(userID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse user_id field: %w", err)
		}
	}

	chatID := form.Get("chat_id")
	if userID != "" {
		message.ChatID, err = strconv.ParseInt(chatID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse chat_id field: %w", err)
		}
	}

	timeStr := form.Get("time")
	if timeStr != "" {
		message.Time, err = time.Parse(timeLayout, timeStr)
		if err != nil {
			return nil, fmt.Errorf("parse time field: %w", err)
		}
	}

	message.Text = form.Get("text")

	return message, nil
}
