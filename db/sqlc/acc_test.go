package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"
)

func creatRadomAccount(t *testing.T) Account {

	user := createRadomUser(t)
	params := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	acc, err := testStore.CreateAccount(context.Background(), params)

	if err != nil {
		fmt.Println('a')

	}
	require.Equal(t, params.Owner, acc.Owner)
	require.Equal(t, params.Balance, acc.Balance)
	require.Equal(t, params.Currency, acc.Currency)
	require.NotZero(t, acc.ID)
	require.NotZero(t, acc.CreatedAt)
	return acc
}

func TestCreateAccount(t *testing.T) {

	creatRadomAccount(t)

}

func TestGetAccount(t *testing.T) {
	acc1 := creatRadomAccount(t)
	acc2, err := testStore.GetAccount(context.Background(), acc1.ID)

	// If there is an error, the test will fail here.
	require.NoError(t, err)
	//This asserts that the retrieved account (acc2) is not empty
	require.NotEmpty(t, acc2)
	// This asserts that the owner of the created account (acc1.Owner) matches the owner of the retrieved account (acc2.Owner)
	require.Equal(t, acc1.Owner, acc2.Owner)
	require.Equal(t, acc1.Balance, acc2.Balance)
	require.Equal(t, acc1.Currency, acc2.Currency)
	//This accounts for any minor delays in database transactions.
	require.WithinDuration(t, acc1.CreatedAt, acc2.CreatedAt, time.Second)

}

func TestUpdateAccount(t *testing.T) {
	acc1 := creatRadomAccount(t)
	arg := UpdateAccountParams{
		ID:      acc1.ID,
		Balance: util.RandomMoney(),
	}

	acc2, err := testStore.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, acc2)
	require.Equal(t, acc1.Owner, acc2.Owner)
	require.Equal(t, arg.Balance, acc2.Balance)
	require.Equal(t, acc1.Currency, acc2.Currency)
	require.WithinDuration(t, acc1.CreatedAt, acc2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	acc1 := creatRadomAccount(t)
	err := testStore.DeleteAccount(context.Background(), acc1.ID)
	require.NoError(t, err)
	acc2, err := testStore.GetAccount(context.Background(), acc1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, acc2)
}

func TestGetListAccounts(t *testing.T) {

	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = creatRadomAccount(t)
	}

	arg := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accs, err := testStore.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accs)
	for _, acc := range accs {
		require.NotEmpty(t, acc)
		require.Equal(t, lastAccount.Owner, acc.Owner)
	}

}
