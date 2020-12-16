package store

import (
	"context"

	"gorm.io/gorm"
)

// User orm struct
type User struct {
	gorm.Model
	Name string `gorm:"not null"`
}

//go:generate mockgen -destination=./mock/user_mock.go -package=mock . UserStore
// UserStore repository interface for user
type UserStore interface {
	List(ctx context.Context) ([]*User, error)
	Get(ctx context.Context, id uint) (*User, error)
	Create(ctx context.Context, user *User) (*User, error)
	Update(ctx context.Context, id uint, user *User) error
	Delete(ctx context.Context, id uint) error
}

// NewUserStore factory methord for UserStore
func NewUserStore(db *gorm.DB) (UserStore, error) {
	if err := db.AutoMigrate(&User{}); err != nil {
		return nil, err
	}
	return &userStore{db}, nil
}

type userStore struct {
	db *gorm.DB
}

func (s userStore) List(ctx context.Context) ([]*User, error) {
	var users []*User
	if err := s.db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (s userStore) Get(ctx context.Context, id uint) (*User, error) {
	var user User
	if err := s.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s userStore) Create(ctx context.Context, user *User) (*User, error) {
	if err := s.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (s userStore) Update(ctx context.Context, id uint, user *User) error {
	return s.db.WithContext(ctx).Model(user).Where("id = ?", id).Updates(user).Error
}

func (s userStore) Delete(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&User{}, id).Error
}
