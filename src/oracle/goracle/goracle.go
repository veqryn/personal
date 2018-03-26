package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "gopkg.in/goracle.v2" // Oracle driver (required)
)

const oracleConnectionString string = `
%s/%s@(DESCRIPTION =
	(SDU=32767)
	(ENABLE=BROKEN)
	(ADDRESS_LIST =
		(ADDRESS =
			(PROTOCOL = TCP)
			(HOST = %s)
			(PORT = %d)
		)
	)
	(CONNECT_DATA =
		(SERVICE_NAME = %s)
		(SERVER = DEDICATED)
	)
)
`

func main() {

	dbUsername := "asuser"
	dbPassword := "Oradoc_db1"
	dbHost := "localhost"
	dbPort := 1521
	dbSID := "ORCLPDB1"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	connString := fmt.Sprintf(oracleConnectionString, dbUsername, dbPassword, dbHost, dbPort, dbSID)
	conn, err := sqlx.ConnectContext(ctx, "goracle", connString) // SQLX will ping after connecting
	if err != nil {
		panic(err)
	}

	var (
		regexp  = `(?m-xis:^X-pid:\s+(\d+)\s*$)`
		realmID = 8806400403
		id      = 0
	)

	// Inserts
	stmt, err := conn.PreparexContext(ctx, `
		INSERT INTO master_owner.mail_sig_regexp_t (
			mail_sig_regexp_id,
			regexp,
			realm_id,
			regexp_match_type,
			status
		) VALUES (
			AS_OWNER.mail_sig_regexp_s.nextval,
			:regexp,
			:realm_id,
			1,
			0
		) RETURNING mail_sig_regexp_id into :mail_sig_regexp_id`)

	if err != nil {
		panic(err)
	}

	// Without naming the variables (just use same order as they appear)
	_, err = stmt.ExecContext(ctx,
		regexp,
		realmID,
		sql.Out{Dest: &id},
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(id)

	// With naming the variables
	_, err = stmt.ExecContext(ctx,
		sql.Named("regexp", regexp),
		sql.Named("realm_id", realmID),
		sql.Named("mail_sig_regexp_id", sql.Out{Dest: &id}),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(id)

	// Queries and scanning:
	type MailSigRegexp struct {
		ID      int64  `db:"MAIL_SIG_REGEXP_ID"`
		Regexp  string `db:"REGEXP"`
		RealmID int64  `db:"REALM_ID"`
		Status  int64  `db:"STATUS"`
	}

	// Query a single struct and scan it in automatically
	var mailSigRegexp MailSigRegexp
	err = conn.GetContext(ctx,
		&mailSigRegexp,
		`SELECT mail_sig_regexp_id, regexp, realm_id, status
		 FROM master_owner.mail_sig_regexp_t
		 WHERE mail_sig_regexp_id = :id`,
		id,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(mailSigRegexp)

	// Query multiple structs and scan them in automatically
	var mailSigRegexps []MailSigRegexp
	err = conn.SelectContext(ctx,
		&mailSigRegexps,
		`SELECT mail_sig_regexp_id, regexp, realm_id, status
		 FROM master_owner.mail_sig_regexp_t`,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(mailSigRegexps)
}
