package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)
	fmt.Println(">> before:", acc1.Balance, acc2.Balance)

	n := 3
	amount := int64(10)
	errsChan := make(chan error)
	resultsChan := make(chan TransferTxResult)
	doneTx := make(map[int]bool)

	// run n concurrent transfer transactions
	for i := 0; i < n; i++ {
		go func() {
			txName := fmt.Sprintf("tx %d", i+1)
			ctx := context.WithValue(context.Background(), txKey, txName)

			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: acc1.ID,
				ToAccountID:   acc2.ID,
				Amount:        amount,
			})
			errsChan <- err
			resultsChan <- result
		}()
	}
	//check results
	for i := 0; i < n; i++ {
		err := <-errsChan
		require.NoError(t, err)

		result := <-resultsChan
		require.NotEmpty(t, result)

		// check transfers
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, acc1.ID, transfer.FromAccountID)
		require.Equal(t, acc2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		//try to get the records
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, acc1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, acc2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)
		//*
		//check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, acc1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, acc2.ID, toAccount.ID)

		// check balances
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)

		diff1 := acc1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - acc2.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		txkey := int(diff1 / amount)
		require.True(t, txkey >= 1 && txkey <= n)
		require.NotContains(t, doneTx, txkey)
		doneTx[txkey] = true
	}
	updatedAcc1, err := store.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	updatedAcc2, err := store.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAcc1.Balance, updatedAcc2.Balance)
	require.Equal(t, acc1.Balance-int64(n)*amount, updatedAcc1.Balance)
	require.Equal(t, acc2.Balance+int64(n)*amount, updatedAcc2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)
	fmt.Println(">> before:", acc1.Balance, acc2.Balance)

	n := 10
	amount := int64(10)
	errsChan := make(chan error)

	// run n concurrent transfer transactions
	for i := 0; i < n; i++ {
		fromAccID := acc1.ID
		toAccID := acc2.ID

		//split tx
		if i%2 == 1 {
			fromAccID = acc2.ID
			toAccID = acc1.ID
		}

		go func() {
			txName := fmt.Sprintf("tx %d", i+1)
			ctx := context.WithValue(context.Background(), txKey, txName)

			_, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccID,
				ToAccountID:   toAccID,
				Amount:        amount,
			})
			errsChan <- err
		}()
	}
	//check results
	for i := 0; i < n; i++ {
		err := <-errsChan
		require.NoError(t, err)
	}
	updatedAcc1, err := store.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	updatedAcc2, err := store.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAcc1.Balance, updatedAcc2.Balance)
	//balance should be the same here: GIGO
	require.Equal(t, acc1.Balance, updatedAcc1.Balance)
	require.Equal(t, acc2.Balance, updatedAcc2.Balance)
}
