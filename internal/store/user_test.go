package store_test

import (
	"github.com/wei840222/go-restful-sample/internal/store"

	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUserStore(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	userStore, _ := store.NewUserStore(db)

	newUser, err := userStore.Create(context.Background(), &store.User{Name: "tester"})
	assert.Nil(t, err)
	assert.NotZero(t, newUser.ID)
	assert.EqualValues(t, "tester", newUser.Name)
	assert.NotZero(t, newUser.CreatedAt)

	users, err := userStore.List(context.Background())
	assert.Nil(t, err)
	assert.Len(t, users, 1)
	assert.EqualValues(t, newUser.ID, users[0].ID)
	assert.EqualValues(t, "tester", users[0].Name)

	user, err := userStore.Get(context.Background(), newUser.ID)
	assert.Nil(t, err)
	assert.EqualValues(t, "tester", user.Name)

	err = userStore.Update(context.Background(), newUser.ID, &store.User{Name: "tester2"})
	assert.Nil(t, err)
	user, _ = userStore.Get(context.Background(), newUser.ID)
	assert.EqualValues(t, "tester2", user.Name)

	err = userStore.Delete(context.Background(), newUser.ID)
	assert.Nil(t, err)
	_, err = userStore.Get(context.Background(), newUser.ID)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}
