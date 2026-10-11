-- name: GetUserByID :one
SELECT id, user_name, created_at, display_name, updated_at
FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT id, user_name, created_at, display_name, updated_at
FROM users
ORDER BY created_at DESC;

-- name: CreateUser :one
INSERT INTO users (
    id, user_name, display_name
) VALUES (
    $1, $2, $3
)
RETURNING id, user_name, created_at, display_name, updated_at;

-- name: UpdateUser :one
UPDATE users
SET user_name = $2,
    display_name = $3,
    updated_at = now()
WHERE id = $1
RETURNING id, user_name, created_at, display_name, updated_at;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: GetDeviceByID :one
SELECT id, user_id, name, memo, created_at
FROM devices
WHERE id = $1 LIMIT 1;

-- name: ListDevicesByUserID :many
SELECT id, user_id, name, memo, created_at
FROM devices
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: CreateDevice :one
INSERT INTO devices (
    id, user_id, name, memo
) VALUES (
    $1, $2, $3, $4
)
RETURNING id, user_id, name, memo, created_at;

-- name: UpdateDevice :one
UPDATE devices
SET name = $2,
    memo = $3
WHERE id = $1
RETURNING id, user_id, name, memo, created_at;

-- name: DeleteDevice :exec
DELETE FROM devices
WHERE id = $1;

-- name: GetRecordByID :one
SELECT id, device_id, type, ip_address, reported_at
FROM records
WHERE id = $1 LIMIT 1;

-- name: ListRecordsByDeviceID :many
SELECT id, device_id, type, ip_address, reported_at
FROM records
WHERE device_id = $1
ORDER BY reported_at DESC;

-- name: CreateRecord :one
INSERT INTO records (
    id, device_id, type, ip_address
) VALUES (
    $1, $2, $3, $4
)
RETURNING id, device_id, type, ip_address, reported_at;

-- name: UpdateRecord :one
UPDATE records
SET type = $2,
    ip_address = $3
WHERE id = $1
RETURNING id, device_id, type, ip_address, reported_at;

-- name: DeleteRecord :exec
DELETE FROM records
WHERE id = $1;

-- name: GetCredentialByUserIDAndProvider :one
SELECT user_id, provider, secret_hash, subject, created_at, updated_at
FROM credentials
WHERE user_id = $1 AND provider = $2 LIMIT 1;

-- name: ListCredentialsByUserID :many
SELECT user_id, provider, secret_hash, subject, created_at, updated_at
FROM credentials
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: CreateCredential :one
INSERT INTO credentials (
    user_id, provider, secret_hash, subject
) VALUES (
    $1, $2, $3, $4
)
RETURNING user_id, provider, secret_hash, subject, created_at, updated_at;

-- name: UpdateCredential :one
UPDATE credentials
SET secret_hash = $3,
    subject = $4,
    updated_at = now()
WHERE user_id = $1 AND provider = $2
RETURNING user_id, provider, secret_hash, subject, created_at, updated_at;

-- name: DeleteCredential :exec
DELETE FROM credentials
WHERE user_id = $1 AND provider = $2;