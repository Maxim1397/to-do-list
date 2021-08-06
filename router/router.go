package router

import (
	"github.com/gorilla/mux"
	"to-do-list/handlers"
)

func Router() *mux.Router {
	//create routes
	router := mux.NewRouter()
	router.HandleFunc("/items", handlers.GetAllItems).Methods("GET", "OPTIONS")
	router.HandleFunc("/items/{id}", handlers.GetItem).Methods("GET", "OPTIONS")
	router.HandleFunc("/items", handlers.CreateItem).Methods("POST", "OPTIONS")
	router.HandleFunc("/items/{id}", handlers.UpdateItemStatus).Methods("PUT", "OPTIONS")
	router.HandleFunc("/items", handlers.UpdateAllItemsStatus).Methods("PUT", "OPTIONS")
	router.HandleFunc("/items/{id}", handlers.DeleteItem).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/items", handlers.DeleteAllItems).Methods("DELETE", "OPTIONS")
	return router
}
