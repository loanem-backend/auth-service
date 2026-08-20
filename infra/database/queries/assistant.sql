-- name: FindAssistantByPhone :one
SELECT * FROM assistants
WHERE phone = $1;

-- name: FindAssistantByID :one
SELECT * FROM assistants
WHERE id = $1;

-- name: InsertAssistant :one
INSERT INTO assistants (name, phone, password, period, email)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: SetAssistantPassword :exec
UPDATE assistants
SET password = $1, updated_at = $2
WHERE id = $3;

-- name: FindActiveAssistants :many
SELECT * FROM assistants
WHERE active = TRUE
ORDER BY period, name;

-- name: DeleteAssistantByID :exec
DELETE FROM assistants
WHERE id = $1;
