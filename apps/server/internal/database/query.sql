--* USER

-- name: GetUser :one
SELECT * FROM user WHERE user.id = ?;

-- name: GetUserByEmail :one
SELECT * FROM user WHERE user.email = ?;

-- name: GetUserBySessionToken :one
SELECT * FROM user
JOIN session on user.id = session.user_id
WHERE session.token = ?;

-- name: CreateUser :one
INSERT INTO user (first_name, last_name, hash, email) VALUES (?, ?, ?, ?) RETURNING *;

-- * Admin

-- name: MakeUserAdmin :exec
INSERT INTO admin (user_id) VALUES (?);

--* Session

-- name: CreateSession :one
INSERT INTO session (
    user_id, token, created_at, expires_at
) VALUES (
    ?, ?, ?, ?
) RETURNING *;

-- name: InvalidateSession :exec
DELETE FROM session WHERE token = ?;

--* Course

-- name: CreateCourse :one
INSERT INTO course (
    uuid, name, description, created_at, updated_at
) VALUES (
    ?, ?, ?, ?, ?
) RETURNING *;

-- name: UpdateCourse :one
UPDATE course SET name = ?, description = ?, updated_at = ? WHERE course.uuid = ? RETURNING *;

-- name: DeleteCourse :execresult
DELETE FROM course WHERE course.uuid = ?;

-- name: GetCourse :one
SELECT uuid, name, description, created_at, updated_at FROM course WHERE course.uuid == ?;

-- name: ListAllCourses :many
SELECT uuid, name, description, created_at, updated_at FROM course;

-- name: CheckCourseExists :one
SELECT EXISTS (SELECT 1 FROM course WHERE uuid = ?) AS course_exists;

--* Material

-- name: CreateMaterial :one
INSERT INTO material (
    uuid, course_uuid, name, description, url, type, favicon_url, byte_size, mime_type, created_at, updated_at
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
) RETURNING *;

-- name: UpdateMaterial :one
UPDATE material SET name = ?, description = ?, url = ? WHERE material.uuid = ? RETURNING *;

-- name: DeleteMaterial :execresult
DELETE FROM material WHERE material.uuid = ?;

-- name: GetMaterial :one
SELECT * FROM material WHERE material.uuid = ?;

-- name: ListAllMaterialsOfCourse :many
SELECT * FROM material WHERE material.course_uuid = ?;

-- name: UpdateMaterialPartial :one
UPDATE material
SET
    name        = COALESCE(sqlc.narg(name), name),
    description = COALESCE(sqlc.narg(description), description),
    url         = COALESCE(sqlc.narg(url), url),
    favicon_url = COALESCE(sqlc.narg(favicon_url), favicon_url),
    byte_size   = COALESCE(sqlc.narg(byte_size), byte_size),
    mime_type   = COALESCE(sqlc.narg(mime_type), mime_type),
    updated_at  = sqlc.arg(updated_at)
WHERE uuid = sqlc.arg(uuid) RETURNING *;
