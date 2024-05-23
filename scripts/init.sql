-- Создание базы данных
CREATE DATABASE myDb;

-- Использование базы данных
\c myDb;

-- Создание таблицы для постов
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    comments_disabled BOOLEAN DEFAULT FALSE
);

-- Создание таблицы для комментариев
CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    post_id INTEGER NOT NULL,
    text TEXT NOT NULL,
    parent_id INTEGER,
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (parent_id) REFERENCES comments(id)
);