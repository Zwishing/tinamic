CREATE TABLE users.user_info(
    id  serial PRIMARY KEY,
    uid uuid NOT NULL UNIQUE,
    name varchar(255) NOT NULL,
    create_at timestamptz,
    update_at timestamptz,
    password_hash varchar(255) NOT NULL,
    user_status int,
    user_role varchar(255)
);