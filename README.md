# to-do-list
API for ToDo list <br />

# IDE for Go:
https://www.jetbrains.com/go/promo/?gclid=Cj0KCQjwu7OIBhCsARIsALxCUaMHcehipAVScSinCF7HRn70vrvrQhFfbvxAbjKU7LAwCCTvEQs8YSIaAi8UEALw_wcB

# How to build Go application:
https://golang.org/doc/tutorial/compile-install <br />

# PostgreSQL:
https://www.postgresql.org/download/ <br />
See .env file for credentials. <br />
Create db: <br />
CREATE DATABASE todolist;
Create table: <br />
- create table items<br />
( <br />
    id serial not null constraint tasks_pkey primary key, <br />
    description text    default ''::text not null, <br />
    status      boolean default false    not null <br />
); <br />
- alter table items owner to postgres

# API :
After application build:
    http://localhost:8084
 - GET: <br />
    Get all items : /items <br />
    Get item by id : /items/{id} <br />
 - POST: <br />
    {
    "description": "Drink water"
    }
    Create new item : /items <br />
 - PUT: <br />
    Update item's status by id : /items/{id} <br />
    Update all item's statuses : /items <br />
 - DELETE: <br />
    Delete item by id : /items/{id} <br />
    Delete all items : /items
 
# Postman (API platform for building and using APIs)
https://www.postman.com/

# Test
    go test
