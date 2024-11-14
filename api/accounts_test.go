package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
	"github.com/thanhphuocnguyen/go-simple-bank/auth"
	mockdb "github.com/thanhphuocnguyen/go-simple-bank/db/mock"
	db "github.com/thanhphuocnguyen/go-simple-bank/db/sqlc"
	utils "github.com/thanhphuocnguyen/go-simple-bank/utils"
)

func TestGetAccountAPI(t *testing.T) {
	user, _, _ := randomUser(t)
	account := randomAccount(user.Username)
	testCases := []struct {
		name            string
		accountID       int64
		setupAuthHeader func(t *testing.T, request *http.Request, tokenGenerator auth.TokenGenerator)
		buildStub       func(store *mockdb.MockStore)
		check           func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			setupAuthHeader: func(t *testing.T, request *http.Request, tokenGenerator auth.TokenGenerator) {
				addAuthHeader(t, request, tokenGenerator, authorizationType, account.Owner, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "Wrong User",
			accountID: account.ID,
			setupAuthHeader: func(t *testing.T, request *http.Request, tokenGenerator auth.TokenGenerator) {
				addAuthHeader(t, request, tokenGenerator, authorizationType, "Hello", time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "Unauthorized",
			accountID: account.ID,
			setupAuthHeader: func(t *testing.T, request *http.Request, tokenGenerator auth.TokenGenerator) {
				// Don't add auth header
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			setupAuthHeader: func(t *testing.T, request *http.Request, tokenGenerator auth.TokenGenerator) {
				addAuthHeader(t, request, tokenGenerator, authorizationType, account.Owner, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalServerError",
			accountID: account.ID,
			setupAuthHeader: func(t *testing.T, request *http.Request, tokenGenerator auth.TokenGenerator) {
				addAuthHeader(t, request, tokenGenerator, authorizationType, account.Owner, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "InternalServerError",
			accountID: 0,
			setupAuthHeader: func(t *testing.T, request *http.Request, tokenGenerator auth.TokenGenerator) {
				addAuthHeader(t, request, tokenGenerator, authorizationType, account.Owner, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStub(store)

			server := createNewServer(t, store)

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuthHeader(t, request, server.tokenGenerator)
			server.router.ServeHTTP(recorder, request)

			tc.check(t, recorder)
		})

	}
}

func TestCreateAccountAPI(t *testing.T) {
	user, _, _ := randomUser(t)
	account := randomAccount(user.Username)
	testCases := []struct {
		name            string
		body            gin.H
		setupAuthHeader func(t *testing.T, request *http.Request, tokenGenerator auth.TokenGenerator)
		buildStub       func(store *mockdb.MockStore)
		check           func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"currency": account.Currency,
			},
			setupAuthHeader: func(t *testing.T, request *http.Request, tokenGenerator auth.TokenGenerator) {
				addAuthHeader(t, request, tokenGenerator, authorizationType, account.Owner, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Owner:    account.Owner,
					Balance:  0,
					Currency: account.Currency,
				}
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(account, nil)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name: "BadRequest",
			body: gin.H{
				"currency": "invalid",
			},
			setupAuthHeader: func(t *testing.T, request *http.Request, tokenGenerator auth.TokenGenerator) {
				addAuthHeader(t, request, tokenGenerator, authorizationType, account.Owner, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(0)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "UniqueViolation",
			body: gin.H{
				"currency": account.Currency,
			},
			setupAuthHeader: func(t *testing.T, request *http.Request, tokenGenerator auth.TokenGenerator) {
				addAuthHeader(t, request, tokenGenerator, authorizationType, account.Owner, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Owner:    account.Owner,
					Balance:  0,
					Currency: account.Currency,
				}
				pgErr := &pgconn.PgError{Code: "23505"}
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Account{}, pgErr)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "ForeignKeyViolation",
			body: gin.H{
				"currency": account.Currency,
			},
			buildStub: func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Owner:    account.Owner,
					Balance:  0,
					Currency: account.Currency,
				}
				pgErr := &pgconn.PgError{Code: "23503"}
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Account{}, pgErr)
			},
			setupAuthHeader: func(t *testing.T, request *http.Request, tokenGenerator auth.TokenGenerator) {
				addAuthHeader(t, request, tokenGenerator, authorizationType, account.Owner, time.Minute)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			body: gin.H{
				"currency": account.Currency,
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			setupAuthHeader: func(t *testing.T, request *http.Request, tokenGenerator auth.TokenGenerator) {
				addAuthHeader(t, request, tokenGenerator, authorizationType, account.Owner, time.Minute)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Unauthorized",
			body: gin.H{
				"currency": account.Currency,
			},
			setupAuthHeader: func(t *testing.T, request *http.Request, tokenGenerator auth.TokenGenerator) {
				// Don't add auth header
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStub(store)

			server := createNewServer(t, store)

			recorder := httptest.NewRecorder()
			body, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(body))
			tc.setupAuthHeader(t, request, server.tokenGenerator)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.check(t, recorder)
		})
	}
}

func TestGetAccountsAPI(t *testing.T) {
	type Query struct {
		Page     int32 `form:"page"`
		PageSize int32 `form:"page_size"`
	}
	user, _, _ := randomUser(t)
	n := 5
	accounts := make([]db.Account, n)
	for i := 0; i < n; i++ {
		accounts[i] = randomAccount(user.Username)
	}

	testCases := []struct {
		name            string
		queries         Query
		setupAuthHeader func(t *testing.T, request *http.Request, tokenGenerator auth.TokenGenerator)
		buildStub       func(store *mockdb.MockStore)
		check           func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			queries: Query{
				Page:     1,
				PageSize: 7,
			},
			setupAuthHeader: func(t *testing.T, request *http.Request, tokenGenerator auth.TokenGenerator) {
				addAuthHeader(t, request, tokenGenerator, authorizationType, user.Username, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				arg := db.ListAccountsParams{
					Owner:  user.Username,
					Limit:  7,
					Offset: 0,
				}
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(accounts, nil)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccounts(t, recorder.Body, accounts)
			},
		},
		{
			name: "PageBadRequest",
			queries: Query{
				Page:     0,
				PageSize: 7,
			},
			setupAuthHeader: func(t *testing.T, request *http.Request, tokenGenerator auth.TokenGenerator) {
				addAuthHeader(t, request, tokenGenerator, authorizationType, user.Username, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "PageSizeBadRequest",
			queries: Query{
				Page:     1,
				PageSize: 4,
			},
			setupAuthHeader: func(t *testing.T, request *http.Request, tokenGenerator auth.TokenGenerator) {
				addAuthHeader(t, request, tokenGenerator, authorizationType, user.Username, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Unauthorized",
			queries: Query{
				Page:     1,
				PageSize: 7,
			},
			setupAuthHeader: func(t *testing.T, request *http.Request, tokenGenerator auth.TokenGenerator) {
				addAuthHeader(t, request, tokenGenerator, "Custom", user.Username, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			queries: Query{
				Page:     1,
				PageSize: 7,
			},
			buildStub: func(store *mockdb.MockStore) {
				arg := db.ListAccountsParams{
					Owner:  user.Username,
					Limit:  7,
					Offset: 0,
				}
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return([]db.Account{}, sql.ErrConnDone)
			},
			setupAuthHeader: func(t *testing.T, request *http.Request, tokenGenerator auth.TokenGenerator) {
				addAuthHeader(t, request, tokenGenerator, authorizationType, user.Username, time.Minute)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStub(store)

			server := createNewServer(t, store)

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/accounts?page=%d&page_size=%d", tc.queries.Page, tc.queries.PageSize)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuthHeader(t, request, server.tokenGenerator)
			server.router.ServeHTTP(recorder, request)

			tc.check(t, recorder)
		})
	}
}

func randomAccount(username string) db.Account {
	return db.Account{
		ID:       utils.RandomInt(1, 1000),
		Owner:    username,
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}

func requireBodyMatchAccounts(t *testing.T, body *bytes.Buffer, accounts []db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccounts []db.Account
	err = json.Unmarshal(data, &gotAccounts)
	require.NoError(t, err)
	require.Equal(t, accounts, gotAccounts)
}
