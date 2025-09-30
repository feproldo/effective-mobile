package dto

import (
	db "github.com/feproldo/effective-mobile/internal/db/generated"
)

type Subscription struct {
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date"`
}

func FromSql(subSql db.Subscription) Subscription {
	temp := Subscription{
		ServiceName: subSql.ServiceName,
		Price:       int(subSql.Price),
		UserID:      subSql.UserID.String(),
		StartDate:   subSql.StartDate.Format("01-2006"),
		EndDate:     nil,
	}

	if subSql.EndDate.Valid == true {
		endDate := subSql.EndDate.Time
		endDateParsed := endDate.Format("01-2006")
		temp.EndDate = &endDateParsed
	}

	return temp
}
