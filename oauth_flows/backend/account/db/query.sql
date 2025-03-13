-- name: InsertAccount :one
INSERT INTO accounts (email, hashed_password)
VALUES ($1, $2)
RETURNING *;

-- name: ByID :one
select * FROM accounts where id = $1;

-- name: ListAccounts :many
SELECT * FROM accounts;

-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1;
