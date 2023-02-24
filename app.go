package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

type Wallet struct {
	MembershipId int     `db:"membership_id"`
	WalletAmount float64 `db:"wallet_amount"`
	Description  string  `db:"description"`
}

func main() {
	dbHost := "db-postgresql-enigma-do-user-279248-0.b.db.ondigitalocean.com"
	dbPort := "25060"
	dbName := "defaultdb"
	dbUser := "doadmin"
	dbPassword := "AVNS_V2yOmmAI1Kd9ZhBpSZJ"
	dataSourceName := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()
	// DDL
	//create table public.mst_membership_wallet
	//(
	//    membership_id integer primary key,
	//    wallet_amount double precision,
	//    description   varchar(20)
	//);
	//tx := db.MustBegin()
	//tx.NamedExec("INSERT INTO mst_membership_wallet VALUES (:membership_id, :wallet_amount, :description)", &Wallet{1, 10000, "Isi wallet"})
	//tx.NamedExec("INSERT INTO mst_membership_wallet VALUES (:membership_id, :wallet_amount, :description)", &Wallet{2, 5000, "Isi wallet"})
	//tx.Commit()

	// Simulasi transfer
	totalAmount := 1000.0
	tx := db.MustBegin()
	_, err = tx.NamedExec("update mst_membership_wallet set wallet_amount=wallet_amount+:wallet_amount,description=:description where membership_id=:membership_id", &Wallet{1, -1 * totalAmount, fmt.Sprintf("Kredit:%v", totalAmount)})
	_, err = tx.NamedExec("update mst_membership_wallet set wallet_amount=wallet_amount+:wallet_amount,description=:description where membership_id=:membership_id", &Wallet{2, totalAmount, fmt.Sprintf("Debit:%v", totalAmount)})
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Println(rollbackErr)
		}
	} else {
		err = tx.Commit()
	}
}
