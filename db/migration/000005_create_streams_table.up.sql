CREATE TABLE streams(
    id serial PRIMARY KEY,
    userID int NOT NULL,
    webhookurl VARCHAR(256) NOT NULL,
    FOREIGN KEY (userID) REFERENCES users (id)

);