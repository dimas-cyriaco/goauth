-- name: InsertAccount :one
INSERT INTO accounts (email, hashed_password)
VALUES ($1, $2)
RETURNING *;

-- name: FindAccountByID :one
select * FROM accounts where id = $1;

-- name: FindAccountByEmail :one
select * FROM accounts where email = $1;

-- name: ListAccounts :many
SELECT * FROM accounts;

-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1;

-- name: MarkAccountEmailAsVerified :exec
UPDATE accounts
set email_verified_at = CURRENT_TIMESTAMP
WHERE email_verified_at is NULL and id = $1;

-- name: InsertSession :one
INSERT INTO sessions (account_id, user_agent, ip_address)
VALUES ($1, $2, $3)
RETURNING *;

-- name: FindSessionByID :one
select * FROM sessions where id = $1;
