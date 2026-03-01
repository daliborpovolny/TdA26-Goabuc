PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS user (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,

    first_name  TEXT NOT NULL,
    last_name   TEXT NOT NULL,
    hash        TEXT NOT NULL,
    email       TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS admin (
    user_id INTEGER PRIMARY KEY,

    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
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
    updated_at INTEGER NOT NULL,

    archived INTEGER NOT NULL DEFAULT 0,
    state TEXT NOT NULL DEFAULT 'preparation' -- preparation | open | closed


);

CREATE TABLE IF NOT EXISTS module (
    uuid TEXT PRIMARY KEY,
    course_uuid TEXT NOT NULL,

    name TEXT NOT NULL,
    description TEXT NOT NULL,
    state TEXT NOT NULL DEFAULT 'preparation',  -- preparation | open | closed
    
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,

    FOREIGN KEY (course_uuid) REFERENCES course(uuid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS material_to_module (
    module_uuid TEXT NOT NULL,
    material_uuid TEXT NOT NULL,

    "order" INTEGER NOT NULL,

    FOREIGN KEY (material_uuid) REFERENCES material(uuid) ON DELETE CASCADE,
    FOREIGN KEY (module_uuid) REFERENCES module(uuid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS quiz_to_module (
    module_uuid TEXT NOT NULL,
    quiz_uuid TEXT NOT NULL,

    "order" INTEGER NOT NULL,

    FOREIGN KEY (quiz_uuid) REFERENCES quiz(uuid) ON DELETE CASCADE,
    FOREIGN KEY (module_uuid) REFERENCES module(uuid) ON DELETE CASCADE
);


-- heading could be in the future used as a general styling unit - 
-- ei variant of X means title, variant of Y means longer text, varint Z horizontal line etc...
-- variant Y red warning text balh blah so on, endless options....
CREATE TABLE IF NOT EXISTS heading (
    uuid TEXT PRIMARY KEY,
    course_uuid TEXT NOT NULL,

    content TEXT NOT NULL,
    variant TEXT NOT NULL, -- "heading" | ..... - this is just an idea for the future 

    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,

    FOREIGN KEY (course_uuid) REFERENCES course(uuid) ON DELETE CASCADE

);

CREATE TABLE IF NOT EXISTS heading_to_module (
    module_uuid TEXT NOT NULL,
    heading_uuid TEXT NOT NULL,

    "order" INTEGER NOT NULL,

    FOREIGN KEY (heading_uuid) REFERENCES heading(uuid) ON DELETE CASCADE,
    FOREIGN KEY (module_uuid) REFERENCES module(uuid) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS material (
    uuid TEXT PRIMARY KEY,
    course_uuid TEXT NOT NULL,
    
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    url TEXT NOT NULL,
    type TEXT NOT NULL,

    favicon_url TEXT,

    mime_type TEXT,
    byte_size INTEGER,

    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,
    
    FOREIGN KEY (course_uuid) REFERENCES course(uuid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS quiz (
    uuid TEXT PRIMARY KEY,
    course_uuid TEXT NOT NULL,


    title TEXT NOT NULL,
    attempts_count INTEGER NOT NULL,
    
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,
    
    FOREIGN KEY (course_uuid) REFERENCES course(uuid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS question (
    uuid TEXT PRIMARY KEY,
    quiz_uuid TEXT NOT NULL,

    question_order INTEGER NOT NULL,

    type TEXT NOT NULL,
    question_text TEXT NOT NULL,
    options TEXT NOT NULL,
    correct_indices TEXT NOT NULL,
    
    FOREIGN KEY (quiz_uuid) REFERENCES quiz(uuid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS answer (
    quiz_uuid TEXT NOT NULL,
    comment TEXT,

    score INTEGER NOT NULL,
    max_score INTEGER NOT NULL,

    user_id INTEGER,
    attempt_number INTEGER NOT NULL,

    submitted_at INTEGER NOT NULL,

    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
    FOREIGN KEY (quiz_uuid) REFERENCES quiz(uuid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS feed_posts (
    uuid TEXT PRIMARY KEY,
    course_uuid TEXT NOT NULL,

    type TEXT NOT NULL, -- 'manual' or 'system'
    message TEXT NOT NULL,
    is_edited BOOLEAN NOT NULL DEFAULT 0,
    
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,

    FOREIGN KEY (course_uuid) REFERENCES course(uuid) ON DELETE CASCADE
);