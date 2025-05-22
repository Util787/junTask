-- name: GetAllUsers :many
SELECT *
FROM users
WHERE 
  (@name::text IS NULL OR name ILIKE '%' || @name || '%') AND
  (@surname::text IS NULL OR surname ILIKE '%' || @surname || '%') AND
  (@patronymic::text IS NULL OR patronymic ILIKE '%' || @patronymic || '%') AND
  (@gender::text IS NULL OR gender = @gender) 
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
