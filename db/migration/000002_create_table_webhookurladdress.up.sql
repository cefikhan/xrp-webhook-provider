CREATE TABLE webhookurladdress(
    id serial PRIMARY KEY,
    url  VARCHAR(256) NOT NULL,
    addresses VARCHAR(256)[] 

);