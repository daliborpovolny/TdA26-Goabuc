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
SELECT * FROM course WHERE course.uuid == ?;

-- name: ListAllCourses :many
SELECT * FROM course;

-- name: CheckCourseExists :one
SELECT EXISTS (SELECT 1 FROM course WHERE uuid = ?) AS course_exists;

-- name: ChangeCourseState :one
UPDATE course
SET
    state = ?,
    updated_at = ?
WHERE uuid = ? RETURNING *;


--* Module

-- name: CreateModule :one
INSERT INTO module (
    uuid, course_uuid, name, description, created_at, updated_at
) VALUES (
    ?, ?, ?, ?, ?, ?
) RETURNING *;


-- name: ChangeModuleState :one
UPDATE module
SET
    state = ?,
    updated_at = ?
WHERE uuid = ? RETURNING *;

-- name: CheckModuleExists :one
SELECT EXISTS (SELECT 1 FROM module WHERE uuid = ?) AS module_exists;

-- name: GetModule :one
SELECT * FROM module WHERE uuid = ? AND course_uuid = ?;

-- name: UpdateModule :one
UPDATE module
SET
    name = ?,
    description = ?
WHERE
    uuid = ? and course_uuid = ?
RETURNING *;

-- name: DeleteModule :exec
DELETE FROM module WHERE uuid = ? AND course_uuid = ?;

--* Heading

-- name: CreateHeading :one
INSERT INTO heading (
    uuid, content, created_at, updated_at
) VALUES (
    ?, ?, ?, ?
) RETURNING *;

-- name: GetHeading :one
SELECT * FROM heading WHERE uuid = ?;

-- name: UpdateHeading :one
UPDATE heading
SET
    content = ?,
    updated_at = ?
WHERE uuid = ?
RETURNING *;

-- name: DeleteHeading :exec
DELETE FROM heading WHERE uuid = ?;



--* Module to Others Pairings

-- HeadingToModule:

-- name: AssignHeadingToModule :one
INSERT INTO heading_to_module (
    module_uuid, heading_uuid, "order"
) VALUES (
    ?, ?, ?
) RETURNING *;

-- name: ChangeHeadingInModuleOrder :one
UPDATE heading_to_module
SET "order" = ?
WHERE heading_uuid = ? AND module_uuid = ?
RETURNING *;

-- name: RemoveHeadingFromModule :exec
DELETE FROM heading_to_module WHERE heading_uuid = ? AND module_uuid = ?;


-- MaterialToModule

-- name: AssignMaterialToModule :one
INSERT INTO material_to_module (
    module_uuid, material_uuid, "order"
) VALUES (
    ?, ?, ?
) RETURNING *;

-- name: ChangeMaterialInModuleOrder :one
UPDATE material_to_module
SET "order" = ?
WHERE material_uuid = ? AND module_uuid = ?
RETURNING *;

-- name: RemoveMaterialFromModule :exec
DELETE FROM material_to_module WHERE material_uuid = ? AND module_uuid = ?;


-- QuizToModule

-- name: AssignQuizToModule :one
INSERT INTO quiz_to_module (
    module_uuid, quiz_uuid, "order"
) VALUES (
    ?, ?, ?
) RETURNING *;

-- name: ChangeQuizInModuleOrder :one
UPDATE quiz_to_module
SET "order" = ?
WHERE quiz_uuid = ? AND module_uuid = ?
RETURNING *;

-- name: RemoveQuizFromModule :exec
DELETE FROM quiz_to_module WHERE quiz_uuid = ? AND module_uuid = ?;


--* ModuleContents

-- name: GetModuleContents :many
SELECT 
    'material' AS item_type,
    m.uuid,
    m.name AS display_name,
    m.description,
    mtm."order"
FROM material_to_module mtm
JOIN material m ON mtm.material_uuid = m.uuid
WHERE mtm.module_uuid = ? AND m.course_uuid = ?

UNION ALL

SELECT 
    'quiz' AS item_type,
    q.uuid,
    q.title AS display_name,
    '' AS description,
    qtm."order"
