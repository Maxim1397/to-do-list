package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"log"
	"os"
	"to-do-list/models"
)

type Repository interface {
	GetAllItems() ([]models.Item, error)
	GetItem(id int64) (models.Item, error)
	InsertItem(item models.Item) (int64, error)
	UpdateItemStatus(id int64) (int64, error)
	UpdateAllItemsStatus() (int64, error)
	DeleteItem(id int64) (int64, error)
	DeleteAllItems() error
}

type Connection struct {
	Pool *pgxpool.Pool
}

//create connection to postgres DB
func NewRepository() Repository {
	//load environment variables file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	pool, err := pgxpool.Connect(context.Background(), os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v\n", err)
	}

	return &Connection{
		pool,
	}
}

//get all items from db
func (c *Connection) GetAllItems() ([]models.Item, error) {
	var items []models.Item
	//create select sql query
	sqlStatement := `SELECT * FROM items`
	//execute sql statement
	rows, err := c.Pool.Query(context.Background(), sqlStatement)
	if err != nil {
		return nil, err
	}

	//close statement
	defer rows.Close()
	//iterate over the rows
	for rows.Next() {
		var item models.Item
		//read value from row to item
		err = rows.Scan(&item.ID, &item.Description, &item.Status)
		if err != nil {
			return nil, err
		}

		//append item to the items
		items = append(items, item)
	}

	return items, err
}

//get one item from the DB by its id
func (c *Connection) GetItem(id int64) (models.Item, error) {
	var item models.Item
	//create query for select
	sqlStatement := `SELECT * FROM items WHERE id=$1`
	//execute sql statement
	row := c.Pool.QueryRow(context.Background(), sqlStatement, id)
	err := row.Scan(&item.ID, &item.Description, &item.Status)
	return item, err
}

//insert new item to db
func (c *Connection) InsertItem(item models.Item) (int64, error) {
	//create query for insert
	sqlStatement := `INSERT INTO items (description) VALUES ($1) RETURNING id`
	var id int64
	//execute the statement
	err := c.Pool.QueryRow(context.Background(), sqlStatement, item.Description).Scan(&id)
	return id, err
}

//update item's status in db by its id
func (c *Connection) UpdateItemStatus(id int64) (int64, error) {
	//create query for update
	sqlStatement := `UPDATE items SET status=TRUE WHERE id=$1`

	//execute the statement
	res, err := c.Pool.Exec(context.Background(), sqlStatement, id)
	if err != nil {
		return -1, err
	}

	//check how many rows affected
	rowsAffected := res.RowsAffected()

	return rowsAffected, nil
}

//update all item's statuses to true
func (c *Connection) UpdateAllItemsStatus() (int64, error) {
	//create query for update
	sqlStatement := `UPDATE items SET status=TRUE WHERE status=FALSE`

	//execute sql statement
	res, err := c.Pool.Exec(context.Background(), sqlStatement)
	if err != nil {
		return -1, err
	}

	//check how many rows affected
	rowsAffected := res.RowsAffected()
	return rowsAffected, nil
}

//delete item from db by its id
func (c *Connection) DeleteItem(id int64) (int64, error) {
	//create query for delete
	sqlStatement := `DELETE FROM items WHERE id=$1`

	//execute the statement
	res, err := c.Pool.Exec(context.Background(), sqlStatement, id)
	if err != nil {
		return -1, err
	}

	//check how many rows affected
	rowsAffected := res.RowsAffected()
	return rowsAffected, nil
}

//delete all items from db
func (c *Connection) DeleteAllItems() error {
	//create query for truncate
	sqlStatement := `TRUNCATE items`
	//execute statement
	_, err := c.Pool.Exec(context.Background(), sqlStatement)
	return err
}
