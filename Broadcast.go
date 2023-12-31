package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AutoGeneratedX struct {
	Result struct {
		Error string `json:"error"`

		Ledger struct {
			Accepted            bool   `json:"accepted"`
			AccountHash         string `json:"account_hash"`
			CloseFlags          int    `json:"close_flags"`
			CloseTime           int    `json:"close_time"`
			CloseTimeHuman      string `json:"close_time_human"`
			CloseTimeResolution int    `json:"close_time_resolution"`
			Closed              bool   `json:"closed"`
			Hash                string `json:"hash"`
			LedgerHash          string `json:"ledger_hash"`
			LedgerIndex         string `json:"ledger_index"`
			ParentCloseTime     int    `json:"parent_close_time"`
			ParentHash          string `json:"parent_hash"`
			SeqNum              string `json:"seqNum"`
			TotalCoins          string `json:"totalCoins"`
			Total_Coins         string `json:"total_coins"`
			TransactionHash     string `json:"transaction_hash"`
			Transactions        []struct {
				Fee             string      `json:"Fee"`
				Account         string      `json:"Account"`
				Amount          interface{} `json:"Amount"`
				Destination     string      `json:"Destination"`
				TransactionType string      `json:"TransactionType"`

				Hash string `json:"hash"`
			}
		} `json:"ledger"`
	} `json:"result"`
}

type TransactionRes struct {
	Fee             string      `json:"Fee"`
	Account         string      `json:"Account"`
	Amount          interface{} `json:"Amount"`
	Destination     string      `json:"Destination"`
	TransactionType string      `json:"TransactionType"`
	Hash            string      `json:"hash"`
}

type LedgerInfoBody struct {
	Method string         `json:"method"`
	Params []paramsstruct `json:"params"`
}

type paramsstruct struct {
	LedgerIndex  string `json:"ledger_index"`
	Accounts     bool   `json:"accounts"`
	Full         bool   `json:"full"`
	Transactions bool   `json:"transactions"`
	Expand       bool   `json:"expand"`
	OwnerFunds   bool   `json:"owner_funds"`
}

func getLedgerInfo(x int) {
	fmt.Println("Calling function for block number : ", x)
	requestMap := &LedgerInfoBody{

		Method: "ledger",
		Params: []paramsstruct{{
			LedgerIndex:  string(strconv.Itoa(x)),
			Accounts:     false,
			Full:         false,
			Transactions: true,

			Expand:     true,
			OwnerFunds: false}},
	}
	requestBytes, err := json.Marshal(requestMap)
	if err != nil {
		fmt.Println("Request Marshal Error", requestMap)
	}
	fmt.Println("Marshelled succesfully ")

	transaction, err := http.Post("https://s1.ripple.com:51234/", "application/json", bytes.NewReader(requestBytes))
	if err != nil {
		fmt.Println("Error making request")
	}

	ledgerDataResponse := AutoGeneratedX{}
	err = json.NewDecoder(transaction.Body).Decode(&ledgerDataResponse)
	switch {
	case err == io.EOF:
		fmt.Println("Empty Body from Response")
	case err != nil:
		fmt.Println("Error in Decoding Ledger DATA")
	}

	w := 0
	if ledgerDataResponse.Result.Error != "" {
		fmt.Println("Check this =====>>>>", ledgerDataResponse.Result.Error)
		fmt.Println("calling function to fetch details again after 5mins  Request")
		beeper.Add(1)
		go waitedRequest(requestBytes, x)

		return
	}

	for i := 0; i < len(ledgerDataResponse.Result.Ledger.Transactions); i++ {
		if ledgerDataResponse.Result.Ledger.Transactions[i].TransactionType == "Payment" {

			switch v := ledgerDataResponse.Result.Ledger.Transactions[i].Amount.(type) {
			case string:
				beeper.Add(1)
				fmt.Println("String: %v ", v)
				// fmt.Println("index : ", i)
				// fmt.Println("Hash ===> : ", ledgerDataResponse.Result.Ledger.Transactions[i].Hash)
				// fmt.Println("Fee ===> : ", ledgerDataResponse.Result.Ledger.Transactions[i].Fee)
				// fmt.Println("Sender ===> : ", ledgerDataResponse.Result.Ledger.Transactions[i].Account)
				// fmt.Println("Amount ====> : ", ledgerDataResponse.Result.Ledger.Transactions[i].Amount)
				// fmt.Println("Destination ===> :", ledgerDataResponse.Result.Ledger.Transactions[i].Destination)
				go broadCast(ledgerDataResponse.Result.Ledger.Transactions[i])

				w = w + 1
			default:
			}

		}
	}
	if w == 0 {
		fmt.Println("None of XRP PAYMENT TRANSACTIONS in ledger Index ", ledgerDataResponse.Result.Ledger.LedgerIndex)
	}

}

