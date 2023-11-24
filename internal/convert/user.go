package convert

import (
	"fmt"
	"net/url"
	"sparky-back/internal/models"
	"strconv"
	"time"
)

const birthdayLayout = "2006-01-02"

func FormToUser(form url.Values) (*models.User, error) {
	user := new(models.User)
	var err error
	//TODO id
	user.Email = form.Get("email")
	user.Password = form.Get("password")
	user.Name = form.Get("name")

	birthday := form.Get("birthday")
	if birthday != "" {
		user.Birthday, err = time.Parse(birthdayLayout, birthday)
		if err != nil {
			return nil, fmt.Errorf("parse birthday field: %w", err)
		}
	}

	sex := form.Get("sex")
	if sex != "" {
		user.Sex, err = strconv.ParseBool(sex)
		if err != nil {
			return nil, fmt.Errorf("parse sex field: %w", err)
		}
	}

	latitude := form.Get("latitude")
	if latitude != "" {
		user.Latitude, err = strconv.ParseFloat(latitude, 64)
		if err != nil {
			return nil, fmt.Errorf("parse latitude field: %w", err)
		}
	}

	longitude := form.Get("longitude")
	if longitude != "" {
		user.Longitude, err = strconv.ParseFloat(longitude, 64)
		if err != nil {
			return nil, fmt.Errorf("parse longitude field: %w", err)
		}
	}

	return user, nil
}
