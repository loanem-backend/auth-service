-- name: FindAssistantByPhone :one
SELECT * FROM assistants
WHERE phone = $1;

-- name: FindAssistantByID :one
SELECT * FROM assistants
WHERE id = $1;

-- name: InsertAssistant :one
INSERT INTO assistants (name, phone, password, period)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: SetAssistantPassword :exec
UPDATE assistants
SET password = $1, updated_at = $2
WHERE id = $3;