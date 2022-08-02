CREATE TABLE users
(
    id       serial primary key not null unique,
    email    varchar(255)       not null,
    login    varchar(30)        not null,
    password varchar(255)       not null
);

CREATE TABLE address
(
    id      serial primary key not null unique,
    country varchar(255)       not null,
    street  varchar(255)       not null,
    city    varchar(255)       not null,
    zip     integer            not null
);

CREATE TABLE users_invoice
(
    user_id    integer references users (id)   not null,
    address_id integer references address (id) not null
);

CREATE TABLE users_shipping
(
    user_id    integer references users (id)   not null,
    address_id integer references address (id) not null
);