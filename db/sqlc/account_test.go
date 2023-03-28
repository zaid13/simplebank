package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"github.com/zaid13/simplebank/db/util"
	"testing"
	"time"
)

func createRandomAccount(t *testing.T) Account{
	arg:=CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account,err:=testQueries.CreateAccount(context.Background(),arg)
	require.NoError(t,err )
	require.NotEmpty(t,account )

	require.Equal(t,arg.Currency , account.Currency )
	require.Equal(t,arg.Balance , account.Balance )
	require.Equal(t,arg.Owner , account.Owner )

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt )


return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T )  {
	//crete account
	account1:=createRandomAccount(t)
	account2,err:=testQueries.GetAccount(context.Background(),account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)


	require.Equal(t, account1.ID,account2.ID)
	require.Equal(t, account1.Currency,account2.Currency)
	require.Equal(t, account1.Balance,account2.Balance)
	require.Equal(t, account1.Owner,account2.Owner)

	require.WithinDuration(t, account1.CreatedAt,account2.CreatedAt,time.Second)





}

func TestUpdateAccount(t *testing.T )  {
	//update account
	account1:=createRandomAccount(t)

	updateArgs:=UpdateAccountParams{
		account1.ID,
		util.RandomMoney(),
	}

	account2,err:=testQueries.UpdateAccount(
		context.Background(),
		updateArgs,
	)

	require.NoError(t, err)
	require.NotEmpty(t, account2)


	require.Equal(t, account1.ID,account2.ID)
	require.Equal(t, account1.Currency,account2.Currency)
	require.Equal(t, updateArgs.Balance,account2.Balance)
	require.Equal(t, account1.Owner,account2.Owner)

	require.WithinDuration(t, account1.CreatedAt,account2.CreatedAt,time.Second)





}

func  TestDeleteAccount(t *testing.T)  {
	//delete account

	account1:=createRandomAccount(t)

	err:=testQueries.DeleteAccount(
		context.Background(),
		account1.ID,
	)

	account2,err:=testQueries.GetAccount(context.Background() ,account1.ID )

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)



}

func  TestListAccounts(t *testing.T)  {
	//delete account


	for i:=0 ; i<10 ; i++{
		createRandomAccount(t)
	}


	accounts, err:=testQueries.ListAccount(context.Background(),ListAccountParams{5,5})


	print(len(accounts))
	require.NoError(t, err)
	require.Equal(t, len(accounts),5)

	for _,accounts:=range accounts{
		require.NotEmpty(t,accounts )
	}

}