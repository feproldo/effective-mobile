package subscriptions

import (
	"context"
	"database/sql"
	"time"

	db "github.com/feproldo/effective-mobile/internal/db/generated"
	"github.com/feproldo/effective-mobile/internal/dto"
	"github.com/google/uuid"
)

type Services struct {
	queries *db.Queries
}

func NewService(queries *db.Queries) *Services {
	return &Services{
		queries: queries,
	}
}

func (s *Services) List(ctx context.Context) (*[]dto.Subscription, error) {
	list, err := s.queries.SubscriptionsList(ctx)
	if err != nil {
		return nil, err
	}

	var subs []dto.Subscription

	for _, el := range list {
		subs = append(subs, dto.FromSql(el))
	}
	return &subs, nil
}

func (s *Services) Create(ctx context.Context, sub dto.Subscription) error {
	userUUID, err := uuid.Parse(sub.UserID)
	if err != nil {
		return err
	}

	startDate, err := time.Parse("01-2006", sub.StartDate)
	if err != nil {
		return err
	}

	endDate := sql.NullTime{
		Time:  time.Now(),
		Valid: false,
	}

	if sub.EndDate != nil && *sub.EndDate != "" {
		endDateParsed, err := time.Parse("01-2006", *sub.EndDate)
		if err == nil {
			endDate.Time = endDateParsed
			endDate.Valid = true
		}
	}

	s.queries.CreateSubscription(ctx, db.CreateSubscriptionParams{
		ServiceName: sub.ServiceName,
		Price:       int32(sub.Price),
		UserID:      userUUID,
		StartDate:   startDate,
		EndDate:     endDate,
	})
	return nil
}

func (s *Services) GetByUserId(ctx context.Context, user_id uuid.UUID) (*[]dto.Subscription, error) {
	userUUID := user_id
	subsSql, err := s.queries.UserSubscriptions(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	var subs []dto.Subscription

	for _, el := range subsSql {
		subs = append(subs, dto.FromSql(el))
	}

	return &subs, nil
}

func (s *Services) Get(ctx context.Context, id int32) (*dto.Subscription, error) {
	subSql, err := s.queries.GetSubscription(ctx, id)
	if err != nil {
		return nil, err
	}

	sub := dto.FromSql(subSql)

	return &sub, nil
}

func (s *Services) Delete(ctx context.Context, id int32) error {
	err := s.queries.DeleteSubscription(ctx, id)
	return err
}

func (s *Services) Update(ctx context.Context, id int32, sub dto.Subscription) error {
	sqlSub, err := sub.ToSql()
	if err != nil {
		return err
	}

	updateSql := db.UpdateSubscriptionParams(*sqlSub)
	updateSql.ID = id

	err = s.queries.UpdateSubscription(ctx, updateSql)
	return err
}

func (s *Services) Sum(ctx context.Context, startDate string, endDate string, userId string, serviceName string) (*int, error) {
	sql := db.GetSubscriptionsWithFilterParams{
		Column1: userId,
		Column2: serviceName,
		Column3: startDate,
		Column4: endDate,
	}

	list, err := s.queries.GetSubscriptionsWithFilter(ctx, sql)
	if err != nil {
		return nil, err
	}

	sum := 0

	for _, el := range list {
		sum += int(el)
	}

	return &sum, nil
}
