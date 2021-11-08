package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	queriesObj := New(tx)
	err = fn(queriesObj)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

//Transfer transaction: create a new transfer record, add 2 new account entries, and update the 2 accounts’ balance within a single database transaction.
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

func (s *Store) TransferTx(ctx context.Context, args TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := s.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(txKey)

		fmt.Println(txName, "create transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: args.FromAccountID,
			ToAccountID:   args.ToAccountID,
			Amount:        args.Amount,
		})

		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.FromAccountID,
			Amount:    -args.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.ToAccountID,
			Amount:    args.Amount,
		})
		if err != nil {
			return err
		}

		//Add update accounts' balance later
		// fmt.Println(txName, "get account 1")
		// acc1, err := q.GetAccountForUpdate(ctx, args.FromAccountID)
		// if err != nil {
		// 	return err
		// }

		//move money out of acc1: sender
		// fmt.Println(txName, "update account 1")
		// result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
		// 	ID:      args.FromAccountID,
		// 	Balance: acc1.Balance - args.Amount,
		// })

		// always update the account with smaller ID first to avoid deadlock
		if args.FromAccountID < args.ToAccountID {

			result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     args.FromAccountID,
				Amount: -args.Amount,
			})
			if err != nil {
				return err
			}
			result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     args.ToAccountID,
				Amount: args.Amount,
			})
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     args.ToAccountID,
				Amount: args.Amount,
			})
			if err != nil {
				return err
			}

			result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     args.FromAccountID,
				Amount: -args.Amount,
			})
			if err != nil {
				return err
			}
		}

		// fmt.Println(txName, "get account 2")
		// acc2, err := q.GetAccountForUpdate(ctx, args.ToAccountID)
		// if err != nil {
		// 	return err
		// }

		//move money into acc2: receiver

		// fmt.Println(txName, "update account 2")
		// result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
		// 	ID:      args.ToAccountID,
		// 	Balance: acc2.Balance + args.Amount,
		// })

		return err
	})
	return result, err
}
