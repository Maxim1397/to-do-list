# to-do-list
API for ToDo list

# How to build Go application:
https://golang.org/doc/tutorial/compile-install

# PostgreSQL:
https://www.postgresql.org/download/
See .env file for credentials.
Create table:
create table items
(
    id          serial                   not null
        constraint tasks_pkey
            primary key,
    description text    default ''::text not null,
    status      boolean default false    not null
);

alter table items
    owner to postgres
