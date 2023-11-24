package convert

import (
	"fmt"
	"net/url"
	"sparky-back/internal/models"
	"strconv"
	"time"
)

func FormToUser(form url.Values) (*models.User, error) {
	user := new(models.User)
	var err error
	//TODO id
	user.Email = form.Get("email")
	user.Password = form.Get("password")
	user.Name = form.Get("name")
	birthday := form.Get("birthday")
	user.Birthday, err = time.Parse("2006-01-02", birthday)
	if err != nil {
		return nil, fmt.Errorf("parse birthday field: %w", err)
	}
	sex := form.Get("sex")
	user.Sex, err = strconv.ParseBool(sex)
	if err != nil {
		return nil, fmt.Errorf("parse sex field: %w", err)
	}

	return user, nil
}
