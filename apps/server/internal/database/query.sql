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
SELECT * FROM material WHERE material.course_uuid = ? ORDER BY created_at DESC;

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

--* Quizz

-- name: CreateQuizz :one
INSERT INTO quizz (
    uuid, course_uuid, title, attempts_count, created_at, updated_at
) VALUES (
    ?, ?, ?, ?, ?, ?
) RETURNING *;

-- name: UpdateQuizz :one
UPDATE quizz
SET
    title =             COALESCE(sqlc.narg(title), title),
    attempts_count =    COALESCE(sqlc.narg(attempts_count), attempts_count),
    updated_at =        COALESCE(sqlc.narg(updated_at), updated_at)
WHERE uuid = sqlc.arg(uuid)
RETURNING *;

-- name: DeleteQuizz :execresult
DELETE FROM quizz WHERE uuid = ?;

-- name: GetQuiz :many
SELECT
    qz.uuid AS quiz_uuid,
    qz.course_uuid AS course_uuid,
    qz.title AS quiz_title,
    qz.attempts_count AS quiz_attempts_count,
    qz.created_at AS quiz_created_at,
    qz.updated_at AS quiz_updated_at,

    qs.uuid AS question_uuid,
    qs.question_order AS question_order,
    qs.type AS question_type,
    qs.question_text AS question_text,
    qs.options AS question_options,
    qs.correct_indices AS question_correct_indices
FROM quizz qz
LEFT JOIN question qs
    ON qs.quizz_uuid = qz.uuid
WHERE qz.uuid = ?
ORDER BY qs.question_order;

-- name: ListQuizzes :many
SELECT
    qz.uuid AS quiz_uuid,
    qz.course_uuid AS course_uuid,
    qz.title AS quiz_title,
    qz.attempts_count AS quiz_attempts_count,
    qz.created_at AS quiz_created_at,
    qz.updated_at AS quiz_updated_at,

    qs.uuid AS question_uuid,
    qs.question_order AS question_order,
    qs.type AS question_type,
    qs.question_text AS question_text,
    qs.options AS question_options,
    qs.correct_indices AS question_correct_indices
FROM quizz qz
LEFT JOIN question qs
    ON qs.quizz_uuid = qz.uuid
ORDER BY qz.uuid ASC, qs.question_order ASC;
--* Question

-- name: CreateQuestion :one
INSERT INTO question (
    uuid, quizz_uuid, question_order, type, question_text, options, correct_indices
) SELECT 
    sqlc.arg(uuid),
    sqlc.arg(quiz_uuid),
    COALESCE(MAX("question_order"), 0) + 1,    
    sqlc.arg(type),
    sqlc.arg(question_text),
    sqlc.arg(options),
    sqlc.arg(correct_indices)
FROM question
WHERE quizz_uuid = sqlc.arg(quiz_uuid)
RETURNING *;

-- name: UpdateQuestion :one
UPDATE question
SET
    question_text = COALESCE(sqlc.narg(question_text), question_text),
    options = COALESCE(sqlc.narg(options), options),
    correct_indices = COALESCE(sqlc.narg(correct_indices), correct_indices)
RETURNING *;

-- name: GetQuestionsOfQuiz :many
SELECT * FROM question WHERE quizz_uuid = ?;

-- name: DeleteQuestionsOfQuiz :execresult
DELETE FROM question WHERE quizz_uuid = ?;

--* Answers

-- name: InsertAnswer :one
INSERT INTO answer (
    quiz_uuid, comment, score, max_score, user_id, attempt_number, submitted_at
) VALUES (
    sqlc.arg(quiz_uuid),
    sqlc.narg(comment),
    
    sqlc.arg(score),
    sqlc.arg(max_score),

    sqlc.narg(user_id),
    CASE
        WHEN sqlc.narg(user_id) IS NULL THEN 0
        ELSE (
            SELECT COALESCE(MAX(answer.attempt_number), 0) + 1
            FROM answer
                WHERE answer.quiz_uuid = sqlc.arg(quiz_uuid) 
                    AND answer.user_id = sqlc.narg(user_id)
        )
    END,

    sqlc.arg(submitted_at)
) RETURNING *;
