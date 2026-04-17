package store

import (
	"errors"
	"sort"
	"strings"
	"sync"
	"time"

	"bookstore/models"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrDuplicateName = errors.New("duplicate name")
	ErrInvalidData   = errors.New("invalid data")
)

type Store struct {
	mu sync.RWMutex

	books      map[int]models.Book
	authors    map[int]models.Author
	categories map[int]models.Category

	// favorites[userID][bookID] = FavoriteBook
	favorites map[int]map[int]models.FavoriteBook

	nextBookID      int
	nextAuthorID    int
	nextCategoryID  int
}

func NewStore() *Store {
	return &Store{
		books:      make(map[int]models.Book),
		authors:    make(map[int]models.Author),
		categories: make(map[int]models.Category),
		favorites:  make(map[int]map[int]models.FavoriteBook),
		nextBookID:     1,
		nextAuthorID:   1,
		nextCategoryID: 1,
	}
}

func (s *Store) SeedDemoData() {
	author1, _ := s.CreateAuthor("George Orwell")
	author2, _ := s.CreateAuthor("J. K. Rowling")

	category1, _ := s.CreateCategory("Fiction")
	category2, _ := s.CreateCategory("Fantasy")

	_, _ = s.CreateBook(models.BookCreateRequest{
		Title:      "1984",
		AuthorID:   author1.ID,
		CategoryID: category1.ID,
		Price:      12.50,
	})

	_, _ = s.CreateBook(models.BookCreateRequest{
		Title:      "Animal Farm",
		AuthorID:   author1.ID,
		CategoryID: category1.ID,
		Price:      10.00,
	})

	_, _ = s.CreateBook(models.BookCreateRequest{
		Title:      "Harry Potter and the Sorcerer's Stone",
		AuthorID:   author2.ID,
		CategoryID: category2.ID,
		Price:      19.99,
	})
}

func containsFold(s, sub string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(sub))
}

func (s *Store) CreateAuthor(name string) (models.Author, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return models.Author{}, ErrInvalidData
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, a := range s.authors {
		if strings.EqualFold(a.Name, name) {
			return models.Author{}, ErrDuplicateName
		}
	}

	author := models.Author{
		ID:   s.nextAuthorID,
		Name: name,
	}
	s.authors[author.ID] = author
	s.nextAuthorID++

	return author, nil
}

func (s *Store) ListAuthors() []models.Author {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]models.Author, 0, len(s.authors))
	for _, a := range s.authors {
		result = append(result, a)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})

	return result
}

func (s *Store) GetAuthor(id int) (models.Author, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	a, ok := s.authors[id]
	return a, ok
}

func (s *Store) CreateCategory(name string) (models.Category, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return models.Category{}, ErrInvalidData
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, c := range s.categories {
		if strings.EqualFold(c.Name, name) {
			return models.Category{}, ErrDuplicateName
		}
	}

	category := models.Category{
		ID:   s.nextCategoryID,
		Name: name,
	}
	s.categories[category.ID] = category
	s.nextCategoryID++

	return category, nil
}

func (s *Store) ListCategories() []models.Category {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]models.Category, 0, len(s.categories))
	for _, c := range s.categories {
		result = append(result, c)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})

	return result
}

func (s *Store) GetCategory(id int) (models.Category, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	c, ok := s.categories[id]
	return c, ok
}

func (s *Store) CreateBook(req models.BookCreateRequest) (models.Book, error) {
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" || req.AuthorID <= 0 || req.CategoryID <= 0 || req.Price < 0.01 {
		return models.Book{}, ErrInvalidData
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.authors[req.AuthorID]; !ok {
		return models.Book{}, ErrNotFound
	}
	if _, ok := s.categories[req.CategoryID]; !ok {
		return models.Book{}, ErrNotFound
	}

	book := models.Book{
		ID:         s.nextBookID,
		Title:      req.Title,
		AuthorID:   req.AuthorID,
		CategoryID: req.CategoryID,
		Price:      req.Price,
	}
	s.books[book.ID] = book
	s.nextBookID++

	return book, nil
}

