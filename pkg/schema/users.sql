-- name: CreateUser :one
INSERT INTO users (
        name
) VALUES ($1) RETURNING *;

-- name: FindUserByID :one
SELECT
    *
FROM
    users
WHERE
    id = $1;

-- name: UpdateUser :one
UPDATE
    users
SET
    name = $1
WHERE
    id = $2
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM
    users
WHERE
    id = $1;
