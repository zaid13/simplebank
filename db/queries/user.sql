-- name: CreateUsers :one
INSERT INTO users (
    username,
    hash_password,
    full_name,
    email,
    password_changed_at
) VALUES (
             $1, $2 , $3, $4, $5
         )
RETURNING *;

-- name: GetUsers :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: GetUsersForUpdate :one
SELECT * FROM users
WHERE username = $1 LIMIT 1
FOR NO KEY UPDATE ;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY username
LIMIT $1
OFFSET $2;

-- name: UpdateUsersName :one
UPDATE users
set full_name = $2
WHERE username = $1
RETURNING *;

-- name: UpdateUsersPassword :one
UPDATE users
set hash_password = $2 ,
    password_changed_at=$1
WHERE username = $1
RETURNING *;



-- name: DeleteUsers :exec
DELETE FROM users
WHERE username = $1;