-- +migrate Up
CREATE TABLE merch (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    price_rubles INT NOT NULL CHECK (price_rubles >= 0)
);

CREATE INDEX idx_merch_name ON merch (name);

INSERT INTO merch (name, price_rubles) VALUES ('t-shirt', 80),
                                              ('cup', 20),
                                              ('book', 50),
                                              ('pen', 10),
                                              ('powerbank',200),
                                              ('hoody', 300),
                                              ('umbrella', 200),
                                              ('socks', 10),
                                              ('wallet', 50),
                                              ('pink-hoody', 500);

CREATE TABLE employees(
    id SERIAL PRIMARY KEY,
    username VARCHAR(127) NOT NULL,
    hashed_password VARCHAR(127) NOT NULL,
);

CREATE INDEX id_employees_username ON employees (username);

CREATE TABLE auth