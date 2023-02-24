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

var (
	walletCreate       = "Create"
	walletUpdateAmount = "UpdateAmount"
	walletQueries      = map[string]string{
		walletCreate:       "INSERT INTO mst_membership_wallet VALUES ($1, $2, $3)",
		walletUpdateAmount: "update mst_membership_wallet set wallet_amount=wallet_amount+$1,description=$2 where membership_id=$3",
	}
)

type WalletRepo struct {
	db *sqlx.DB
	ps map[string]*sqlx.Stmt
}

func NewWalletRepo(db *sqlx.DB) *WalletRepo {
	ps := make(map[string]*sqlx.Stmt, len(walletQueries))
	for idx, v := range walletQueries {
		p, err := db.Preparex(v)
		if err != nil {
			log.Fatalln(err)
		}

		ps[idx] = p
	}

	return &WalletRepo{db: db, ps: ps}
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

	walletRepo := NewWalletRepo(db)
	//err = WithStmtTransaction(db, func(tx *sqlx.Tx) error {
	//	wallet1 := Wallet{1, 10000, "Isi wallet"}
	//	wallet2 := Wallet{2, 5000, "Isi wallet"}
	//	_, err := tx.Stmtx(walletRepo.ps[walletCreate]).Exec(wallet1.MembershipId, wallet1.WalletAmount, wallet1.Description)
	//	_, err = tx.Stmtx(walletRepo.ps[walletCreate]).Exec(wallet2.MembershipId, wallet2.WalletAmount, wallet2.Description)
	//	return err
	//})

	// Simulasi transfer
	totalAmount := 2000.0
	err = WithStmtTransaction(db, func(tx *sqlx.Tx) error {
		wallet1 := Wallet{1, -1 * totalAmount, fmt.Sprintf("Kredit:%v", totalAmount)}
		wallet2 := Wallet{2, totalAmount, fmt.Sprintf("Debit:%v", totalAmount)}
		_, err = tx.Stmtx(walletRepo.ps[walletUpdateAmount]).Exec(wallet1.WalletAmount, wallet1.Description, wallet1.MembershipId)
		_, err = tx.Stmtx(walletRepo.ps[walletUpdateAmount]).Exec(wallet2.WalletAmount, wallet2.Description, wallet2.MembershipId)
		return err
	})
	if err != nil {
		log.Println(err.Error())
	}
}