FROM quiz_to_module qtm
JOIN quiz q ON qtm.quiz_uuid = q.uuid
WHERE qtm.module_uuid = ? AND q.course_uuid = ?

UNION ALL

SELECT 
    'heading' AS item_type,
    h.uuid,
    h.content AS display_name,
    '' AS description,
    htm."order"
FROM heading_to_module htm
JOIN heading h ON htm.heading_uuid = h.uuid
WHERE htm.module_uuid = ? AND h.course_uuid = ?

ORDER BY "order" ASC;



--* Material

-- name: CreateMaterial :one
INSERT INTO material (
    uuid, course_uuid, name, description, url, type, favicon_url, byte_size, mime_type, created_at, updated_at
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
) RETURNING *;

-- name: UpdateMaterial :one
UPDATE material
SET 
    name = ?,
    description = ?,
    url = ?
WHERE material.uuid = ? RETURNING *;

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

--* Quiz

-- name: CreateQuiz :one
INSERT INTO quiz (
    uuid, course_uuid, title, attempts_count, created_at, updated_at
) VALUES (
    ?, ?, ?, ?, ?, ?
) RETURNING *;

-- name: UpdateQuiz :one
UPDATE quiz
SET
    title =             COALESCE(sqlc.narg(title), title),
    attempts_count =    COALESCE(sqlc.narg(attempts_count), attempts_count),
    updated_at =        COALESCE(sqlc.narg(updated_at), updated_at)
WHERE uuid = sqlc.arg(uuid)
RETURNING *;

-- name: IncrementQuizAttemptsCount :exec
UPDATE quiz
SET
    attempts_count = attempts_count + 1
WHERE uuid = ?;

-- name: DeleteQuiz :execresult
DELETE FROM quiz WHERE uuid = ?;

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
FROM quiz qz
LEFT JOIN question qs
    ON qs.quiz_uuid = qz.uuid
WHERE qz.uuid = ?
ORDER BY qs.question_order;

-- name: ListQuizes :many
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
FROM quiz qz
LEFT JOIN question qs
    ON qs.quiz_uuid = qz.uuid
WHERE qz.course_uuid = ?
ORDER BY qz.uuid ASC, qs.question_order ASC;


--* Question

-- name: CreateQuestion :one
INSERT INTO question (
    uuid, quiz_uuid, question_order, type, question_text, options, correct_indices
) SELECT 
    sqlc.arg(uuid),
    sqlc.arg(quiz_uuid),
    COALESCE(MAX("question_order"), 0) + 1,    
    sqlc.arg(type),
    sqlc.arg(question_text),
    sqlc.arg(options),
    sqlc.arg(correct_indices)
FROM question
WHERE quiz_uuid = sqlc.arg(quiz_uuid)
RETURNING *;

-- name: UpdateQuestion :one
UPDATE question
SET
    question_text = COALESCE(sqlc.narg(question_text), question_text),
    options = COALESCE(sqlc.narg(options), options),
    correct_indices = COALESCE(sqlc.narg(correct_indices), correct_indices)
RETURNING *;

-- name: GetQuestionsOfQuiz :many
SELECT * FROM question WHERE quiz_uuid = ?;

-- name: DeleteQuestionsOfQuiz :execresult
DELETE FROM question WHERE quiz_uuid = ?;

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

-- TODO RETURN ANSWERS AND SHOW THEM ON FRONTEND

-- name: GetAnswersOfQuiz :many
SELECT * FROM answer WHERE quiz_uuid = ? ORDER BY answer.submitted_at DESC;

--* Posts

-- name: GetPostsByCourse :many
SELECT * FROM feed_posts
WHERE course_uuid = ?
ORDER BY created_at DESC;

-- name: GetPost :one
SELECT * FROM feed_posts
WHERE uuid = ?;

-- name: CreatePost :one
INSERT INTO feed_posts (uuid, course_uuid, type, message, is_edited, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdatePost :one
UPDATE feed_posts
SET message = ?, is_edited = 1, updated_at = ?
WHERE uuid = ?
RETURNING *;

-- name: DeletePost :exec
DELETE FROM feed_posts
WHERE uuid = ?;