# to-do-list
API for ToDo list <br />

# How to build Go application:
https://golang.org/doc/tutorial/compile-install <br />

# PostgreSQL:
https://www.postgresql.org/download/ <br />
See .env file for credentials. <br />
- Create table: <br />
create table items<br />
( <br />
    id serial not null constraint tasks_pkey primary key, <br />
    description text    default ''::text not null, <br />
    status      boolean default false    not null <br />
); <br />
alter table items owner to postgres
