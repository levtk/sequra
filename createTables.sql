CREATE TABLE IF NOT EXISTS DISBURSEMENT (
    id INTEGER PRIMARY KEY,
    record_uuid TEXT NOT NULL,
    disbursement_group_id TEXT,
    transaction_id TEXT, -- not implemented. when payout is confirmed by payment method/system the id for the payment transaction should be saved
    merchReference TEXT NOT NULL,
    order_id TEXT NOT NULL UNIQUE ,
    order_fee INT NOT NULL,
    order_fee_running_total INT,
    payout_date TEXT,
    payout_running_total INT,
    payout_total INT,
    is_paid_out INT,
    createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP);

CREATE TABLE IF NOT EXISTS ORDERS (
    id char(12) PRIMARY KEY ,
    merchant_reference TEXT NOT NULL,
    merchant_id TEXT NOT NULL,
    amount INT NOT NULL,
    created_at TEXT);

CREATE TABLE IF NOT EXISTS MERCHANTS (
    id varchar(128) PRIMARY KEY,
    reference varchar(50),
    email varchar(128),
    live_on date,
    disbursement_frequency varchar(6),
    minimum_monthly_fee varchar(5));

CREATE TABLE IF NOT EXISTS MONTHLY (
    id varchar(128) primary key,
    merchant_id TEXT,
    merchant_reference TEXT,
    monthly_fee_date date,
    did_pay_fee INT,
    monthly_fee INT,
    total_order_amt INT,
    order_fee_total INT,
    amt_monthly_fee_paid INT GENERATED ALWAYS AS (monthly_fee-order_fee_total) VIRTUAL,
    createdAt timestamp,
    updatedAt timestamp
);