package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, acc1 Account, acc2 Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Amount:        10,
	}

	tranfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, tranfer)
	require.Equal(t, tranfer.FromAccountID, arg.FromAccountID)
	require.Equal(t, tranfer.ToAccountID, arg.ToAccountID)
	require.Equal(t, tranfer.Amount, arg.Amount)
	require.NotZero(t, tranfer.ID)
	require.NotZero(t, tranfer.CreatedAt)

	return tranfer

}

func TestCreatTransfer(t *testing.T) {
	acc1 := creatRadomAccount(t)
	acc2 := creatRadomAccount(t)
	createRandomTransfer(t, acc1, acc2)
}

func TestGetTransfer(t *testing.T) {
	acc1 := creatRadomAccount(t)
	acc2 := creatRadomAccount(t)
	transfer1 := createRandomTransfer(t, acc1, acc2)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)

	require.NoError(t, err)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, transfer1.CreatedAt, time.Second)

}

func TestListTrasnfers(t *testing.T) {
	acc1 := creatRadomAccount(t)
	acc2 := creatRadomAccount(t)
	for i := 0; i < 5; i++ {
		createRandomTransfer(t, acc1, acc2)
		createRandomTransfer(t, acc2, acc1)
	}

	arg := ListTransfersParams{
		FromAccountID: acc1.ID,
		ToAccountID:   acc1.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)

	require.Len(t, transfers, 5)

	for _, tr := range transfers {
		require.NotEmpty(t, tr)
		require.True(t, tr.FromAccountID == acc1.ID || tr.ToAccountID == acc1.ID)
	}

}
