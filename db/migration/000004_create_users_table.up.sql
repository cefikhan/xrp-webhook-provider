CREATE TABLE users(
    id serial PRIMARY KEY,
    username VARCHAR(256) NOT NULL,
    email VARCHAR(256) NOT NULL,
    userpassword VARCHAR(256) NOT NULL
);