CREATE TABLE items (
    id serial primary key,
    description varchar(255) not null,
    cost int not null,
    timestamp timestamp not null default now()
);
