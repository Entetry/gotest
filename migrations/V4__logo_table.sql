CREATE TABLE logo
(
    id         uuid  NOT NULL PRIMARY KEY,
    company_id uuid  NOT NULL,
    image      bytea NOT NULL
);
