CREATE EXTENSION IF NOT EXISTS citext;
CREATE TABLE IF NOT EXISTS carts (
user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
books text[] not null,
book_id bigint NOT NULL REFERENCES books ON DELETE CASCADE,
quantity int not null,
total_price int not null)