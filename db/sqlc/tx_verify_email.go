package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type VerrifyEmailTxParams struct {
	EmailId    int64
	SecretCode string
}

type VerrifyEmailTxResult struct {
	User        User
	VerifyEmail VerifyEmail
}

// Transfer performs a maonry transfer from one account to the other.
// It creates a transfer record , add accs entries, and update accounts' balance within a single transaction
func (store *SQLStore) VerrifyEmailTx(ctx context.Context, arg VerrifyEmailTxParams) (VerrifyEmailTxResult, error) {
	var result VerrifyEmailTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.VerifyEmail, err = q.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID:         arg.EmailId,
			SecretCode: arg.SecretCode,
		})
		if err != nil {
			return err

		}
		result.User, err = q.UpdateUser(ctx, UpdateUserParams{
			Username: result.User.Username,
			IsEmailVerified: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
		})
		return err

	})
	return result, err
}
