CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    uuid VARCHAR(255),
    username VARCHAR(255),
    email VARCHAR(255),
    password VARCHAR(255),
    role VARCHAR(255),
    token VARCHAR(255)
);

INSERT INTO users (uuid, username, email, password, role) VALUES
    ('123e4567-e89b-12d3-a456-426614174000', 'mwazovzky', 'mike@example.com', 'secret', 'admin');