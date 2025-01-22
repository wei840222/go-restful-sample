package storage

import (
	"context"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Title       string
	Description string
	Completed   bool
}

type TodoStorage interface {
	Create(ctx context.Context, todo *Todo) error
	List(ctx context.Context) ([]*Todo, error)
	Get(ctx context.Context, id uint) (*Todo, error)
	Update(ctx context.Context, id uint, todo *Todo) error
	Delete(ctx context.Context, id uint) error
}

type todoStorage struct {
	db *gorm.DB
}

func NewTodoStorage(lc fx.Lifecycle, db *gorm.DB) TodoStorage {

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return db.AutoMigrate(&Todo{})
		},
	})

	return &todoStorage{db: db}
}

func (s *todoStorage) Create(ctx context.Context, todo *Todo) error {
	return s.db.WithContext(ctx).Create(todo).Error
}

func (s *todoStorage) List(ctx context.Context) ([]*Todo, error) {
	var todos []*Todo
	return todos, s.db.WithContext(ctx).Find(&todos).Error
}

func (s *todoStorage) Get(ctx context.Context, id uint) (*Todo, error) {
	var todo Todo
	return &todo, s.db.WithContext(ctx).First(&todo, id).Error
}

func (s *todoStorage) Update(ctx context.Context, id uint, todo *Todo) error {
	return s.db.WithContext(ctx).Model(&Todo{}).Where("id = ?", id).Updates(todo).Error
}

func (s *todoStorage) Delete(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&Todo{}, id).Error
}
