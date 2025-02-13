CREATE TABLE merch
(
    id    SERIAL PRIMARY KEY,
    name  VARCHAR(255) NOT NULL UNIQUE,
    price INT          NOT NULL CHECK (price >= 0)
);

CREATE INDEX idx_merch_name ON merch USING HASH (name);

INSERT INTO merch (name, price)
VALUES ('t-shirt', 80),
       ('cup', 20),
       ('book', 50),
       ('pen', 10),
       ('powerbank', 200),
       ('hoody', 300),
       ('umbrella', 200),
       ('socks', 10),
       ('wallet', 50),
       ('pink-hoody', 500);

CREATE TABLE employees
(
    id              SERIAL PRIMARY KEY,
    username        VARCHAR(127) NOT NULL UNIQUE,
    hashed_password TEXT         NOT NULL,
    coins           INT          NOT NULL DEFAULT 1000 CHECK (coins >= 0)
);

INSERT INTO employees (username, hashed_password)
VALUES ('shop', 'SHOP_HASH');

CREATE INDEX idx_employees_username ON employees USING HASH (username);

CREATE TABLE inventory
(
    id          SERIAL PRIMARY KEY,
    employee_id INT NOT NULL,
    merch_id    INT NOT NULL,
    quantity    INT NOT NULL CHECK (quantity > 0),
    UNIQUE (employee_id, merch_id),
    FOREIGN KEY (employee_id) REFERENCES employees (id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (merch_id) REFERENCES merch (id) ON DELETE CASCADE ON UPDATE CASCADE
);


CREATE TABLE coin_transactions
(
    id         SERIAL PRIMARY KEY,
    sender       INT NOT NULL,
    recipient           INT NOT NULL,
    amount     INT NOT NULL CHECK (amount >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    FOREIGN KEY (sender) REFERENCES employees (id) ON DELETE CASCADE,
    FOREIGN KEY (recipient) REFERENCES employees (id) ON DELETE CASCADE
);

CREATE INDEX idx_coin_transactions_sender ON coin_transactions (sender);
CREATE INDEX idx_coin_transactions_recipient ON coin_transactions (recipient);
