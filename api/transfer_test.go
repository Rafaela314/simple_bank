package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	mockdb "simple_bank/db/mock"
	db "simple_bank/db/sqlc"
	"simple_bank/util"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateTransferAPI(t *testing.T) {
	account1 := createRandomAccount()
	account2 := createRandomAccount()
	// Use a supported currency (USD/EUR); CAD is not accepted by validCurrency.
	account1.Currency = util.USD
	account2.Currency = util.USD
	// Amount must be >= 1 (binding "min=1"); RandomMoney() can return 0.
	amount := util.RandomInt(1, 100)

	testCases := []struct {
		name          string
		body          createTransferRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: createTransferRequest{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
				Currency:      account1.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, nil)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).
					Return(account2, nil)
				transfer := db.Transfer{
					ID:            util.RandomInt(1, 1000),
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store.EXPECT().
					CreateTransfer(gomock.Any(), gomock.Eq(db.CreateTransferParams{
						FromAccountID: account1.ID,
						ToAccountID:   account2.ID,
						Amount:        amount,
					})).
					Times(1).
					Return(transfer, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTransfer(t, recorder.Body, db.Transfer{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				})
			},
		},
		{
			name: "FromAccountNotFound",
			body: createTransferRequest{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
				Currency:      account1.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(0)
				store.EXPECT().
					CreateTransfer(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				require.Contains(t, recorder.Body.String(), "error")
			},
		},
		{
			name: "ToAccountNotFound",
			body: createTransferRequest{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
				Currency:      account1.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, nil)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
				store.EXPECT().
					CreateTransfer(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				require.Contains(t, recorder.Body.String(), "error")
			},
		},
		{
			name: "InvalidRequest",
			body: createTransferRequest{
				FromAccountID: 0,
				ToAccountID:   account2.ID,
				Amount:        amount,
				Currency:      account1.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					CreateTransfer(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				require.Contains(t, recorder.Body.String(), "error")
			},
		},
		{
			name: "InternalErrorOnCreateTransfer",
			body: createTransferRequest{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
				Currency:      account1.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, nil)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).
					Return(account2, nil)
				store.EXPECT().
					CreateTransfer(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Transfer{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
				require.Contains(t, recorder.Body.String(), "error")
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/transfers", bytes.NewReader(data))
			require.NoError(t, err)
			request.Header.Set("Content-Type", "application/json")

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func requireBodyMatchTransfer(t *testing.T, body *bytes.Buffer, expected db.Transfer) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	var gotTransfer db.Transfer
	err = json.Unmarshal(data, &gotTransfer)
	require.NoError(t, err)
	require.Equal(t, expected.FromAccountID, gotTransfer.FromAccountID)
	require.Equal(t, expected.ToAccountID, gotTransfer.ToAccountID)
	require.Equal(t, expected.Amount, gotTransfer.Amount)
	if expected.ID != 0 {
		require.Equal(t, expected.ID, gotTransfer.ID)
	}
}