func (s *Store) UpdateBook(id int, req models.BookUpdateRequest) (models.Book, error) {
	req.Title = strings.TrimSpace(req.Title)
	if id <= 0 || req.Title == "" || req.AuthorID <= 0 || req.CategoryID <= 0 || req.Price < 0.01 {
		return models.Book{}, ErrInvalidData
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.books[id]; !ok {
		return models.Book{}, ErrNotFound
	}
	if _, ok := s.authors[req.AuthorID]; !ok {
		return models.Book{}, ErrNotFound
	}
	if _, ok := s.categories[req.CategoryID]; !ok {
		return models.Book{}, ErrNotFound
	}

	book := models.Book{
		ID:         id,
		Title:      req.Title,
		AuthorID:   req.AuthorID,
		CategoryID: req.CategoryID,
		Price:      req.Price,
	}
	s.books[id] = book

	return book, nil
}

func (s *Store) DeleteBook(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.books[id]; !ok {
		return false
	}

	delete(s.books, id)

	for userID := range s.favorites {
		delete(s.favorites[userID], id)
		if len(s.favorites[userID]) == 0 {
			delete(s.favorites, userID)
		}
	}

	return true
}

func (s *Store) GetBookByID(id int) (models.Book, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	book, ok := s.books[id]
	return book, ok
}

func (s *Store) bookViewLocked(book models.Book) (models.BookView, bool) {
	author, ok := s.authors[book.AuthorID]
	if !ok {
		return models.BookView{}, false
	}
	category, ok := s.categories[book.CategoryID]
	if !ok {
		return models.BookView{}, false
	}

	return models.BookView{
		ID:           book.ID,
		Title:        book.Title,
		AuthorID:     author.ID,
		AuthorName:   author.Name,
		CategoryID:   category.ID,
		CategoryName: category.Name,
		Price:        book.Price,
	}, true
}

func (s *Store) GetBookView(id int) (models.BookView, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	book, ok := s.books[id]
	if !ok {
		return models.BookView{}, false
	}

	return s.bookViewLocked(book)
}

func (s *Store) ListBooks(filter models.BookFilter) ([]models.BookView, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]models.BookView, 0, len(s.books))

	for _, book := range s.books {
		view, ok := s.bookViewLocked(book)
		if !ok {
			continue
		}

		if filter.Title != "" && !containsFold(view.Title, filter.Title) {
			continue
		}
		if filter.Author != "" && !containsFold(view.AuthorName, filter.Author) {
			continue
		}
		if filter.Category != "" && !containsFold(view.CategoryName, filter.Category) {
			continue
		}
		if filter.MinPrice > 0 && view.Price < filter.MinPrice {
			continue
		}
		if filter.MaxPrice > 0 && view.Price > filter.MaxPrice {
			continue
		}

		items = append(items, view)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})

	total := len(items)

	page := filter.Page
	limit := filter.Limit
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	start := (page - 1) * limit
	if start >= total {
		return []models.BookView{}, total
	}

	end := start + limit
	if end > total {
		end = total
	}

	return items[start:end], total
}

func (s *Store) AddFavorite(userID, bookID int) (models.FavoriteBook, error) {
	if userID <= 0 || bookID <= 0 {
		return models.FavoriteBook{}, ErrInvalidData
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.books[bookID]; !ok {
		return models.FavoriteBook{}, ErrNotFound
	}

	if _, ok := s.favorites[userID]; !ok {
		s.favorites[userID] = make(map[int]models.FavoriteBook)
	}

	if existing, ok := s.favorites[userID][bookID]; ok {
		return existing, nil
	}

	fav := models.FavoriteBook{
		UserID:    userID,
		BookID:    bookID,
		CreatedAt: time.Now(),
	}

	s.favorites[userID][bookID] = fav
	return fav, nil
}

func (s *Store) RemoveFavorite(userID, bookID int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	userFavorites, ok := s.favorites[userID]
	if !ok {
		return false
	}

	if _, ok := userFavorites[bookID]; !ok {
		return false
	}

	delete(userFavorites, bookID)
	if len(userFavorites) == 0 {
		delete(s.favorites, userID)
	}

	return true
}

func (s *Store) ListFavorites(userID, page, limit int) ([]models.FavoriteBookView, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	userFavorites, ok := s.favorites[userID]
	if !ok {
		return []models.FavoriteBookView{}, 0
	}

	items := make([]models.FavoriteBookView, 0, len(userFavorites))

	for _, fav := range userFavorites {
		book, ok := s.books[fav.BookID]
		if !ok {
			continue
		}

		view, ok := s.bookViewLocked(book)
		if !ok {
			continue
		}

		items = append(items, models.FavoriteBookView{
			ID:           view.ID,
			Title:        view.Title,
			AuthorID:     view.AuthorID,
			AuthorName:   view.AuthorName,
			CategoryID:   view.CategoryID,
			CategoryName: view.CategoryName,
			Price:        view.Price,
			CreatedAt:    fav.CreatedAt,
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt.After(items[j].CreatedAt)
	})

	total := len(items)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	start := (page - 1) * limit
	if start >= total {
		return []models.FavoriteBookView{}, total
	}

	end := start + limit
	if end > total {
		end = total
	}

	return items[start:end], total
}