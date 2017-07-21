package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-oci8"
	//_ "gopkg.in/rana/ora.v4"
)

func main() {

	// connect
	conn, err := sql.Open("oci8", "my_owner/my_owner@oracle-db-dev:1521/ORCL.localdomain") // mattn
	//conn, err := sql.Open("ora", "my_owner/my_owner@oracle-db-dev:1521/ORCL.localdomain") // rana
	if err != nil {
		panic(err)
	}

	// ping
	if err = conn.Ping(); err != nil {
		panic(err)
	}

	// try query
	stmt := `SELECT id, ip FROM my_owner.ips`
	rows, err := conn.Query(stmt)
	if err != nil {
		panic(err)
	}
	defer closeRows(rows)

	for rows.Next() {
		if rows.Err() != nil {
			break
		}
		var (
			id int64
			ip      int64
		)
		if err = rows.Scan(&id, &ip); err != nil {
			panic(err)
		}
		fmt.Printf("%d, %d\n", id, ip)
	}
	if err = rows.Err(); err != nil {
		panic(err)
	}
}

func closeRows(rows *sql.Rows) {
	if closeErr := rows.Close(); closeErr != nil {
		panic(closeErr)
	}
}
