CREATE TABLE orders (
    order_uid           TEXT PRIMARY KEY,
    track_number        TEXT NOT NULL,
    entry               TEXT,
    locale              TEXT,
    internal_signature  TEXT,
    customer_id         TEXT,
    delivery_service    TEXT,
    shardkey            TEXT,
    sm_id               INTEGER,
    date_created        TIMESTAMP,
    oof_shard           TEXT,

    delivery_name       TEXT,
    delivery_phone      TEXT,
    delivery_zip        TEXT,
    delivery_city       TEXT,
    delivery_address    TEXT,
    delivery_region     TEXT,
    delivery_email      TEXT,

    payment_transaction     TEXT,
    payment_request_id      TEXT,
    payment_currency        TEXT,
    payment_provider        TEXT,
    payment_amount          INTEGER,
    payment_dt              BIGINT,
    payment_bank            TEXT,
    payment_delivery_cost   INTEGER,
    payment_goods_total     INTEGER,
    payment_custom_fee      INTEGER
);

CREATE TABLE items (
    id            BIGSERIAL PRIMARY KEY,
    order_uid     TEXT,
    chrt_id       INTEGER,
    track_number  TEXT,
    price         INTEGER,
    rid           TEXT,
    name          TEXT,
    sale          INTEGER,
    size          TEXT,
    total_price   INTEGER,
    nm_id         INTEGER,
    brand         TEXT,
    status        INTEGER
);
