CREATE TABLE IF NOT EXISTS users (
    id SERIAL,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(80) NOT NULL,
    hashed_password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    -- Define CONSTRAINTS
    CONSTRAINT users_id_pk PRIMARY KEY (id),
    CONSTRAINT users_username_uk UNIQUE (username),
    CONSTRAINT users_email_uk UNIQUE (email)
);