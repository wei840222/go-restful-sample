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
	Completed   *bool `gorm:"default:false"`
}

type TodoStorage interface {
	Get(ctx context.Context, id int) (Todo, error)
	List(ctx context.Context) ([]Todo, error)
	Create(ctx context.Context, todo *Todo) error
	Update(ctx context.Context, id int, todo Todo) error
	Delete(ctx context.Context, id int) error
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

func (s *todoStorage) List(ctx context.Context) ([]Todo, error) {
	var todos []Todo
	if err := s.db.WithContext(ctx).Find(&todos).Error; err != nil {
		return nil, err
	}
	return todos, nil
}

func (s *todoStorage) Get(ctx context.Context, id int) (Todo, error) {
	var todo Todo
	if err := s.db.WithContext(ctx).First(&todo, id).Error; err != nil {
		return todo, err
	}
	return todo, nil
}

func (s *todoStorage) Update(ctx context.Context, id int, todo Todo) error {
	if _, err := s.Get(ctx, id); err != nil {
		return err
	}
	return s.db.WithContext(ctx).Model(&Todo{}).Where("id = ?", id).Updates(todo).Error
}

func (s *todoStorage) Delete(ctx context.Context, id int) error {
	if _, err := s.Get(ctx, id); err != nil {
		return err
	}
	return s.db.WithContext(ctx).Delete(&Todo{}, id).Error
}
