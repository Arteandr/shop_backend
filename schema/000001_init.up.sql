CREATE TABLE categories
(
    id   serial primary key not null unique,
    name varchar(255)       not null unique
);

CREATE TABLE items
(
    id          serial primary key                 not null unique,
    name        varchar(255)                       not null,
    description varchar(255)                       not null,
    category_id integer references categories (id) not null,
    price       decimal                            not null,
    sku         varchar(255)                       not null unique,
    created_at  timestamp default now()
);

CREATE TABLE tags
(
    id      serial primary key                          not null,
    item_id int references items (id) on delete cascade not null,
    name    varchar(255)                                not null
);