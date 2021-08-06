# to-do-list
API for ToDo list

# How to build Go application:
https://golang.org/doc/tutorial/compile-install

# PostgreSQL:
https://www.postgresql.org/download/ <br />
See .env file for credentials. <br />
- Create table: <br />
create table items<br />
( <br />
    id serial not null constraint tasks_pkey primary key,
    description text    default ''::text not null,
    status      boolean default false    not null
);
alter table items owner to postgres
