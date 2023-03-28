package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferTxDeadLock(t *testing.T) {

	store :=NewStore(testDB)
	account1:=createRandomAccount(t)
	account2:=createRandomAccount(t)

	fmt.Println(">>balance ", account1.Balance , account2.Balance)


	//run n concurrent transfer transactions

	n:=10

	amount:=int64(10)

	errs:=make(chan error)



	for i:=0 ; i<n ; i++{
		fromAccountID:=account1.ID
		toAccountID:=account2.ID
		//txName:=fmt.Sprintf("tx %d",i+1)

		if i%2 ==1{
			fromAccountID=account2.ID
			toAccountID=account1.ID

		}
		go func() {
			ctx := context.Background()//context.WithValue(context.Background(),txKey,txName)

			_,err:=store.TransferTx(ctx , TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:  toAccountID,
				Amount:        amount,
			})
			errs <- err

		}()
	}

	 //	check results


	for i:=0 ; i<n ; i++ {
		err := <-errs

		require.NoError(t, err)

	}

		//check the final updated balance
		updatedAccount1,err:=testQueries.GetAccount(context.Background() , account1.ID)
		require.NoError(t, err)

		updatedAccount2,err:=testQueries.GetAccount(context.Background() , account2.ID)
		require.NoError(t, err)

		fmt.Println(">>updated balance ", updatedAccount1.Balance , updatedAccount2.Balance)

		require.Equal(t, account1.Balance  , updatedAccount1.Balance)
		require.Equal(t, account2.Balance  , updatedAccount2.Balance)


		// todo check account balance





}
