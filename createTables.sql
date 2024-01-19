CREATE TABLE IF NOT EXISTS DISBURSEMENT (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
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
    createdAt timestamp default (strftime('%s', 'now')));

CREATE TABLE IF NOT EXISTS ORDERS (
    id TEXT NOT NULL PRIMARY KEY,
    merchant_reference TEXT NOT NULL,
    merchant_id TEXT NOT NULL,
    amount INT NOT NULL,
    created_at TEXT);

CREATE TABLE IF NOT EXISTS MERCHANTS (
    id TEXT PRIMARY KEY,
    reference TEXT,
    email TEXT,
    live_on TEXT,
    disbursement_frequency TEXT,
    minimum_monthly_fee TEXT);