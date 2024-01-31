CREATE TABLE IF NOT EXISTS "user"
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(50)  NOT NULL,
    email      VARCHAR(255) NOT NULL UNIQUE,
    password   VARCHAR(72)  NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS currency
(
    id     SERIAL PRIMARY KEY,
    name   VARCHAR(50) NOT NULL,
    symbol VARCHAR(3)  NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS account
(
    id          SERIAL PRIMARY KEY,
    user_id     INT            NOT NULL,
    currency_id INT            NOT NULL,
    amount      DECIMAL(10, 2) NOT NULL DEFAULT 0 CHECK (amount >= 0),
    FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE,
    FOREIGN KEY (currency_id) REFERENCES currency (id)
);

CREATE TABLE IF NOT EXISTS transaction
(
    id              SERIAL PRIMARY KEY,
    from_account_id INT,
    to_account_id   INT,
    currency_id     INT            NOT NULL,
    amount          DECIMAL(10, 2) NOT NULL,
    created_at      TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (from_account_id) REFERENCES account (id),
    FOREIGN KEY (to_account_id) REFERENCES account (id),
    FOREIGN KEY (currency_id) REFERENCES currency (id)
);

INSERT INTO currency (name, symbol)
VALUES ('Russian Ruble', 'RUB'),
       ('Euro', 'EUR'),
       ('US Dollar', 'USD'),
       ('British Pound', 'GBP'),
       ('Japanese Yen', 'JPY')
ON CONFLICT (symbol) DO NOTHING;