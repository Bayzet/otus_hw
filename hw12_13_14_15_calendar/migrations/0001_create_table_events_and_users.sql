-- +goose Up
create table users (
    id int primary key,
    name text default '' not null
);

create table events (
    id uuid not null primary key,
    title text not null,
    date      timestamp not null,
    user_id int not null
);

-- +goose Down
drop table events;
drop table users;