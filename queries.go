package main

import (
	"fmt"

	"github.com/go-pg/pg"
)

type SelectUserQueryResult struct {
	Id           int    `sql:"id"`
	Userpassword string `sql:"userpassword"`
}

type SelectUserStreamQueryResult struct {
	Id         int    `sql:"id"`
	Webhookurl string `sql:"webhookurl"`
}

type SelectAddresseswithIDQueryResult struct {
	ID      int    `sql:"id"`
	Address string `sql:"address"`
}

func Select(db *pg.DB) ([]xrpledger_entries, error) {
	ledgerIndices := []xrpledger_entries{}
	_, err := db.Query(&ledgerIndices, `SELECT ledgerno FROM xrpledger_entries`)
	if err != nil {
		return nil, err
	}
	return ledgerIndices, nil
}

func getLatestLedgerNofromDb(db *pg.DB) (int, error) {
	var currentLedgerIndex int
	_, err := db.Query(&currentLedgerIndex, `SELECT ledgerno FROM xrpledger_entries`)
	if err != nil {
		return 0, err

	}
	return currentLedgerIndex, nil
}

func Insert(db *pg.DB, ledgerIndex int) {

	_, err := db.Exec(`INSERT INTO xrpledger_entries (ledgerno) VALUES (?)
	RETURNING id`, ledgerIndex)
	if err != nil {
		fmt.Println("Insert Query Error")
	}

}

func UpdateLedgerNoInDb(db *pg.DB, new int) (string, error) {

	var currentLedgerIndex int
	_, err := db.Query(&currentLedgerIndex, `SELECT ledgerno FROM xrpledger_entries`)
	if err != nil {
		return "", err
	}

	_, err = db.Exec(`UPDATE xrpledger_entries set ledgerno = ? where ledgerno = ?`, new, currentLedgerIndex)
	if err != nil {
		return "", err

	}
	return "", nil
}

func InsertAddressQuery(db *pg.DB, Address string) error {

	_, err := db.Exec(`INSERT INTO addresses (address) VALUES (?)
	RETURNING id`, Address)
	if err != nil {
		return err
	}
	return nil

}

func getAddressesfromDbQuery(db *pg.DB) ([]string, error) {
	var addresses []string
	_, err := db.Query(&addresses, `SELECT address FROM addresses`)
	if err != nil {
		return nil, err

	}
	return addresses, nil
}

func getURLfromDbQuery(db *pg.DB) (string, error) {
	var urls string
	_, err := db.Query(&urls, `SELECT url FROM webhookurls`)
	if err != nil {
		return "", err

	}
	return urls, nil
}

func UpdateURLfromDb(db *pg.DB, new, old string) (string, error) {

	_, err := db.Exec(`UPDATE webhookurls set url = ? where url = ?`, new, old)
	if err != nil {
		return "", err

	}
	return "", nil
}

func InsertUrlQuery(db *pg.DB, url string) error {

	_, err := db.Exec(`INSERT INTO webhookurls (url) VALUES (?)
	RETURNING id`, url)
	if err != nil {
		return err
	}
	return nil

}

//Register URL into urlAddress table
func InsertWebhookUrlAddressQuery(db *pg.DB, url string) error {

	fmt.Println("URL -====>>>>", url)
	_, err := db.Exec(`INSERT INTO webhookurladdress (url) VALUES (?)
	RETURNING id`, url)
	if err != nil {
		return err
	}
	return nil

}

func InsertAddressInWebhookUrlsQuery(db *pg.DB, Address string, id int) error {

	_, err := db.Exec(`UPDATE webhookurladdress SET addresses = array_append(addresses,?) WHERE id = ?`, Address, id)
	if err != nil {
		return err
	}
	return nil

}

func getwebhookURLsfromDbQuery(db *pg.DB) ([]string, error) {
	var urls []string
	_, err := db.Query(&urls, `SELECT url FROM webhookurladdress`)
	if err != nil {
		return nil, err

	}
	return urls, nil
}

