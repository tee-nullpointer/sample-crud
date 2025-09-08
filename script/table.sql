create table sample.products
(
    id         SERIAL
        constraint products_pk
            primary key,
    name       varchar,
    created_at timestamp,
    updated_at timestamp
);