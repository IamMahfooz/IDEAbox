CREATE DATABASE ideabox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE ideabox;
CREATE TABLE snippets (
                          id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
                          title VARCHAR(100) NOT NULL,
                          content TEXT NOT NULL,
                          created DATETIME NOT NULL,
                          expires DATETIME NOT NULL,
                          byy VARCHAR(255) NOT NULL
);
CREATE TABLE users (
                       id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
                       name VARCHAR(255) NOT NULL,
                       email VARCHAR(255) NOT NULL,
                       hashed_password CHAR(60) NOT NULL,
                       created DATETIME NOT NULL,
                       active BOOLEAN NOT NULL DEFAULT TRUE
);
ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);
CREATE INDEX idx_snippets_created ON snippets(created);
CREATE USER 'web'@'localhost';
GRANT SELECT, INSERT, UPDATE ON ideabox.* TO 'web'@'localhost';
-- Important: Make sure to swap 'pass' with a password of your own choosing.
ALTER USER 'web'@'localhost' IDENTIFIED BY 'mahfooz';
