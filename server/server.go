package server

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/spf13/cast"
	"github.com/wx-satellite/bookstore/store"
	"net/http"
	"time"
)

type Server struct {
	store  store.Store
	server *http.Server
}

func New(addr string, store store.Store) (srv *Server) {
	srv = &Server{
		store:  store,
		server: &http.Server{Addr: addr},
	}
	router := mux.NewRouter()

	// 注册路由

	srv.server.Handler = router
	return
}

// ListenAndServe 启动服务
func (s *Server) ListenAndServe() (ch <-chan error, err error) {
	errCh := make(chan error)
	go func() {
		err = s.server.ListenAndServe()
		errCh <- err
	}()
	ticker := time.NewTicker(time.Second)
	select {
	case err = <-errCh:
		ticker.Stop()
		return
	case <-ticker.C:
		return errCh, nil
	}
}

// Shutdown 优雅退出
func (s *Server) Shutdown(ctx context.Context) (err error) {
	return s.server.Shutdown(ctx)
}

// createBookHandler 创建
func (s *Server) createBookHandler(w http.ResponseWriter, r *http.Request) {
	var book store.Book
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.store.Create(&book); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// updateBookHandler 更新
func (s *Server) updateBookHandler(w http.ResponseWriter, r *http.Request) {
	id := cast.ToInt64(mux.Vars(r)["id"])
	if id <= 0 {
		http.Error(w, "no id found in request", http.StatusBadRequest)
		return
	}
	var book store.Book
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	book.Id = id
	if err := s.store.Update(&book); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// findBookHandler 获取单个
func (s *Server) findBookHandler(w http.ResponseWriter, r *http.Request) {
	id := cast.ToInt64(mux.Vars(r)["id"])
	if id <= 0 {
		http.Error(w, "no id found in request", http.StatusBadRequest)
		return
	}
	data, err := s.store.Find(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response(w, data)
}

// getBooksHandler 获取所有
func (s *Server) getBooksHandler(w http.ResponseWriter, r *http.Request) {
	books, _, err := s.store.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response(w, books)
	return
}

// deleteBookHandler 删除
func (s *Server) deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	id := cast.ToInt64(mux.Vars(r)["id"])
	if id <= 0 {
		http.Error(w, "no id found in request", http.StatusBadRequest)
		return
	}
	err := s.store.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func response(w http.ResponseWriter, data interface{}) {
	dataBs, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	_, _ = w.Write(dataBs)
}
