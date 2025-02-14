CREATE TABLE merch
(
    id    SERIAL PRIMARY KEY,
    name  VARCHAR(255) NOT NULL UNIQUE,
    price INT          NOT NULL CHECK (price >= 0)
);


CREATE TABLE employees
(
    id              SERIAL PRIMARY KEY,
    username        VARCHAR(127) NOT NULL UNIQUE,
    hashed_password TEXT         NOT NULL,
    coins           INT          NOT NULL DEFAULT 1000 CHECK (coins >= 0)
);

CREATE TABLE inventory
(
    id          SERIAL PRIMARY KEY,
    employee_id INT NOT NULL,
    merch_id    INT NOT NULL,
    quantity    INT NOT NULL CHECK (quantity > 0),
    UNIQUE (employee_id, merch_id),
    FOREIGN KEY (employee_id) REFERENCES employees (id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (merch_id) REFERENCES merch (id) ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE TABLE coin_transactions
(
    id         SERIAL PRIMARY KEY,
    sender     INT NULL,
    recipient  INT NULL,
    amount     INT NOT NULL CHECK (amount >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    FOREIGN KEY (sender) REFERENCES employees (id) ON DELETE SET NULL ON UPDATE CASCADE,
    FOREIGN KEY (recipient) REFERENCES employees (id) ON DELETE SET NULL ON UPDATE CASCADE
);

CREATE INDEX idx_coin_transactions_sender_time ON coin_transactions (sender, created_at DESC);
CREATE INDEX idx_coin_transactions_recipient_time ON coin_transactions (recipient, created_at DESC);

-- static row in db to make shop transactions correct
INSERT INTO employees (username, hashed_password)
VALUES ('shop', 'SHOP_HASH');

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
