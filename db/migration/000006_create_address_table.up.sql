CREATE TABLE addresses(
    id serial PRIMARY KEY,
    streamid INT NOT NULL,
    address VARCHAR(256) NOT NULL,
    FOREIGN KEY (streamid) REFERENCES streams (id)

);