-- name: FindAssistantByPhone :one
SELECT * FROM assistants
WHERE phone = $1;

-- name: FindAssistantByID :one
SELECT * FROM assistants
WHERE id = $1;

-- name: SetAssistantPassword :exec
UPDATE assistants
SET password = $1
WHERE id = $2;