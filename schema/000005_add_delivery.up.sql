CREATE TABLE delivery_company
(
    id   serial primary key not null unique,
    name varchar(255) unique
);

CREATE TABLE delivery
(
    id         serial primary key                                     not null unique,
    name       varchar(255)                                           not null,
    company_id int references delivery_company (id) on delete cascade not null,
    price      decimal                                                not null
);

