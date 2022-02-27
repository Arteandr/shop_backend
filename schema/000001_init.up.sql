CREATE TABLE categories
(
    id   serial primary key not null unique,
    name varchar(255)       not null
);

CREATE TABLE items
(
    id          serial primary key                                   not null unique,
    name        varchar(255)                                         not null,
    price       decimal(10, 2)                                       not null check ( price > 0 ),
    description varchar(255)                                         not null,
    category_id integer references categories (id) on delete cascade not null,
    tags        varchar(255)[]                                       not null,
    created_at  timestamp                                            not null
);