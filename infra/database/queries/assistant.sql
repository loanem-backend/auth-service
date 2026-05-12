-- name: FindAssistantByPhone :one
SELECT * FROM assistants
WHERE phone = $1;