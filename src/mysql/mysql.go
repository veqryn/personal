package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	// connect
	conn, err := sql.Open("mysql", "root:password@/my_schema")
	if err != nil {
		panic(err)
	}

	// ping
	if err := conn.Ping(); err != nil {
		panic(err)
	}

	// try query
	stmt := `SELECT id, ip FROM my_schema.ips`
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
		if err := rows.Scan(&id, &ip); err != nil {
			panic(err)
		}
		fmt.Printf("%d, %d\n", id, ip)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
}

func closeRows(rows *sql.Rows) {
	if closeErr := rows.Close(); closeErr != nil {
		panic(closeErr)
	}
}
