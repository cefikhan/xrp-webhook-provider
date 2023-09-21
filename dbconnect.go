package main

import (
	"github.com/go-pg/pg"
)

type xrpledger_entries struct {
	Id       int `json:"id"`
	Ledgerno int `json:"ledgerno"`
}

func new(addr, user, password, name string) (*pg.DB, error) {

	opts := &pg.Options{
		Addr:     addr,
		User:     user,
		Password: password,
		Database: name,
	}

	db := pg.Connect(opts)
	_, err := db.Exec("SELECT 1")
	if err != nil {
		return nil, err
	}

	return db, nil
}
