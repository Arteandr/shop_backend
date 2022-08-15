CREATE TABLE images
(
    id         serial primary key not null,
    filename   varchar(255)       not null,
    created_at timestamp default now()
);

CREATE TABLE items_images
(
    id       serial primary key                           not null,
    item_id  int references items (id) on delete cascade  not null,
    image_id int references images (id) on delete cascade not null
);

CREATE TABLE categories_images
(
    id          serial primary key                               not null,
    category_id int references categories (id) on delete cascade not null,
    image_id    int references images (id) on delete cascade     not null
)