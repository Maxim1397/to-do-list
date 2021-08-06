package main

import (
	"github.com/gorilla/mux" //used for routes
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"to-do-list/handlers"
)

func TestGetAllItems(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/items", handlers.GetAllItems)
	ts := httptest.NewServer(r)
	defer ts.Close()
	res, err := http.Get(ts.URL + "/items")
	if err != nil {
		t.Errorf("Expected nil, received %s", err.Error())
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected %d, received %d", http.StatusOK, res.StatusCode)
	}

}

func TestGetItem(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/items/{id}", handlers.GetItem)
	ts := httptest.NewServer(r)
	defer ts.Close()
	t.Run("not found", func(t *testing.T) {
		res, err := http.Get(ts.URL + "/items/1")
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
		}

		if res.StatusCode != http.StatusNotFound {
			t.Errorf("Expected %d, received %d", http.StatusNotFound, res.StatusCode)
		}

	})

	t.Run("found", func(t *testing.T) {
		res, err := http.Get(ts.URL + "/items/70")
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
		}

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected %d, received %d", http.StatusOK, res.StatusCode)
		}
	})
}

func TestCreateItem(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/items", handlers.CreateItem)
	ts := httptest.NewServer(r)
	defer ts.Close()
	res, err := http.Post(ts.URL+"/items", "application/json", strings.NewReader(`{"description":"Drink water"}`))
	if err != nil {
		t.Errorf("Expected nil, received %s", err.Error())
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected %d, received %d", http.StatusOK, res.StatusCode)
	}

}

func TestUpdateItemStatus(t *testing.T) {
	client := &http.Client{}
	r := mux.NewRouter()
	r.HandleFunc("/items/"+"{id}", handlers.UpdateItemStatus)
	ts := httptest.NewServer(r)
	defer ts.Close()
	t.Run("not found", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPut, ts.URL+"/items/1", nil)
		if err != nil {
			t.Errorf(err.Error())
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Errorf(err.Error())
			return
		}

		defer resp.Body.Close()
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected %d, received %d", http.StatusNotFound, resp.StatusCode)
		}
	})
	t.Run("found", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPut, ts.URL+"/items/60", nil)
		if err != nil {
			t.Errorf(err.Error())
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Errorf(err.Error())
			return
		}

		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected %d, received %d", http.StatusOK, resp.StatusCode)
		}
	})
}

func TestUpdateAllItemsStatus(t *testing.T) {
	client := &http.Client{}
	r := mux.NewRouter()
	r.HandleFunc("/items", handlers.UpdateAllItemsStatus)
	ts := httptest.NewServer(r)
	defer ts.Close()
	req, err := http.NewRequest(http.MethodPut, ts.URL+"/items", nil)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected %d, received %d", http.StatusOK, resp.StatusCode)
	}
}

func TestDeleteItem(t *testing.T) {
	client := &http.Client{}
	r := mux.NewRouter()
	r.HandleFunc("/items/"+"{id}", handlers.DeleteItem)
	ts := httptest.NewServer(r)
	defer ts.Close()
	t.Run("not found", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, ts.URL+"/items/60", nil)
		if err != nil {
			t.Errorf(err.Error())
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Errorf(err.Error())
			return
		}

		defer resp.Body.Close()
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected %d, received %d", http.StatusNotFound, resp.StatusCode)
		}
	})
	t.Run("found", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, ts.URL+"/items/58", nil)
		if err != nil {
			t.Errorf(err.Error())
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Errorf(err.Error())
			return
		}

		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected %d, received %d", http.StatusOK, resp.StatusCode)
		}
	})
}

func TestDeleteAllItems(t *testing.T) {
	client := &http.Client{}
	r := mux.NewRouter()
	r.HandleFunc("/items", handlers.DeleteAllItems)
	ts := httptest.NewServer(r)
	defer ts.Close()
	req, err := http.NewRequest(http.MethodDelete, ts.URL+"/items", nil)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected %d, received %d", http.StatusOK, resp.StatusCode)
	}
}
