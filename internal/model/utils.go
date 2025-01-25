package db_model

import "gorm.io/gorm"

// AddQueryParamsToDB adds ordering, pagination, and other common query parameters to a GORM query.
func AddQueryParamsToDB(query *gorm.DB, order string, limit, page, offset int) *gorm.DB {
	// Apply order if provided
	if order != "" {
		query = query.Order(order)
	}

	// Apply limit
	if limit > 0 {
		query = query.Limit(limit)
	}

	// Apply page and offset if provided
	if page > 0 {
		if offset == 0 {
			offset = (page - 1) * limit
		}
		query = query.Offset(offset)
	} else if offset > 0 {
		query = query.Offset(offset)
	}

	return query
}

// AuthPostRequestParams is the struct for the request body of the POST auth endpoint
type AuthPostRequestParams struct {
	UsernameOrEmail string `json:"username_or_email"`
	Password        string `json:"password"`
}
