CREATE TABLE users
(
    username TEXT PRIMARY KEY UNIQUE NOT NULL,
    password TEXT                    NOT NULL
);