-- go-uni PostgreSQL schema

BEGIN;

CREATE TABLE IF NOT EXISTS students (
    id BIGSERIAL PRIMARY KEY,
    first_name TEXT NOT NULL CHECK (char_length(btrim(first_name)) > 0),
    last_name TEXT NOT NULL CHECK (char_length(btrim(last_name)) > 0),
    email TEXT NOT NULL UNIQUE CHECK (position('@' IN email) > 1)
);

CREATE TABLE IF NOT EXISTS teachers (
    id BIGSERIAL PRIMARY KEY,
    first_name TEXT NOT NULL CHECK (char_length(btrim(first_name)) > 0),
    last_name TEXT NOT NULL CHECK (char_length(btrim(last_name)) > 0),
    department TEXT NOT NULL CHECK (char_length(btrim(department)) > 0)
);

CREATE TABLE IF NOT EXISTS courses (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL CHECK (char_length(btrim(title)) > 0),
    description TEXT NOT NULL DEFAULT '',
    teacher_id BIGINT NOT NULL,
    CONSTRAINT fk_courses_teacher
        FOREIGN KEY (teacher_id)
        REFERENCES teachers(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS enrollments (
    student_id BIGINT NOT NULL,
    course_id BIGINT NOT NULL,
    PRIMARY KEY (student_id, course_id),
    CONSTRAINT fk_enrollments_student
        FOREIGN KEY (student_id)
        REFERENCES students(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    CONSTRAINT fk_enrollments_course
        FOREIGN KEY (course_id)
        REFERENCES courses(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

COMMIT;
