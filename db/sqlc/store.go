package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

type SQLStore struct {
	*Queries         // Embeds Queries struct to inherit its methods
	db       *sql.DB // Holds a reference to the database connection
}

// It initializes a new Store instance with a given database connection.
func NewsStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
	//Returns: A pointer to a new Store instance.
}

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	//Creates a new Queries instance using the transaction
	q := New(tx)
	//Excute queries
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountId int64 `json:"from_account_id"`
	ToAccountId   int64 `json:"to_account_id"`
	Amount        int64 `json: "amount`
}

type TransferTxResult struct {
	Transfer    Transfer `json: "transfer"`
	FromAccount Account  `json: "from_account"`
	ToAccount   Account  `json:"to_account`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `jsonL"to_entry`
}

// Transfer performs a maonry transfer from one account to the other.
// It creates a transfer record , add accs entries, and update accounts' balance within a single transaction
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID:   arg.ToAccountId,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		//Add account entries
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountId,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountId,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		//Update accounts' balance

		if arg.FromAccountId < arg.ToAccountId {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountId, -arg.Amount, arg.ToAccountId, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountId, arg.Amount, arg.FromAccountId, -arg.Amount)

		}

		return nil
	})
	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddUpdateAccount(ctx, AddUpdateAccountParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}
	account2, err = q.AddUpdateAccount(ctx, AddUpdateAccountParams{
		ID:     accountID2,
		Amount: amount2,
	})

	return

}
