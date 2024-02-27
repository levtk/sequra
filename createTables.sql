CREATE TABLE IF NOT EXISTS DISBURSEMENT (
    record_uuid UUID PRIMARY KEY ,
    disbursement_group_id UUID,
    transaction_id UUID, -- not implemented. when payout is confirmed by payment method/system the id for the payment transaction should be saved
    merchReference varchar(255) NOT NULL,
    order_id char(12) NOT NULL UNIQUE ,
    order_fee INT NOT NULL,
    order_fee_running_total INT,
    payout_date datetime,
    payout_running_total INT,
    payout_total INT,
    is_paid_out INT,
    createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP);

CREATE TABLE IF NOT EXISTS ORDERS (
    id char(12) PRIMARY KEY ,
    merchant_reference varchar(255) NOT NULL,
    merchant_id uuid NOT NULL,
    amount INT NOT NULL,
    created_at datetime);

CREATE TABLE IF NOT EXISTS MERCHANTS (
    id char(128) PRIMARY KEY,
    reference varchar(255),
    email varchar(255),
    live_on date,
    disbursement_frequency varchar(6),
    minimum_monthly_fee varchar(5));

CREATE TABLE IF NOT EXISTS MONTHLY (
    id UUID primary key,
    merchant_id UUID,
    merchant_reference varchar(255),
    monthly_fee_date date,
    did_pay_fee INT,
    monthly_fee INT,
    total_order_amt INT,
    order_fee_total INT,
    amt_monthly_fee_paid INT GENERATED ALWAYS AS (monthly_fee-order_fee_total) VIRTUAL,
    createdAt datetime,
    updatedAt datetime
);