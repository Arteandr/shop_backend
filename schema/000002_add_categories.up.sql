CREATE TABLE item_categories
(
    item_id int4 references items(id),
    category_id int4 references categories(id)
);