CREATE TABLE images
(
    id         serial primary key not null,
    filename   varchar(255)       not null,
    created_at timestamp default now()
);