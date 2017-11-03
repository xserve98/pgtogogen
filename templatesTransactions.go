package main

const BASE_TRANSACTIONS = `package {{.PackageName}}

/* ************************************************************* */
/* This file was automatically generated by pgtogogen.           */
/* Do not modify this file unless you know what you are doing.   */
/* ************************************************************* */

import (
	"context"
	pgx "{{.PgxImport}}"	
)

//
// DB transaction-related types and functionality
//

// Transaction isolation levels for the pgx package

const IsoLevelSerializable = pgx.Serializable
const IsoLevelRepeatableRead = pgx.RepeatableRead
const IsoLevelReadCommitted = pgx.ReadCommitted
const IsoLevelReadUncommitted = pgx.ReadUncommitted

// Wrapper structure over the pgx transaction package, so we don't need to import
// that package in the generated table-to-struct files.
type Transaction struct {
	Tx *pgx.Tx
}

// Commits the current transaction
func (t *Transaction) Commit() error {
	if t.Tx == nil {
		return NewModelsErrorLocal("Transaction.Commit()", "The inner Tx transaction is nil")
	}
	return t.Tx.Commit()
}

// Attempts to rollback the current transaction
func (t *Transaction) Rollback() error {
	if t.Tx == nil {
		return NewModelsErrorLocal("Transaction.Rollback()", "The inner Tx transaction is nil")
	}
	return t.Tx.Rollback()
}

/* BEGIN Transactions utility functions */

// Begins and returns a transaction using the default isolation level.
// Unlike TxWrap, it is the responsibility of the caller to commit and
// rollback the transaction if necessary.
func TxBegin() (*Transaction, error) {

	txWrapper := &Transaction{}
	tx, err := GetDb().Begin()

	if err != nil {
		return nil, err
	} else {
		txWrapper.Tx = tx
		return txWrapper, nil
	}

}

// Begins and returns a transaction using the specified isolation level.
// The following global constants can be passed (residing in the same package):
//  IsoLevelSerializable
//  IsoLevelRepeatableRead
//  IsoLevelReadCommitted
//  IsoLevelReadUncommitted
func TxBeginIso(isolationLevel pgx.TxIsoLevel) (*Transaction, error) {

	txWrapper := &Transaction{}
	tx, err := GetDb().BeginEx(context.Background(), &pgx.TxOptions{IsoLevel: isolationLevel})

	if err != nil {
		return nil, err
	} else {
		txWrapper.Tx = tx
		return txWrapper, nil
	}
}

/* This method helps wrap the transaction inside a closure function. Additional arguments can be passed
 along to the closure via a variadic list of interface{} parameters. 
 TxWrap automatically handles commit and rollback, in case of error. 
 It returns an error in case of failure, or nil, in case of success.

 Example:

	// define the transaction functionlity in this wrapper closure
	var transactionFunc = func(tx *models.Transaction, arguments ...interface{}) (interface{}, error) {

		// assuming the generated package is named models and
		// there is a TestEvent struct corresponding to a test_event table in the database
		newTestEvent := models.Tables.TestEvent.New()

		// load the event name as passed via the variadic arguments
		newTestEvent.SetEventName(arguments[0].(string))
		newTestEvent.SetEventOverview(arguments[1].(string), true)

		newTestEvent, err := tx.InsertTestEvent(newTestEvent)
		if err != nil {
			return nil, models.NewModelsError("insert event tx error:", err)
		}

		// any other transaction operations...

		// at the end, we return nil for a successful operation
		return newTestEvent, nil
	}

	// define some parameters to be passed inside the transaction
	eventName := "Donald Duck Anniversary"
	eventDescription := "Where is the party ?"

	// we defined the transaction functionality, let's run it with the event name argument
	returnedNewEvent, err := models.TxWrap(transactionFunc, eventName, eventDescription)
	if err != nil {
		fmt.Println("FAIL:", err.Error())
	} else {
		if returnedNewEvent == nil {
			fmt.Printf("OK. But newlyInsertedEvent is nil \r\n")
		} else {
			// we need to make sure to convert the resulting type to the needs of this particular transaction
			fmt.Printf("OK. newlyInsertedEvent overview: " + returnedNewEvent.(*models.TestEvent).EventOverview + "  \r\n")
		}
	} */
func TxWrap(wrapperFunc func(tx *Transaction, args ...interface{}) (interface{}, error), arguments ...interface{}) (interface{}, error) {

	var errorPrefix = "TxWrap() ERROR: "

	realTx, err := GetDb().Begin()
	if err != nil {
		return nil, NewModelsError(errorPrefix+"GetDb().Begin() error: ", err)
	}

	// pgx package note: Rollback is safe to call even if the tx is already closed,
	// so if the tx commits successfully, this is a no-op
	defer realTx.Rollback()

	// wrap the real tx into our wrapper
	tx := &Transaction{Tx: realTx}

	result, err := wrapperFunc(tx, arguments...)
	if err != nil {
		return nil, NewModelsError(errorPrefix+"inner wrapperFunc() error - will return and rollback: ", err)
	}

	err = realTx.Commit()
	if err != nil {
		return nil, NewModelsError(errorPrefix+"tx.Commit() error: ", err)
	}

	return result, nil
}

/* END Transactions utility functions */

`
