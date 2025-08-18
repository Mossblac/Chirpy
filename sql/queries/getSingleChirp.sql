-- name: GetSingleChirpByID :one
SELECT * FROM chirps
WHERE $1 = id;