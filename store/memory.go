package store

import "sync"

func init() {
	Register("mem", &MemStore{
		books: make(map[int64]Book),
	})
}

type MemStore struct {
	sync.RWMutex
	books map[int64]Book
}

func (s *MemStore) Create(book *Book) (err error) {
	s.RLock()
	if _, exist := s.books[book.Id]; exist {
		s.RUnlock()
		return ErrExist
	}
	s.RUnlock()

	s.Lock()
	defer s.Unlock()
	s.books[book.Id] = *book
	return
}

func (s *MemStore) Update(book *Book) (err error) {
	s.RLock()
	obj, exist := s.books[book.Id]
	s.RUnlock()
	if !exist {
		return ErrNotFound
	}
	if book.Name != "" {
		obj.Name = book.Name
	}
	if len(book.Authors) <= 0 {
		obj.Authors = book.Authors
	}
	if book.Press != "" {
		obj.Press = book.Press
	}

	s.Lock()
	defer s.Unlock()
	s.books[book.Id] = obj
	return
}

func (s *MemStore) Find(id int64) (obj Book, err error) {
	s.RLock()
	obj, exist := s.books[id]
	s.RUnlock()
	if !exist {
		err = ErrNotFound
		return
	}
	return
}

func (s *MemStore) GetAll() (objs []Book, count int64, err error) {
	s.RLock()
	defer s.RUnlock()

	// 不是直接 objs = s.books
	for _, obj := range s.books {
		objs = append(objs, obj)
	}
	count = int64(len(s.books))
	return
}

func (s *MemStore) Delete(id int64) (err error) {
	s.RLock()
	_, exist := s.books[id]
	s.RUnlock()
	if !exist {
		err = ErrNotFound
		return
	}
	s.Lock()
	defer s.Unlock()
	delete(s.books, id)
	return
}
