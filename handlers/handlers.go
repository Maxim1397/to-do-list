package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv" //read .env file
	_ "github.com/lib/pq"      //postgres golang driver
	"log"
	"net/http"
	"os"
	"strconv"
	"to-do-list/models" //models package with Item model
)

//response format
type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

//create connection to postgres DB
func createConnection() *pgxpool.Pool {
	//load environment variables file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	pool, err := pgxpool.Connect(context.Background(), os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v\n", err)
	}

	return pool
}

//GetAllItems send response with all the items
func GetAllItems(w http.ResponseWriter, r *http.Request) {
	items, err := getAllItems()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server error"))
		log.Println(err.Error())
		return
	}

	//send all items response
	json.NewEncoder(w).Encode(items)
}

//get all items from db
func getAllItems() ([]models.Item, error) {
	//create connection with database
	pool := createConnection()
	defer pool.Close()

	var items []models.Item
	//create select sql query
	sqlStatement := `SELECT * FROM items`
	//execute sql statement
	rows, err := pool.Query(context.Background(), sqlStatement)
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

//GetItem send response with one item by its id
func GetItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	//get item id from the request params, and convert it to string
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request"))
		log.Println(err.Error())
		return
	}

	item, err := getItem(int64(id))
	switch err {
	case pgx.ErrNoRows:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Item not found"))
		log.Println(err.Error())
		return
	case nil:
		json.NewEncoder(w).Encode(item)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server error"))
		log.Println(err.Error())
		return
	}
}

//get one item from the DB by its id
func getItem(id int64) (models.Item, error) {
	//create connection with database
	pool := createConnection()
	defer pool.Close()
	var item models.Item
	//create query for select
	sqlStatement := `SELECT * FROM items WHERE id=$1`
	//execute sql statement
	row := pool.QueryRow(context.Background(), sqlStatement, id)
	err := row.Scan(&item.ID, &item.Description, &item.Status)
	return item, err
}

//CreateItem send response after new item creation
func CreateItem(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	//decode json request to item
	err := json.NewDecoder(r.Body).Decode(&item)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request"))
		log.Println(err.Error())
		return
	}

	insertID, err := insertItem(item)
	res := response{
		ID:      insertID,
		Message: "Item created successfully",
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server error"))
		log.Println(err.Error())
		return
	}

	//send response
	json.NewEncoder(w).Encode(res)
}

//insert new item to db
func insertItem(item models.Item) (int64, error) {
	//create connection with database
	pool := createConnection()
	defer pool.Close()

	//create query for insert
	sqlStatement := `INSERT INTO items (description) VALUES ($1) RETURNING id`
	var id int64
	//execute the statement
	err := pool.QueryRow(context.Background(), sqlStatement, item.Description).Scan(&id)
	return id, err
}

//UpdateItemStatus send response after item update
func UpdateItemStatus(w http.ResponseWriter, r *http.Request) {
	//get item id from the request params, and convert it to string
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request"))
		log.Println(err.Error())
		return
	}

	updatedRows, err := updateItemStatus(int64(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server error"))
		log.Println(err.Error())
		return
	}
	//format  message string
	msg := "Item updated successfully."
	if updatedRows == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Item not found"))
		return
	}

	//format response message
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	//send response
	json.NewEncoder(w).Encode(res)
}

//update item's status in db by its id
func updateItemStatus(id int64) (int64, error) {
	//create connection with database
	pool := createConnection()
	defer pool.Close()

	//create query for update
	sqlStatement := `UPDATE items SET status=TRUE WHERE id=$1`

	//execute the statement
	res, err := pool.Exec(context.Background(), sqlStatement, id)
	if err != nil {
		return -1, err
	}

	//check how many rows affected
	rowsAffected := res.RowsAffected()

	return rowsAffected, nil
}

//UpdateAllItemsStatus send request after all items update
func UpdateAllItemsStatus(w http.ResponseWriter, r *http.Request) {
	updatedRows, err := updateAllItemsStatus()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server error"))
		log.Println(err.Error())
		return
	}
	//format message string
	msg := fmt.Sprintf("Items updated successfully. Total rows/record affected %v", updatedRows)
	//send response
	json.NewEncoder(w).Encode(msg)
}

//update all item's statuses to true
func updateAllItemsStatus() (int64, error) {
	//create connection with postgres db
	pool := createConnection()
	defer pool.Close()

	//create query for update
	sqlStatement := `UPDATE items SET status=TRUE WHERE status=FALSE`

	//execute sql statement
	res, err := pool.Exec(context.Background(), sqlStatement)
	if err != nil {
		return -1, err
	}

	//check how many rows affected
	rowsAffected := res.RowsAffected()
	return rowsAffected, nil
}

//DeleteItem send response after item remove
func DeleteItem(w http.ResponseWriter, r *http.Request) {
	//get item id from the request params and convert int to string
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request"))
		log.Println(err.Error())
		return
	}

	deletedRows, err := deleteItem(int64(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server error"))
		log.Println(err.Error())
		return
	}

	//format message string
	msg := "Item deleted successfully."
	if deletedRows == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Item not found"))
		return
	}

	//format response message
	res := response{
		ID:      int64(id),
		Message: msg,
	}
	//send response
	json.NewEncoder(w).Encode(res)
}

//delete item from db by its id
func deleteItem(id int64) (int64, error) {
	//create connection with database
	pool := createConnection()
	defer pool.Close()

	//create query for delete
	sqlStatement := `DELETE FROM items WHERE id=$1`

	//execute the statement
	res, err := pool.Exec(context.Background(), sqlStatement, id)
	if err != nil {
		return -1, err
	}

	//check how many rows affected
	rowsAffected := res.RowsAffected()
	return rowsAffected, nil
}

//DeleteAllItems send response after all item removed
func DeleteAllItems(w http.ResponseWriter, r *http.Request) {
	err := deleteAllItems()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server error"))
		log.Println(err.Error())
		return
	}

	//format response message
	res := response{
		ID:      1,
		Message: "All Items deleted successfully.",
	}
	//send response
	json.NewEncoder(w).Encode(res)
}

//delete all items from db
func deleteAllItems() error {
	//create connection with database
	pool := createConnection()
	defer pool.Close()

	//create query for truncate
	sqlStatement := `TRUNCATE items`
	//execute statement
	_, err := pool.Exec(context.Background(), sqlStatement)
	return err
}
