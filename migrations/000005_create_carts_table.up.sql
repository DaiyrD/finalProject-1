CREATE EXTENSION IF NOT EXISTS citext;
CREATE TABLE IF NOT EXISTS carts (
    id bigserial not null primary key,
email citext NOT NULL,
books text[] not null,
book_id bigint NOT NULL REFERENCES books ON DELETE CASCADE,
quantity int not null,
total_quantity int not null,
total_price int not null)