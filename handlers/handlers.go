package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	_ "github.com/lib/pq" //postgres golang driver
	"log"
	"net/http"
	"strconv"
	"to-do-list/models" //models package with Item model
	"to-do-list/repository"
)

//response format
type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

type Handler interface {
	GetAllItems(w http.ResponseWriter, r *http.Request)
	GetItem(w http.ResponseWriter, r *http.Request)
	CreateItem(w http.ResponseWriter, r *http.Request)
	UpdateItemStatus(w http.ResponseWriter, r *http.Request)
	UpdateAllItemsStatus(w http.ResponseWriter, r *http.Request)
	DeleteItem(w http.ResponseWriter, r *http.Request)
	DeleteAllItems(w http.ResponseWriter, r *http.Request)
}

type Repository struct {
	Repository repository.Repository
}

func NewHandler() Handler {
	return &Repository{
		repository.NewRepository(),
	}
}

//GetAllItems send response with all the items
func (repo *Repository) GetAllItems(w http.ResponseWriter, r *http.Request) {
	items, err := repo.Repository.GetAllItems()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server error"))
		log.Println(err.Error())
		return
	}

	//send all items response
	json.NewEncoder(w).Encode(items)
}

//GetItem send response with one item by its id
func (repo *Repository) GetItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	//get item id from the request params, and convert it to string
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request"))
		log.Println(err.Error())
		return
	}

	item, err := repo.Repository.GetItem(int64(id))
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

//CreateItem send response after new item creation
func (repo *Repository) CreateItem(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	//decode json request to item
	err := json.NewDecoder(r.Body).Decode(&item)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request"))
		log.Println(err.Error())
		return
	}

	insertID, err := repo.Repository.InsertItem(item)
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

//UpdateItemStatus send response after item update
func (repo *Repository) UpdateItemStatus(w http.ResponseWriter, r *http.Request) {
	//get item id from the request params, and convert it to string
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request"))
		log.Println(err.Error())
		return
	}

	updatedRows, err := repo.Repository.UpdateItemStatus(int64(id))
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

//UpdateAllItemsStatus send request after all items update
func (repo *Repository) UpdateAllItemsStatus(w http.ResponseWriter, r *http.Request) {
	updatedRows, err := repo.Repository.UpdateAllItemsStatus()
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

//DeleteItem send response after item remove
func (repo *Repository) DeleteItem(w http.ResponseWriter, r *http.Request) {
	//get item id from the request params and convert int to string
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request"))
		log.Println(err.Error())
		return
	}

	deletedRows, err := repo.Repository.DeleteItem(int64(id))
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

//DeleteAllItems send response after all item removed
func (repo *Repository) DeleteAllItems(w http.ResponseWriter, r *http.Request) {
	err := repo.Repository.DeleteAllItems()
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
