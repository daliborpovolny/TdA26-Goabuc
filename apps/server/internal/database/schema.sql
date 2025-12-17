PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS user (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    first_name  TEXT NOT NULL,
    last_name   TEXT NOT NULL,
    hash        TEXT NOT NULL,
    email       TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS session (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id     INTEGER NOT NULL,
    token      TEXT NOT NULL,
    created_at  INTEGER NOT NULL,
    expires_at  INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS course (
    uuid TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS material (
    uuid TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    url TEXT NOT NULL,
    courseUuid TEXT NOT NULL,
    FOREIGN KEY (courseUuid) REFERENCES course(uuid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS quizz (
    uuid TEXT PRIMARY KEY,
    courseUuid TEXT NOT NULL,
    title TEXT NOT NULL,
    attemptsCount INTEGER NOT NULL,
    FOREIGN KEY (courseUuid) REFERENCES course(uuid) ON DELETE CASCADE
)

CREATE TABLE IF NOT EXISTS  question (
    uuid TEXT PRIMARY KEY,
    quizzUuid TEXT NOT NULL,
    type TEXT NOT NULL,
    question_text TEXT NOT NULL,
    options TEXT NOT NULL,
    correct_indices TEXT NOT NULL,
    FOREIGN KEY (quizzUuid) REFERENCES quizz(uuid) ON DELETE CASCADE
)