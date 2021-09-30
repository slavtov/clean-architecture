DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS articles CASCADE;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


CREATE TABLE users (
    id          uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    email       varchar(50) UNIQUE NOT NULL CHECK (email <> ''),
    "password"  varchar(250) NOT NULL CHECK ("password" <> ''),
    updated_at  timestamp with time zone NOT NULL DEFAULT current_timestamp,
    created_at  timestamp with time zone NOT NULL DEFAULT current_timestamp
);

CREATE TABLE articles (
    id          uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    author_id   uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE,
    title       varchar(250) NOT NULL CHECK (title <> ''),
    "desc"      text NOT NULL CHECK ("desc" <> ''),
    updated_at  timestamp with time zone NOT NULL DEFAULT current_timestamp,
    created_at  timestamp with time zone NOT NULL DEFAULT current_timestamp
);
