-- name: UpdateUser :exec
UPDATE users
SET updated_at = $1, name = $2, surname = $3, patronymic = $4, age = $5, gender = $6, nationality = $7
WHERE id = $8;