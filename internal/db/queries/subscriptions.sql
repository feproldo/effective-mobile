-- name: SubscriptionsList :many
SELECT * FROM subscriptions;

-- name: CreateSubscription :exec
INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date) VALUES ($1, $2, $3, $4, $5);

-- name: UserSubscriptions :many
SELECT * FROM subscriptions WHERE user_id = $1;

-- name: GetSubscription :one
SELECT * FROM subscriptions WHERE id = $1;

-- name: DeleteSubscription :exec
DELETE FROM subscriptions WHERE id = $1;

-- name: UpdateSubscription :exec
UPDATE subscriptions SET service_name = $2, price = $3, user_id = $4, start_date = $5, end_date = $6 WHERE id = $1;

-- name: GetSubscriptionsWithFilter :many
SELECT price FROM subscriptions WHERE ($1 = '' OR user_id = $1::uuid) AND ($2 = '' OR service_name = $2::text) AND ($3 = '' OR start_date >= TO_DATE($3, 'MM-YYYY')) AND ($4 = '' OR end_date <= TO_DATE($4, 'MM-YYYY'));
