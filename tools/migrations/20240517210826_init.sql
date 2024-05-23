-- +goose Up
-- +goose StatementBegin
CREATE TABLE delivery_service
(
    id            SERIAL PRIMARY KEY NOT NULL,
    delivery_name TEXT NOT NULL
);

CREATE TABLE delivery
(
    id                  SERIAL PRIMARY KEY NOT NULL,
    zip                 TEXT NOT NULL,
    city                TEXT NOT NULL,
    address             TEXT NOT NULL,
    region              TEXT NOT NULL,
    delivery_service_id BIGINT NOT NULL,
    CONSTRAINT delivery_service_id_fk FOREIGN KEY (delivery_service_id)
        REFERENCES delivery_service(id) ON DELETE CASCADE
);

CREATE TABLE bank
(
    id        SERIAL PRIMARY KEY NOT NULL,
    bank_name TEXT NOT NULL
);

CREATE TABLE currency
(
    id            SERIAL PRIMARY KEY NOT NULL,
    currency_name TEXT NOT NULL
);

CREATE TABLE payment
(
    id            SERIAL PRIMARY KEY,
    transaction   TEXT NOT NULL,
    request_id    TEXT,
    provider      TEXT NOT NULL,
    amount        BIGINT NOT NULL,
    payment_dt    TIMESTAMP NOT NULL,
    delivery_cost BIGINT NOT NULL,
    goods_total   BIGINT NOT NULL,
    custom_fee    BIGINT NOT NULL,
    bank_id       BIGINT NOT NULL,
    currency_id   BIGINT NOT NULL,
    CONSTRAINT bank_id_fk FOREIGN KEY (bank_id) REFERENCES bank(id)
        ON DELETE CASCADE,
    CONSTRAINT currency_id_fk FOREIGN KEY (currency_id) REFERENCES currency(id)
        ON DELETE CASCADE
);

CREATE TABLE items
(
    id    SERIAL PRIMARY KEY NOT NULL,
    price BIGINT NOT NULL,
    rid   TEXT NOT NULL,
    name  TEXT NOT NULL,
    size  TEXT NOT NULL,
    nm_id BIGINT NOT NULL,
    brand TEXT NOT NULL
);

CREATE TABLE status
(
    id          SERIAL PRIMARY KEY NOT NULL,
    name_status BIGINT NOT NULL
);

CREATE TABLE customer
(
    id          SERIAL PRIMARY KEY NOT NULL,
    name        TEXT NOT NULL,
    phone       TEXT NOT NULL,
    email       TEXT NOT NULL,
    external_customer_id TEXT NOT NULL
);

CREATE TABLE orders
(
    id                 SERIAL PRIMARY KEY,
    order_uid          TEXT NOT NULL,
    track_number       TEXT NOT NULL,
    chrt_id            BIGINT NOT NULL,
    entry              TEXT NOT NULL,
    locale             TEXT NOT NULL,
    internal_signature TEXT,
    shardkey           TEXT NOT NULL,
    sm_id              BIGINT NOT NULL,
    date_created       TIMESTAMP NOT NULL,
    off_shard          TEXT NOT NULL,
    payment_id         BIGINT NOT NULL,
    status_id          BIGINT NOT NULL,
    delivery_id        BIGINT NOT NULL,
    customer_id        BIGINT NOT NULL,

    CONSTRAINT payment_id_fk FOREIGN KEY (payment_id) REFERENCES payment(id) ON
        DELETE CASCADE,
    CONSTRAINT status_id_fk FOREIGN KEY (status_id) REFERENCES status(id) ON
        DELETE CASCADE,
    CONSTRAINT delivery_id_fk FOREIGN KEY (delivery_id) REFERENCES delivery(id)
        ON DELETE CASCADE,
    CONSTRAINT customer_id_fk FOREIGN KEY (customer_id) REFERENCES customer(id)
        ON DELETE CASCADE
);

CREATE TABLE order_items
(
    id          SERIAL PRIMARY KEY,
    sale        BIGINT NOT NULL,
    total_price BIGINT NOT NULL,
    item_id     BIGINT NOT NULL,
    order_id    BIGINT NOT NULL,
    CONSTRAINT item_id_fk FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE
        CASCADE,
    CONSTRAINT order_id_fk FOREIGN KEY (order_id) REFERENCES orders(id) ON
        DELETE CASCADE
);

ALTER TABLE bank ADD CONSTRAINT unique_bank_name UNIQUE (bank_name);

ALTER TABLE currency ADD CONSTRAINT unique_currency_name UNIQUE (currency_name);

ALTER TABLE customer ADD CONSTRAINT unique_external_customer_id UNIQUE (external_customer_id);

ALTER TABLE delivery ADD CONSTRAINT unique_delivery_all UNIQUE (zip, city, address, region, delivery_service_id);

ALTER TABLE delivery_service ADD CONSTRAINT unique_delivery_name UNIQUE (delivery_name);

ALTER TABLE items ADD CONSTRAINT unique_rid UNIQUE (rid);

ALTER TABLE orders ADD CONSTRAINT unique_order_uid UNIQUE (order_uid);

ALTER TABLE payment ADD CONSTRAINT unique_transaction UNIQUE (transaction);

ALTER TABLE status ADD CONSTRAINT unique_status_name UNIQUE (name_status);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE order_items, orders, customer, status, items, payment, currency, bank, delivery, delivery_service;
-- +goose StatementEnd
