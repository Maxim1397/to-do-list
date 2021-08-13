package main

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"to-do-list/models"
	"to-do-list/repository"
	"zombiezen.com/go/postgrestest"
)

const createTable = `create table items
					(
					id serial not null constraint tasks_pkey primary key,
					description text default ''::text not null,
					status boolean default false not null
					);`

func NewMock() sqlmock.Sqlmock {
	_, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return mock
}

func TestCreateItem(t *testing.T) {
	mock := NewMock()
	ctx := context.Background()
	srv, err := postgrestest.Start(ctx)
	assert.NoError(t, err)
	pool, err := pgxpool.Connect(context.Background(), srv.DefaultDatabase())
	assert.NoError(t, err)
	_, err = pool.Exec(ctx, createTable)
	assert.NoError(t, err)
	repo := repository.Connection{
		Pool: pool,
	}
	defer repo.Pool.Close()

	item := models.Item{
		ID:          1,
		Description: "description",
		Status:      false,
	}
	query := "INSERT INTO items \\(id, description, status\\) VALUES \\(\\?, \\?, \\?\\)"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(item.ID, item.Description, item.Status).WillReturnResult(sqlmock.NewResult(0, 1))
	id, err := repo.InsertItem(item)
	assert.Equal(t, id, item.ID)
	assert.NoError(t, err)
	t.Cleanup(srv.Cleanup)
}

func TestGetItem(t *testing.T) {
	mock := NewMock()
	ctx := context.Background()
	srv, err := postgrestest.Start(ctx)
	assert.NoError(t, err)
	pool, err := pgxpool.Connect(context.Background(), srv.DefaultDatabase())
	assert.NoError(t, err)
	_, err = pool.Exec(ctx, createTable)
	assert.NoError(t, err)
	repo := repository.Connection{
		Pool: pool,
	}
	defer repo.Pool.Close()

	item := models.Item{
		ID:          1,
		Description: "description",
		Status:      false,
	}

	//not to get "no rows in result set"
	_, err = repo.InsertItem(item)
	assert.NoError(t, err)

	query := "SELECT id, description, status FROM items WHERE id = \\?"
	rows := sqlmock.NewRows([]string{"id", "description", "status"})
	mock.ExpectQuery(query).WithArgs(item.ID).WillReturnRows(rows)

	getItem, err := repo.GetItem(item.ID)
	assert.NoError(t, err)
	assert.Equal(t, getItem, item)
	t.Cleanup(srv.Cleanup)
}

func TestUpdateItemStatus(t *testing.T) {
	mock := NewMock()
	ctx := context.Background()
	srv, err := postgrestest.Start(ctx)
	assert.NoError(t, err)
	pool, err := pgxpool.Connect(context.Background(), srv.DefaultDatabase())
	assert.NoError(t, err)
	_, err = pool.Exec(ctx, createTable)
	assert.NoError(t, err)
	repo := repository.Connection{
		Pool: pool,
	}
	defer repo.Pool.Close()

	item := models.Item{
		ID:          1,
		Description: "description",
		Status:      false,
	}

	//not to get "no rows in result set"
	_, err = repo.InsertItem(item)
	assert.NoError(t, err)

	query := "UPDATE items SET status=TRUE WHERE id=\\?"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(item.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	id, err := repo.UpdateItemStatus(item.ID)
	assert.NoError(t, err)
	assert.Equal(t, id, item.ID)

	t.Cleanup(srv.Cleanup)
}

func TestDeleteItem(t *testing.T) {
	mock := NewMock()
	ctx := context.Background()
	srv, err := postgrestest.Start(ctx)
	assert.NoError(t, err)
	pool, err := pgxpool.Connect(context.Background(), srv.DefaultDatabase())
	assert.NoError(t, err)
	_, err = pool.Exec(ctx, createTable)
	assert.NoError(t, err)
	repo := repository.Connection{
		Pool: pool,
	}
	defer repo.Pool.Close()

	item := models.Item{
		ID:          1,
		Description: "description",
		Status:      false,
	}

	//not to get "no rows in result set"
	_, err = repo.InsertItem(item)
	assert.NoError(t, err)

	query := "DELETE FROM items WHERE id = \\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(item.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	id, err := repo.DeleteItem(item.ID)
	assert.NoError(t, err)
	assert.Equal(t, id, item.ID)

	t.Cleanup(srv.Cleanup)
}
