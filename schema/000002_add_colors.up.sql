CREATE TABLE colors
(
    id    serial primary key not null unique,
    name  varchar(255)       not null,
    hex   varchar(255)       not null,
    price decimal(10, 2)     not null
);


CREATE TABLE item_colors
(
    item_id  int4 references items (id),
    color_id int4 references colors (id)
);