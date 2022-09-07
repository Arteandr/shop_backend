CREATE TABLE users
(
    id           serial primary key not null unique,
    email        varchar(255)       not null unique,
    login        varchar(30)        not null unique,
    password     varchar(255)       not null,
    first_name   varchar(20),
    last_name    varchar(20),
    company_name varchar(30),
    admin        boolean default false,
    completed    boolean default false
);

CREATE TABLE phone_numbers
(
    user_id integer references users (id) on delete cascade not null unique,
    code    varchar(5),
    number  varchar(15)
);

CREATE TABLE address
(
    id      serial primary key not null unique,
    country varchar(255)       not null,
    city    varchar(255)       not null,
    street  varchar(255)       not null,
    zip     varchar(255)       not null
);

CREATE TABLE users_invoice
(
    user_id    integer references users (id) on delete cascade not null,
    address_id integer references address (id) on delete cascade
);

CREATE TABLE users_shipping
(
    user_id    integer references users (id) on delete cascade not null,
    address_id integer references address (id) on delete cascade
);

CREATE TABLE sessions
(
    user_id       integer references users (id) on delete cascade not null,
    refresh_token varchar(255)                                    not null,
    expires_at    timestamp                                       not null
);