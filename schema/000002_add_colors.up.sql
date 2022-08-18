CREATE TABLE colors
(
    id    serial primary key not null unique,
    name  varchar(30)        not null,
    hex   varchar(10)        not null,
    price decimal(10, 2)     not null
);


CREATE TABLE items_colors
(
    id       serial primary key                           not null,
    item_id  int references items (id) on delete cascade  not null,
    color_id int references colors (id) on delete cascade not null
);

ALTER TABLE items_colors
    ADD CONSTRAINT items_colors_item_color_id_c UNIQUE (item_id, color_id);