package dto

import (
	"database/sql"
	"time"

	db "github.com/feproldo/effective-mobile/internal/db/generated"
	"github.com/google/uuid"
)

const TIME_FORMAT = "01-2006"

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
		StartDate:   subSql.StartDate.Format(TIME_FORMAT),
		EndDate:     nil,
	}

	if subSql.EndDate.Valid == true {
		endDate := subSql.EndDate.Time
		endDateParsed := endDate.Format(TIME_FORMAT)
		temp.EndDate = &endDateParsed
	}

	return temp
}

func (sub *Subscription) ToSql() (*db.Subscription, error) {
	userUUID, err := uuid.Parse(sub.UserID)
	if err != nil {
		return nil, err
	}

	startDate, err := time.Parse(TIME_FORMAT, sub.StartDate)
	if err != nil {
		return nil, err
	}

	var endDate sql.NullTime = sql.NullTime{Valid: false, Time: time.Now()}

	if sub.EndDate != nil {

		endTime, err := time.Parse(TIME_FORMAT, sub.StartDate)
		if err == nil {
			endDate = sql.NullTime{
				Valid: true,
				Time:  endTime,
			}
		}
	}

	temp := db.Subscription{
		ServiceName: sub.ServiceName,
		Price:       int32(sub.Price),
		UserID:      userUUID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	return &temp, nil
}
