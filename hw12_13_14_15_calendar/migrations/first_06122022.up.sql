CREATE DATABASE calendar;

create table calendar.users
(
    id int,
    name varchar(255) default '' not null
);

create unique index users_id_uindex
    on users (id);

alter table users modify id int auto_increment;

create table calendar.events
(
    id binary(36) not null,
    title varchar(255) default '' not null,
    date datetime not null,
    user_id int not null,
    constraint events_pk
        primary key (id)
);

create unique index events_id_uindex
    on events (id);

