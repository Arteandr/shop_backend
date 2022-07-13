CREATE TABLE colors
(
    id    serial primary key not null unique,
    name  varchar(255)       not null,
    hex   varchar(255)       not null,
    price decimal(10, 2)     not null
);


CREATE TABLE item_colors
(
    id       serial primary key                           not null,
    item_id  int references items (id) on delete cascade  not null,
    color_id int references colors (id) on delete cascade not null
);