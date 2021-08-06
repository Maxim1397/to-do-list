# to-do-list
API for ToDo list

# How to build Go application:
https://golang.org/doc/tutorial/compile-install __

# PostgreSQL:
https://www.postgresql.org/download/ __
See .env file for credentials. __
- Create table: __
create table items __
( __
    id serial not null constraint tasks_pkey primary key, __
    description text    default ''::text not null, __
    status      boolean default false    not null __
); __
__
alter table items owner to postgres
