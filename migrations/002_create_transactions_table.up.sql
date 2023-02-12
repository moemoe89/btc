CREATE TABLE transactions (
    time TIMESTAMPTZ NOT NULL,
    user_id INTEGER NOT NULL,
    amount FLOAT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

SELECT create_hypertable('transactions', 'time');
