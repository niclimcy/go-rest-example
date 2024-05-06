-- name: GetUserById :one
SELECT * FROM user_account
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM user_account
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM user_account;

-- name: CreateUser :one
INSERT INTO user_account (
  name, email, phone_number, address, user_type
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UpdateUser :exec
UPDATE user_account
  set name = $2,
  email = $3,
  phone_number = $4,
  address = $5,
  user_type = $6
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM user_account
WHERE id = $1;