func waitedRequest(requestBytes []byte, x int) {
	time.Sleep(10 * time.Second)
	transaction, err := http.Post("https://s1.ripple.com:51234/", "application/json", bytes.NewReader(requestBytes))
	if err != nil {
		fmt.Println("Error making request")
	}

	ledgerDataResponse := AutoGeneratedX{}
	err = json.NewDecoder(transaction.Body).Decode(&ledgerDataResponse)
	switch {
	case err == io.EOF:
		fmt.Println("Empty Body from Response")
	case err != nil:
		fmt.Println("Error in Decoding Ledger DATA")
	}

	w := 0
	if ledgerDataResponse.Result.Error != "" {
		fmt.Println("Check this =====>>>>", ledgerDataResponse.Result.Error)
		return
	}

	for i := 0; i < len(ledgerDataResponse.Result.Ledger.Transactions); i++ {
		if ledgerDataResponse.Result.Ledger.Transactions[i].TransactionType == "Payment" {

			switch v := ledgerDataResponse.Result.Ledger.Transactions[i].Amount.(type) {
			case string:
				beeper.Add(1)
				fmt.Printf("String: %v", v)
				// fmt.Println("index : ", i)
				// fmt.Println("Hash ===> : ", ledgerDataResponse.Result.Ledger.Transactions[i].Hash)
				// fmt.Println("Fee ===> : ", ledgerDataResponse.Result.Ledger.Transactions[i].Fee)
				// fmt.Println("Sender ===> : ", ledgerDataResponse.Result.Ledger.Transactions[i].Account)
				// fmt.Println("Amount ====> : ", ledgerDataResponse.Result.Ledger.Transactions[i].Amount)
				// fmt.Println("Destination ===> :", ledgerDataResponse.Result.Ledger.Transactions[i].Destination)
				go broadCast(ledgerDataResponse.Result.Ledger.Transactions[i])
				w = w + 1
			default:
			}

		}
	}
	if w == 0 {
		fmt.Println("None of XRP PAYMENT TRANSACTIONS")
	}
	beeper.Done()

}

func broadCast(tx TransactionRes) {

	res, err := GetAllStreamQuery(db)
	if err != nil {
		fmt.Println("Error in getting url from DB")
	}

	for i := 0; i < len(res); i++ {
		beeper.Add(1)
		go makeRequestX(res[i].Webhookurl, res[i].Id, tx)

	}

	beeper.Done()

}

//make request to addresses of particular url
func makeRequestX(url string, streamID int, tx TransactionRes) {

	//get all the address for streamID
	res, err := getAddressfromStreamQuery(db, streamID)
	if err != nil {
		fmt.Println("Error in getting address for url from db")
	}

	for i := 0; i < len(res); i++ {

		if tx.Destination == res[i] {
			requestBytes, err := json.Marshal(tx)

			if err != nil {
				fmt.Println("Error Marshalling tx")
			}

			fmt.Println("Marshelled succesfully ")
			beeper.Add(1)
			POST(url, requestBytes)
		}

	}

	beeper.Done()
}

//make request to addresses of particular url
func makeRequest(url string, tx TransactionRes) {

	res, err := getAddressesforUrlfromDbQuery(db, url)

	if err != nil {
		fmt.Println("Error in getting address for url from db")
	}

	res1 := strings.Split(res[0], ",")

	t := ""
	for i := 0; i < len(res1); i++ {
		if i == 0 {
			t = strings.Replace(res1[i], "{", "", -1)
			fmt.Println(t)

		} else if i == len(res1)-1 {

			t = strings.Replace(res1[i], "}", "", -1)
			fmt.Println(t)

		} else {
			t = res1[i]
		}

		if tx.Destination == t {
			requestBytes, err := json.Marshal(tx)

			if err != nil {
				fmt.Println("Error Marshalling tx")
			}

			fmt.Println("Marshelled succesfully ")
			beeper.Add(1)
			POST(url, requestBytes)
		}

	}

	beeper.Done()
}

func POST(url string, requestBytes []byte) {

	fmt.Println("Broadcasting to Url", url)
	_, err := http.Post(url, "application/json", bytes.NewReader(requestBytes))
	if err != nil {
		fmt.Println("Error making request to webhook url")
	}
	beeper.Done()
}

// b, err := io.ReadAll(transaction.Body)
// // b, err := ioutil.ReadAll(resp.Body)  Go.1.15 and earlier
// if err != nil {
// 	log.Fatalln(err)
// }
// fmt.Println("LALALALALALALALALALALA ======================>>>>>>>>>>>>>>>>>>")

// fmt.Println(string(b))
