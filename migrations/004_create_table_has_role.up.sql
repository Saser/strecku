CREATE TYPE user_role AS ENUM (
    'customer'
);

CREATE TABLE has_role (
    "user" UUID -- Quoted since `user` is a reserved keyword in SQL.
        NOT NULL
        REFERENCES users (uuid),
    store UUID
        NOT NULL
        REFERENCES stores (uuid),
    role USER_ROLE
        NOT NULL,

    PRIMARY KEY ("user", store)
);
