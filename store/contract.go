package store

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrExist = errors.New("exist")
)
type Book struct {
	Id      int64   `json:"id"`      // 图书ISBN ID
	Name    string   `json:"name"`    // 图书名称
	Authors []string `json:"authors"` // 图书作者
	Press   string   `json:"press"`   // 出版社
}


type Store interface {
	Create(book *Book)error
	Update(book *Book)error
	Find(int64)(Book,error)
	GetAll()([]Book,int64,error)
	Delete(int64)error
}
