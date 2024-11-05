package http

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"bookloop.net/config"
	ucMock "bookloop.net/internal/books/mock"
	"bookloop.net/internal/models"
	"bookloop.net/pkg/utils"
	"bookloop.net/pkg/validator"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestBooksHandlers_List(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	uc := ucMock.NewMockUseCase(ctl)

	cfg := config.LoadConfig()
	sl := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)

	h := booksHandlers{
		cfg:     cfg,
		booksUC: uc,
		logger:  sl,
	}

	type testCase struct {
		name           string
		queryParams    string
		mockBooks      []*models.Book
		mockPagination utils.Pagination
		mockError      error
		expectedStatus int
		expectedBody   string
	}

	testCases := []testCase{
		{
			name:           "valid request",
			queryParams:    "page=1&page_size=2&sort=-created_at",
			mockBooks:      []*models.Book{{ID: 1, Title: "Test Book 1", Author: "Author1", Version: 1}, {ID: 2, Title: "Test Book 2", Author: "Author2", Version: 2}},
			mockPagination: utils.Pagination{CurrentPage: 1, PageSize: 2, FirstPage: 1, LastPage: 1, TotalRecords: 2},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"books":[{"id":1,"title":"Test Book 1","author":"Author1","version":1},{"id":2,"title":"Test Book 2","author":"Author2","version":2}],"metadata":{"current_page":1,"page_size":2,"total_records":2,"first_page":1,"last_page":1}}`,
		},
		{
			name:           "invalid page size",
			queryParams:    "page_size=150",
			mockBooks:      nil,
			mockPagination: utils.Pagination{},
			mockError: &validator.ValidationError{
				Errors: map[string]string{"page_size": "must be a maximum of 100"},
				Err:    validator.ErrJSONIsNotValid,
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"error":{"page_size":"must be a maximum of 100"}}`,
		},
		{
			name:           "invalid sort parameter",
			queryParams:    "sort=-invalid_field",
			mockBooks:      nil,
			mockPagination: utils.Pagination{},
			mockError: &validator.ValidationError{
				Errors: map[string]string{"sort": "invalid sort value"},
				Err:    validator.ErrJSONIsNotValid,
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"error":{"sort":"invalid sort value"}}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uc.
				EXPECT().
				List(gomock.Any(), gomock.Any()).
				Return(tc.mockBooks, tc.mockPagination, tc.mockError)

			req, err := http.NewRequest(http.MethodGet, "/books?"+tc.queryParams, nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()

			handler := h.List()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
			assert.JSONEq(t, tc.expectedBody, rr.Body.String())
		})
	}
}
