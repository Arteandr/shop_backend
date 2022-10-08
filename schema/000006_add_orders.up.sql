CREATE TABLE statuses
(
    id   serial primary key not null,
    name varchar(30)        not null
);
INSERT INTO statuses (name)
VALUES ('Waiting for payment');
INSERT INTO statuses (name)
VALUES ('Processing');
INSERT INTO statuses (name)
VALUES ('Queued');
INSERT INTO statuses (name)
VALUES ('Completed');
INSERT INTO statuses (name)
VALUES ('Canceled');

CREATE TABLE payment_methods
(
    id          serial primary key not null,
    name        varchar(30) unique not null,
    description varchar(30)        not null,
    active      bool               not null
);

CREATE TABLE payment_methods_images
(
    payment_method_id int references payment_methods (id) on delete cascade unique not null,
    image_id          int references images (id) on delete cascade                 not null
);

CREATE TABLE orders
(
    id          serial primary key                          not null,
    user_id     int references users (id) on delete cascade not null,
    delivery_id int references delivery (id)                not null,
    status_id   int references statuses (id) default 1      not null,
    payment_id  int references payment_methods (id)         not null,
    comment     varchar(255)                                not null,
    created_at  timestamp                    default now()
);

CREATE TABLE order_items
(
    order_id int references orders (id) not null,
    item_id  int references items (id)  not null,
    color_id int references colors (id) not null,
    quantity int                        not null,
    check ( quantity > 0 )
);

