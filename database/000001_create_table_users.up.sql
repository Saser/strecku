BEGIN;

CREATE TABLE IF NOT EXISTS users (
    PRIMARY KEY (uuid),     -- surrogate key
    UNIQUE (email_address), -- natural key
    uuid          UUID    NOT NULL,
    deleted       BOOLEAN NOT NULL,
    email_address TEXT    NOT NULL,
                          CONSTRAINT email_address_not_empty
                          CHECK (email_address <> ''),
    display_name  TEXT    NOT NULL,
                          CONSTRAINT display_name_not_empty
                          CHECK (display_name <> '')
);

COMMIT;
