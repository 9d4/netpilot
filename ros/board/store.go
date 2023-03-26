package board

import "gorm.io/gorm"

type Store interface {
	Create(board *Board) error
	FindAll() ([]*Board, error)
	FindByID(id uint) (*Board, error)
	FindByUUID(uuid string) (*Board, error)
	Update(board *Board) error
	Delete(board *Board) error
	GetDB() *gorm.DB
}
type store struct {
	db *gorm.DB
}

func NewBoardStore(db *gorm.DB) Store {
	return &store{db: db}
}

func (s *store) Create(board *Board) error {
	return s.db.Create(board).Error
}

func (s *store) FindAll() ([]*Board, error) {
	var boards []*Board
	err := s.db.Find(&boards).Error
	return boards, err
}

func (s *store) FindByID(id uint) (*Board, error) {
	var board Board
	err := s.db.First(&board, id).Error
	if err != nil {
		return nil, err
	}
	return &board, nil
}

func (s *store) FindByUUID(uuid string) (*Board, error) {
	var board Board
	err := s.db.Where("uuid = ?", uuid).First(&board).Error
	if err != nil {
		return nil, err
	}
	return &board, nil
}

func (s *store) Update(board *Board) error {
	return s.db.Save(board).Error
}

func (s *store) Delete(board *Board) error {
	return s.db.Delete(board).Error
}

func (s *store) GetDB() *gorm.DB {
	return s.db
}
