package convert

import (
	"fmt"
	"net/url"
	"sparky-back/internal/models"
	"strconv"
)

func FormToFilter(form url.Values) (*models.Filter, error) {
	filter := new(models.Filter)
	var err error

	userID := form.Get("user_id")
	if userID != "" {
		filter.UserID, err = strconv.ParseInt(userID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse user_id field: %w", err)
		}
	}

	sex := form.Get("sex")
	if sex != "" {
		filter.Sex, err = strconv.ParseBool(sex)
		if err != nil {
			return nil, fmt.Errorf("parse sex field: %w", err)
		}
	}

	minAge := form.Get("min_age")
	if minAge != "" {
		filter.MinAge, err = strconv.Atoi(minAge)
		if err != nil {
			return nil, fmt.Errorf("parse min_age field: %w", err)
		}
	}

	maxAge := form.Get("max_age")
	if minAge != "" {
		filter.MaxAge, err = strconv.Atoi(maxAge)
		if err != nil {
			return nil, fmt.Errorf("parse max_age field: %w", err)
		}
	}

	distance := form.Get("distance")
	if distance != "" {
		filter.Distance, err = strconv.ParseFloat(distance, 64)
		if err != nil {
			return nil, fmt.Errorf("parse distance field: %w", err)
		}
	}

	limit := form.Get("limit")
	if limit != "" {
		filter.Limit, err = strconv.Atoi(limit)
		if err != nil {
			return nil, fmt.Errorf("parse limit field: %w", err)
		}
	}

	return filter, nil
}
