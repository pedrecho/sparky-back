package convert

import (
	"fmt"
	"net/url"
	"sparky-back/internal/models"
	"strconv"
)

func FormToReaction(form url.Values) (*models.Reaction, error) {
	reaction := new(models.Reaction)
	var err error

	userID := form.Get("user_id")
	if userID != "" {
		reaction.UserID, err = strconv.ParseInt(userID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse user_id field: %w", err)
		}
	}

	toID := form.Get("to_id")
	if toID != "" {
		reaction.ToID, err = strconv.ParseInt(toID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse to_id field: %w", err)
		}
	}

	like := form.Get("like")
	if like != "" {
		reaction.Like, err = strconv.ParseBool(like)
		if err != nil {
			return nil, fmt.Errorf("parse like field: %w", err)
		}
	}

	return reaction, nil
}
