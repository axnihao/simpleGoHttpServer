package server

import (
	"bookstore/server/middleware"
	"bookstore/store"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

const (
	ContentType     = "Content-Type"
	ApplicationJson = "application/json"
)

type BookStoreServer struct {
	s   store.Store
	srv *http.Server
}

func (bs *BookStoreServer) createBookHandler(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var book store.Book
	if err := decoder.Decode(&book); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if err := bs.s.Create(&book); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

}

func (bs *BookStoreServer) updateBookHandler(writer http.ResponseWriter, request *http.Request) {
	id, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, "no id found in request", http.StatusBadRequest)
	}

	decoder := json.NewDecoder(request.Body)
	var book store.Book
	if err := decoder.Decode(&book); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	book.Id = id
	if err := bs.s.Update(&book); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
}

func (bs *BookStoreServer) getBookHandler(writer http.ResponseWriter, request *http.Request) {
	id, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, "no id found in request", http.StatusBadRequest)
	}
	book, err := bs.s.Get(id)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	response(writer, book)
}

func (bs *BookStoreServer) getAllBookHandler(writer http.ResponseWriter, request *http.Request) {
	books, _ := bs.s.GetAll()
	response(writer, books)
}

func (bs *BookStoreServer) deleteBookHandler(writer http.ResponseWriter, request *http.Request) {
	id, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, "no id found in request", http.StatusBadRequest)
		return
	}
	if bs.s.Delete(id) != nil {
		http.Error(writer, "no book found in request", http.StatusBadRequest)
	}
}

func NewBookStoreServer(addr string, s store.Store) *BookStoreServer {
	srv := &BookStoreServer{
		s: s,
		srv: &http.Server{
			Addr: addr,
		},
	}
	router := mux.NewRouter()
	router.HandleFunc("/book", srv.createBookHandler).Methods(http.MethodPost)
	router.HandleFunc("/book/{id}", srv.updateBookHandler).Methods(http.MethodPost)
	router.HandleFunc("/book/{id}", srv.getBookHandler).Methods(http.MethodGet)
	router.HandleFunc("/book", srv.getAllBookHandler).Methods(http.MethodGet)
	router.HandleFunc("/book/{id}", srv.deleteBookHandler).Methods(http.MethodDelete)

	srv.srv.Handler = middleware.Logging(middleware.Validating(router))
	return srv

}

func (bs *BookStoreServer) ListenAndServe() (<-chan error, error) {
	var err error
	errChan := make(chan error)
	go func() {
		err = bs.srv.ListenAndServe()
		errChan <- err
	}()

	select {
	case err = <-errChan:
		return nil, err
	case <-time.After(time.Second):
		return errChan, err
	}

}

func (bs *BookStoreServer) Shutdown(ctx context.Context) error {
	return bs.srv.Shutdown(ctx)
}

func response(writer http.ResponseWriter, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	writer.Header().Set(ContentType, ApplicationJson)
	writer.Write(data)
}
