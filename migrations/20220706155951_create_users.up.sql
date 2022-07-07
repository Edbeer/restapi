DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS news CASCADE;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE TABLE users
(
    user_id      UUID PRIMARY KEY            DEFAULT uuid_generate_v4(),
    first_name   VARCHAR(32)                 NOT NULL check ( first_name <> '' ),
    last_name    VARCHAR(32)                 NOT NULL check ( last_name <> '' ),
    email        VARCHAR(64)                 NOT NULL check ( email <> '' ),
    password     VARCHAR(250)                NOT NULL CHECK ( octet_length(password) <> 0 ),
    role         VARCHAR(10)                 NOT NULL DEFAULT 'user',
    avatar       VARCHAR(512),
    phone_number VARCHAR(20),
    address      VARCHAR(250),
    city         VARCHAR(30),
    country      VARCHAR(30),
    postcode     SMALLINT,
    created_at   TIMESTAMP                   NOT NULL DEFAULT now(),
    updated_at   TIMESTAMP                            DEFAULT current_timestamp
);

CREATE TABLE news
(
    news_id    UUID PRIMARY KEY      DEFAULT uuid_generate_v4(),
    author_id  UUID         NOT NULL REFERENCES users (user_id),
    title      VARCHAR(250) NOT NULL check ( title <> '' ),
    content    TEXT         NOT NULL check ( content <> '' ),
    image_url  VARCHAR(1024) check ( image_url <> '' ),
    category   VARCHAR(250),
    created_at TIMESTAMP    NOT NULL DEFAULT now(),
    updated_at TIMESTAMP             DEFAULT current_timestamp
);