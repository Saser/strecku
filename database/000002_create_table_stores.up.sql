BEGIN;

CREATE TABLE IF NOT EXISTS stores (
    PRIMARY KEY (uuid),
    uuid         UUID    NOT NULL,
    deleted      BOOLEAN NOT NULL,
    display_name TEXT    NOT NULL,
                         CONSTRAINT display_name_not_empty
                         CHECK (display_name <> '')
);

COMMIT;
