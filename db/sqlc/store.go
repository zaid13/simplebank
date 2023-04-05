package db

import (
	"context"
	"database/sql"
	"fmt"
)

//store provides all function to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

//NewStore creates a new store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

//execTx executes creates a function inside a database transaction
func (Store *Store) execTx(ctx context.Context, fn func(queries *Queries) error) error {

	tx, err := Store.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		rberr := tx.Rollback()
		if rberr != nil {
			return fmt.Errorf("tx err:%v , rb err:%v", err, rberr)
		}

	}
	return tx.Commit()
}

//TransferTxParams contain the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

//TransferTxResult contain the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

//TransferTx  performs a money transfer from one account to other.
//It creates a transfer record , add account enteries , and update accounts' balance within a single database transaction

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {

	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{arg.FromAccountID, arg.ToAccountID, arg.Amount})
		if err != nil {
			return nil
		}

		result.FromEntry, err = q.CreateEntries(ctx, CreateEntriesParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return nil
		}

		result.ToEntry, err = q.CreateEntries(ctx, CreateEntriesParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return nil
		}

		if arg.FromAccountID > arg.ToAccountID {

			result.FromAccount, result.ToAccount, err = AddMoney(ctx, q, arg.FromAccountID, arg.ToAccountID, -arg.Amount, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = AddMoney(ctx, q, arg.ToAccountID, arg.FromAccountID, arg.Amount, -arg.Amount)

		}
		//result.FromAccount, err =q.AddAccountBalance(ctx,AddAccountBalanceParams{
		//	ID:     arg.FromAccountID,
		//	Amount: -arg.Amount,
		//})
		//if err != nil {
		//	return err
		//}
		//
		//result.ToAccount, err =q.AddAccountBalance(ctx,AddAccountBalanceParams{
		//	ID:     arg.ToAccountID,
		//	Amount: arg.Amount,
		//})
		//if err != nil {
		//	return err
		//}

		return nil
	})
	return result, err
}

func AddMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	accountID2 int64,
	amount1 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {

	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	if err != nil {
		return
	}

	return

}
