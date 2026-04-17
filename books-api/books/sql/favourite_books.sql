CREATE TABLE favorite_books (
    user_id    INT NOT NULL,
    book_id    INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, book_id)
);

CREATE INDEX idx_favorite_books_user_created_at
ON favorite_books (user_id, created_at DESC);