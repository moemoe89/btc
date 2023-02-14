CREATE TABLE transactions (
    datetime TIMESTAMPTZ NOT NULL,
    user_id INTEGER NOT NULL,
    amount DECIMAL NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

SELECT create_hypertable('transactions', 'datetime');
