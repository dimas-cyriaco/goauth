-- name: InsertAccount :one
INSERT INTO accounts (email, hashed_password)
VALUES ($1, $2)
RETURNING *;

-- name: ByID :one
select * FROM accounts where id = $1;

-- name: ByEmail :one
select * FROM accounts where email = $1;

-- name: ListAccounts :many
SELECT * FROM accounts;

-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1;

-- name: MarkEmailAsVerified :exec
UPDATE accounts
set email_verified_at = CURRENT_TIMESTAMP
WHERE email_verified_at is NULL and id = $1;
