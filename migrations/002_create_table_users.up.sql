CREATE TABLE users (
    uuid UUID
        PRIMARY KEY,
    email NONEMPTYTEXT
        NOT NULL
        UNIQUE,
    password_hash NONEMPTYTEXT
        NOT NULL
);
