CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


CREATE TABLE articles (
    id          uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    author_id   uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE,
    title       varchar(250) NOT NULL CHECK (title <> ''),
    "desc"      text NOT NULL CHECK ("desc" <> ''),
    updated_at  timestamp with time zone NOT NULL DEFAULT current_timestamp,
    created_at  timestamp with time zone NOT NULL DEFAULT current_timestamp
);
