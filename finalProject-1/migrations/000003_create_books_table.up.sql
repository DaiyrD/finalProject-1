-- null values are not appreciated in GoLang
-- so all columns either not null or have default vals
CREATE TABLE IF NOT EXISTS books (
    -- id column is a 64-bit auto-incrementing integer & primary key
                                      id bigserial PRIMARY KEY,
                                      created_at timestamp(0) with time zone not null default NOW(),
    title text not null,
    year integer not null,
    author text not null,
    genres text[] not null,
    price integer not null,
    version integer not null default 1
    );