func getAddressesforUrlfromDbQuery(db *pg.DB, url string) ([]string, error) {
	var addresses []string
	_, err := db.Query(&addresses, `SELECT addresses FROM webhookurladdress where url = ?`, url)
	if err != nil {
		return nil, err

	}
	return addresses, nil
}

func RegisterUserQuery(db *pg.DB, Username string, Email string, Password string) error {

	_, err := db.Exec(`INSERT INTO users (username,email,userpassword) VALUES (?,?,?)
	RETURNING id`, Username, Email, Password)
	if err != nil {
		return err
	}
	return nil

}

func GetUserDetailsQuery(db *pg.DB, Username string) (SelectUserQueryResult, error) {
	userResult := SelectUserQueryResult{}
	_, err := db.Query(&userResult, `SELECT id, userpassword FROM users where username = ?`, Username)
	if err != nil {
		return SelectUserQueryResult{}, err

	}
	return userResult, nil
}

func CreateStreamQuery(db *pg.DB, userID int, webhookurl string) error {
	_, err := db.Exec(`INSERT INTO streams (userid,webhookurl) VALUES (?,?)
	RETURNING id`, userID, webhookurl)
	if err != nil {
		return err
	}
	return nil

}

func UpdateStreamQuery(db *pg.DB, streamID int, webhookurl string) error {
	_, err := db.Exec(`UPDATE streams set webhookurl = ? where id = ?`, webhookurl, streamID)
	if err != nil {
		return err

	}
	return nil

}

func GetAllStreamQuery(db *pg.DB) ([]SelectUserStreamQueryResult, error) {
	streamResult := []SelectUserStreamQueryResult{}
	_, err := db.Query(&streamResult, `SELECT id, webhookurl FROM streams`)
	if err != nil {
		return nil, err

	}
	return streamResult, nil
}

func GetUserStreamQuery(db *pg.DB, userID int) (SelectUserStreamQueryResult, error) {
	streamResult := SelectUserStreamQueryResult{}
	_, err := db.Query(&streamResult, `SELECT id, webhookurl FROM streams where userid = ?`, userID)
	if err != nil {
		return SelectUserStreamQueryResult{}, err

	}
	return streamResult, nil
}

func GetUserAllStreamQuery(db *pg.DB, userID int) ([]SelectUserStreamQueryResult, error) {
	streamResult := []SelectUserStreamQueryResult{}
	_, err := db.Query(&streamResult, `SELECT id, webhookurl FROM streams where userid = ?`, userID)
	if err != nil {
		return nil, err

	}
	return streamResult, nil
}

func registerAddressToStreamQuery(db *pg.DB, streamID int, address string) error {
	_, err := db.Exec(`INSERT INTO addresses (streamid,address) VALUES (?,?)
	RETURNING id`, streamID, address)
	if err != nil {
		return err
	}
	return nil
}

func getAddressfromStreamQuery(db *pg.DB, streamID int) ([]string, error) {
	var streamResult []string
	_, err := db.Query(&streamResult, `SELECT address FROM addresses where streamid = ?`, streamID)
	if err != nil {
		return nil, err

	}
	return streamResult, nil
}

func getAddresswithIDsfromStreamQuery(db *pg.DB, streamID int) ([]SelectAddresseswithIDQueryResult, error) {
	streamResult := []SelectAddresseswithIDQueryResult{}
	_, err := db.Query(&streamResult, `SELECT id,address FROM addresses where streamid = ?`, streamID)
	if err != nil {
		return nil, err

	}
	return streamResult, nil
}

func deleteAddressfromStreamQuery(db *pg.DB, streamID int, addressID int) error {
	_, err := db.Exec(`DELETE FROM addresses where streamid = ? and id=?`, streamID, addressID)
	if err != nil {
		return err

	}
	return nil
}
