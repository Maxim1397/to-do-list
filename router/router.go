package router

import (
	"github.com/gorilla/mux"
	"to-do-list/handlers"
)

func Router() *mux.Router {
	//create routes
	router := mux.NewRouter()
	//create handler
	handler := handlers.NewHandler()
	router.HandleFunc("/items", handler.GetAllItems).Methods("GET", "OPTIONS")
	router.HandleFunc("/items/{id}", handler.GetItem).Methods("GET", "OPTIONS")
	router.HandleFunc("/items", handler.CreateItem).Methods("POST", "OPTIONS")
	router.HandleFunc("/items/{id}", handler.UpdateItemStatus).Methods("PUT", "OPTIONS")
	router.HandleFunc("/items", handler.UpdateAllItemsStatus).Methods("PUT", "OPTIONS")
	router.HandleFunc("/items/{id}", handler.DeleteItem).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/items", handler.DeleteAllItems).Methods("DELETE", "OPTIONS")
	return router
}
