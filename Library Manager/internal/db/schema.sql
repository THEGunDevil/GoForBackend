-- db/schema.sql
CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    year INT,
    is_borrowed BOOLEAN DEFAULT FALSE
);
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL CHECK (email = LOWER(email))
);

-- Create case-insensitive unique index
CREATE UNIQUE INDEX unique_email_lower_idx ON users (LOWER(email));



CREATE TABLE borrows (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    book_id INT REFERENCES books(id) ON DELETE CASCADE,
    borrowed_at TIMESTAMP DEFAULT NOW(),
    returned_at TIMESTAMP
);
