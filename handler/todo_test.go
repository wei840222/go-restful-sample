package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"

	"github.com/wei840222/go-restful-sample/storage"
	"github.com/wei840222/go-restful-sample/storage/mock"
)

func TestTodoHandler_Get(t *testing.T) {
	gin.SetMode(gin.TestMode)

	Convey("Given a TodoHandler with mock storage", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStorage := mock.NewMockTodoStorage(ctrl)
		e := gin.Default()
		RegisterTodoHandler(e, mockStorage)

		Convey("When getting a todo with valid ID", func() {
			now := time.Now()
			completed := false
			mockTodo := storage.Todo{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Title:       "Test Todo",
				Description: "Test Description",
				Completed:   &completed,
			}

			mockStorage.EXPECT().
				Get(gomock.Any(), gomock.Eq(1)).
				Return(mockTodo, nil).
				Times(1)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/todos/1", nil)
			e.ServeHTTP(w, req)

			Convey("Then it should return 200 status code", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
			})

			Convey("And return the todo", func() {
				var res GetTodoRes
				err := json.Unmarshal(w.Body.Bytes(), &res)
				So(err, ShouldBeNil)
				So(res.ID, ShouldEqual, 1)
				So(res.Title, ShouldEqual, "Test Todo")
				So(res.Description, ShouldEqual, "Test Description")
				So(res.Completed, ShouldEqual, false)
				So(res.CreatedAt, ShouldEqual, now)
				So(res.UpdatedAt, ShouldEqual, now)
			})
		})

		Convey("When getting a todo with invalid ID format", func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/todos/invalid", nil)
			e.ServeHTTP(w, req)

			Convey("Then it should return 400 status code", func() {
				So(w.Code, ShouldEqual, http.StatusBadRequest)
			})
		})

		Convey("When todo is not found", func() {
			mockStorage.EXPECT().
				Get(gomock.Any(), gomock.Eq(999)).
				Return(storage.Todo{}, gorm.ErrRecordNotFound).
				Times(1)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/todos/999", nil)
			e.ServeHTTP(w, req)

			Convey("Then it should return 404 status code", func() {
				So(w.Code, ShouldEqual, http.StatusNotFound)
			})
		})

		Convey("When storage returns an error", func() {
			mockStorage.EXPECT().
				Get(gomock.Any(), gomock.Any()).
				Return(storage.Todo{}, errors.New("internal server error")).
				Times(1)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/todos/1", nil)
			e.ServeHTTP(w, req)

			Convey("Then it should return 500 status code", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})
}

func TestTodoHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)

	Convey("Given a TodoHandler with mock storage", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStorage := mock.NewMockTodoStorage(ctrl)
		e := gin.Default()
		RegisterTodoHandler(e, mockStorage)

		Convey("When listing todos successfully", func() {
			now := time.Now()
			completed := false
			notCompleted := true
			mockTodos := []storage.Todo{
				{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: now,
						UpdatedAt: now,
					},
					Title:       "First Todo",
					Description: "First Description",
					Completed:   &completed,
				},
				{
					Model: gorm.Model{
						ID:        2,
						CreatedAt: now,
						UpdatedAt: now,
					},
					Title:       "Second Todo",
					Description: "Second Description",
					Completed:   &notCompleted,
				},
			}

			mockStorage.EXPECT().
				List(gomock.Any()).
				Return(mockTodos, nil).
				Times(1)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/todos", nil)
			e.ServeHTTP(w, req)

			Convey("Then it should return 200 status code", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
			})

			Convey("And return the list of todos", func() {
				var res ListTodoRes
				err := json.Unmarshal(w.Body.Bytes(), &res)
				So(err, ShouldBeNil)
				So(len(res), ShouldEqual, 2)

				So(res[0].ID, ShouldEqual, 1)
				So(res[0].Title, ShouldEqual, "First Todo")
				So(res[0].Description, ShouldEqual, "First Description")
				So(res[0].Completed, ShouldEqual, false)
				So(res[0].CreatedAt, ShouldEqual, now)
				So(res[0].UpdatedAt, ShouldEqual, now)

				So(res[1].ID, ShouldEqual, 2)
				So(res[1].Title, ShouldEqual, "Second Todo")
				So(res[1].Description, ShouldEqual, "Second Description")
				So(res[1].Completed, ShouldEqual, true)
				So(res[1].CreatedAt, ShouldEqual, now)
				So(res[1].UpdatedAt, ShouldEqual, now)
			})
		})

		Convey("When listing todos with empty result", func() {
			mockStorage.EXPECT().
				List(gomock.Any()).
				Return([]storage.Todo{}, nil).
				Times(1)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/todos", nil)
			e.ServeHTTP(w, req)

			Convey("Then it should return 200 status code", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
			})

			Convey("And return empty array", func() {
				var res ListTodoRes
				err := json.Unmarshal(w.Body.Bytes(), &res)
				So(err, ShouldBeNil)
				So(len(res), ShouldEqual, 0)
				So(w.Body.String(), ShouldEqual, "[]")
			})
		})

		Convey("When storage returns an error", func() {
			mockStorage.EXPECT().
				List(gomock.Any()).
				Return(nil, errors.New("internal server error")).
				Times(1)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/todos", nil)
			e.ServeHTTP(w, req)

			Convey("Then it should return 500 status code", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})
}

func TestTodoHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	Convey("Given a TodoHandler with mock storage", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStorage := mock.NewMockTodoStorage(ctrl)
		e := gin.Default()
		RegisterTodoHandler(e, mockStorage)

		Convey("When creating a new todo with valid input", func() {
			now := time.Now()

			mockStorage.EXPECT().
				Create(gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ any, todo *storage.Todo) error {
					completed := false
					todo.ID = 1
					todo.Completed = &completed
					todo.CreatedAt = now
					todo.UpdatedAt = now
					return nil
				}).
				Times(1)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/todos", bytes.NewBufferString(`{
																							"title": "Test Todo",
																							"description": "Test Description"
																						}`))
			req.Header.Set("Content-Type", "application/json")
			e.ServeHTTP(w, req)

			Convey("Then it should return 201 status code", func() {
				So(w.Code, ShouldEqual, http.StatusCreated)
			})

			Convey("And return the created todo", func() {
				var res GetTodoRes
				err := json.Unmarshal(w.Body.Bytes(), &res)
				So(err, ShouldBeNil)
				So(res.ID, ShouldEqual, 1)
				So(res.Title, ShouldEqual, "Test Todo")
				So(res.Description, ShouldEqual, "Test Description")
				So(res.Completed, ShouldEqual, false)
			})
		})

		Convey("When creating a todo with invalid JSON", func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/todos", bytes.NewBufferString("{invalid json}"))
			req.Header.Set("Content-Type", "application/json")
			e.ServeHTTP(w, req)

			Convey("Then it should return 400 status code", func() {
				So(w.Code, ShouldEqual, http.StatusBadRequest)
			})
		})

		Convey("When storage returns an error", func() {
			mockStorage.EXPECT().
				Create(gomock.Any(), gomock.Any()).
				Return(errors.New("internal server error")).
				Times(1)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/todos", bytes.NewBufferString(`{
																							"title": "Test Todo",
																							"description": "Test Description"
																						}`))
			req.Header.Set("Content-Type", "application/json")
			e.ServeHTTP(w, req)

			Convey("Then it should return 500 status code", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})
}

func TestTodoHandler_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)

	Convey("Given a TodoHandler with mock storage", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStorage := mock.NewMockTodoStorage(ctrl)
		e := gin.Default()
		RegisterTodoHandler(e, mockStorage)

		Convey("When updating a todo with valid input", func() {
			mockStorage.EXPECT().
				Update(gomock.Any(), gomock.Eq(1), gomock.Any()).
				DoAndReturn(func(_ any, _ int, todo storage.Todo) error {
					So(todo.Title, ShouldEqual, "Updated Todo")
					So(todo.Description, ShouldEqual, "Updated Description")
					So(*todo.Completed, ShouldEqual, true)
					return nil
				}).
				Times(1)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPatch, "/todos/1", bytes.NewBufferString(`{
																							"title": "Updated Todo",
																							"description": "Updated Description",
																							"completed": true
																						}`))
			req.Header.Set("Content-Type", "application/json")
			e.ServeHTTP(w, req)

			Convey("Then it should return 204 status code", func() {
				So(w.Code, ShouldEqual, http.StatusNoContent)
			})
		})

		Convey("When updating a todo with invalid ID format", func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPatch, "/todos/invalid", bytes.NewBufferString(`{
																									"title": "Updated Todo",
																									"description": "Updated Description",
																									"completed": true
																								}`))
			req.Header.Set("Content-Type", "application/json")
			e.ServeHTTP(w, req)

			Convey("Then it should return 400 status code", func() {
				So(w.Code, ShouldEqual, http.StatusBadRequest)
			})
		})

		Convey("When updating with invalid JSON", func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPatch, "/todos/1", bytes.NewBufferString("{invalid json}"))
			req.Header.Set("Content-Type", "application/json")
			e.ServeHTTP(w, req)

			Convey("Then it should return 400 status code", func() {
				So(w.Code, ShouldEqual, http.StatusBadRequest)
			})
		})

		Convey("When todo is not found", func() {
			mockStorage.EXPECT().
				Update(gomock.Any(), gomock.Eq(999), gomock.Any()).
				Return(gorm.ErrRecordNotFound).
				Times(1)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPatch, "/todos/999", bytes.NewBufferString(`{
																								"title": "Updated Todo",
																								"description": "Updated Description",
																								"completed": true
																							}`))
			req.Header.Set("Content-Type", "application/json")
			e.ServeHTTP(w, req)

			Convey("Then it should return 404 status code", func() {
				So(w.Code, ShouldEqual, http.StatusNotFound)
			})
		})

		Convey("When storage returns an error", func() {
			mockStorage.EXPECT().
				Update(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(errors.New("internal server error")).
				Times(1)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPatch, "/todos/1", bytes.NewBufferString(`{
																							"title": "Updated Todo",
																							"description": "Updated Description",
																							"completed": true
																						}`))
			req.Header.Set("Content-Type", "application/json")
			e.ServeHTTP(w, req)

			Convey("Then it should return 500 status code", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})
}

func TestTodoHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	Convey("Given a TodoHandler with mock storage", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStorage := mock.NewMockTodoStorage(ctrl)
		e := gin.Default()
		RegisterTodoHandler(e, mockStorage)

		Convey("When deleting a todo with valid ID", func() {
			mockStorage.EXPECT().
				Delete(gomock.Any(), gomock.Eq(1)).
				Return(nil).
				Times(1)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodDelete, "/todos/1", nil)
			e.ServeHTTP(w, req)

			Convey("Then it should return 204 status code", func() {
				So(w.Code, ShouldEqual, http.StatusNoContent)
			})
		})

		Convey("When deleting a todo with invalid ID format", func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodDelete, "/todos/invalid", nil)
			e.ServeHTTP(w, req)

			Convey("Then it should return 400 status code", func() {
				So(w.Code, ShouldEqual, http.StatusBadRequest)
			})
		})

		Convey("When todo is not found", func() {
			mockStorage.EXPECT().
				Delete(gomock.Any(), gomock.Eq(999)).
				Return(gorm.ErrRecordNotFound).
				Times(1)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodDelete, "/todos/999", nil)
			e.ServeHTTP(w, req)

			Convey("Then it should return 404 status code", func() {
				So(w.Code, ShouldEqual, http.StatusNotFound)
			})
		})

		Convey("When storage returns an error", func() {
			mockStorage.EXPECT().
				Delete(gomock.Any(), gomock.Any()).
				Return(errors.New("internal server error")).
				Times(1)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodDelete, "/todos/1", nil)
			e.ServeHTTP(w, req)

			Convey("Then it should return 500 status code", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})
}
