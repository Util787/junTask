-- name: UserExists :one
SELECT EXISTS (
    SELECT 1
    FROM users
    WHERE name=$1 AND surname=$2 AND patronymic=$3 
);