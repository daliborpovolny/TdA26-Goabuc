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

--* Session

-- name: CreateSession :one
INSERT INTO session (
    user_id, token, created_at, expires_at
) VALUES (
    ?, ?, ?, ?
) RETURNING *;

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

--* Material

-- name: CreateMaterial :one
INSERT INTO material (
    uuid, name, description, url, courseUuid
) VALUES (
    ?, ?, ?, ?, ?
) RETURNING *;

-- name: UpdateMaterial :one
UPDATE material SET name = ?, description = ?, url = ? WHERE material.uuid = ? RETURNING *;

-- name: DeleteMaterial :execresult
DELETE FROM material WHERE material.uuid = ?;

-- name: GetMaterial :one
SELECT uuid, name, description, url, courseUuid FROM material WHERE material.uuid = ?;

-- name: ListAllMaterialsOfCourse :many
SELECT uuid, name, description, url, courseUuid FROM material WHERE material.courseUuid = ?;