ALTER TABLE books
ADD CONSTRAINT genres_length_check CHECK (
        array_length(genres, 1) BETWEEN 1 AND 5
    );