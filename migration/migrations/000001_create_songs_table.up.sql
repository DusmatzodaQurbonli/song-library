CREATE TABLE songs (
    id SERIAL PRIMARY KEY,
    group_name TEXT NOT NULL,
    title TEXT NOT NULL,
    release_date TIMESTAMP NOT NULL,
    text TEXT NOT NULL,
    link TEXT NOT NULL
);