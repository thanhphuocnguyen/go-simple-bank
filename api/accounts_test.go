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

	"github.com/golang/mock/gomock"
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
