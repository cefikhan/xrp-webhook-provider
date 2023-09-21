package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type LedgerCurrentBody struct {
	Method string `json:"method"`
}

type LedgerCurrentResultStruct struct {
	LedgerCurrentIndex int    `json:"ledger_current_index"`
	Status             string `json:"status"`
}

type LedgerResponse struct {
	Result LedgerCurrentResultStruct `json:"result"`
}

func XRPLedgerService() {

	//Get current block number from database

	currentLedger, err := getLatestLedgerNofromDb(db)
	if err != nil {
		fmt.Println("Error Fetchin latest ledger no from db")
	}
	if currentLedger == 0 {
		//if blocknumber from db is less than currentblock
		//fetch it from Blockchain Node
		currentLedger, err = getLatestLedgerNofromNode()
		if err != nil {
			fmt.Println("Error")
		}

	}
	currentLedgerNode, err := getLatestLedgerNofromNode()
	if err != nil {
		fmt.Println("Error")
	}

	if currentLedger > currentLedgerNode {
		fmt.Println("Syncing the node sleeping for 5s")
		time.Sleep(5 * time.Second)

	}

	//passing current block number Fetch blockdetails
	getLedgerInfo(currentLedger)
	currentLedger = currentLedger + 1
	UpdateLedgerNoInDb(db, currentLedger)
	beeper.Wait()

}

func getLatestLedgerNofromNode() (int, error) {

	requestMap := &LedgerCurrentBody{

		Method: "ledger_current",
	}
	requestBytes, err := json.Marshal(requestMap)
	if err != nil {
		return 0, err
	}

	transaction, err := http.Post("https://s1.ripple.com:51234/", "application/json", bytes.NewReader(requestBytes))
	if err != nil {
		fmt.Println("Error making request")
	}
	defer transaction.Body.Close()
	ledgerResponse := LedgerResponse{}

	err = json.NewDecoder(transaction.Body).Decode(&ledgerResponse)
	if err != nil {
		return 0, err
	}

	return ledgerResponse.Result.LedgerCurrentIndex, nil

}
