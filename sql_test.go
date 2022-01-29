package golangdatabase

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestExecSql(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	insert_data := "INSERT INTO customer VALUES ('pambudi', 'Pambudi')"
	_, err := db.ExecContext(ctx, insert_data)

	if err != nil {
		panic(err)
	}
	fmt.Println("Insert Data Success")

}

func TestQuerySql(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	sql := "SELECT id, name FROM customer"
	rows, err := db.QueryContext(ctx, sql)
	if err != nil {
		panic(err)
	}

	/* Iterate data rows */
	for rows.Next() {
		var id, name string
		err := rows.Scan(&id, &name)
		if err != nil {
			panic(err)
		}

		fmt.Println("Id", id)
		fmt.Println("Name", name)

	}
	defer rows.Close()
}

func TestQuerySqlComplex(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "SELECT id, name, email, balance, rating, birth_date, married, created_at FROM customer"
	rows, err := db.QueryContext(ctx, script)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, name string
		var email sql.NullString
		var balance int32
		var rating sql.NullFloat64
		var birthDate sql.NullTime
		var createdAt time.Time
		var married bool

		err = rows.Scan(&id, &name, &email, &balance, &rating, &birthDate, &married, &createdAt)
		if err != nil {
			panic(err)
		}
		fmt.Println("================")
		fmt.Println("Id:", id)
		fmt.Println("Name:", name)
		if email.Valid {
			fmt.Println("Email:", email.String)
		}
		fmt.Println("Balance:", balance)
		if rating.Valid {
			fmt.Println("Rating", rating.Float64)
		}
		if birthDate.Valid {
			fmt.Println("Birth Date:", birthDate.Time)
		}
		fmt.Println("Married:", married)
		fmt.Println("Created At:", createdAt)
	}
}

func TestSQLInjection(t *testing.T) {
	db := GetConnection()
	defer db.Close()
	ctx := context.Background()

	username := "admin';#"
	password := "admin"

	sql := "SELECT username FROM user WHERE username = '" + username + "' AND password='" + password + "' LIMIT 1"
	fmt.Println(sql)
	rows, err := db.QueryContext(ctx, sql)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	if rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			panic(err)
		}

		fmt.Println("Login Success!")
	} else {
		fmt.Println("Login Failed!")
	}

}

func TestSQLInjectionSafe(t *testing.T) {
	db := GetConnection()
	defer db.Close()
	ctx := context.Background()

	username := "admin"
	password := "admin"

	sql := "SELECT username FROM user WHERE username = ? AND password = ? LIMIT 1"
	fmt.Println(sql)
	rows, err := db.QueryContext(ctx, sql, username, password)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	if rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			panic(err)
		}

		fmt.Println("Login Success!")
	} else {
		fmt.Println("Login Failed!")
	}

}

func TestExecSqlParamater(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	username := "eko'; DROP TABLE user; #"
	password := "eko"

	insert_data := "INSERT INTO user VALUES (?, ?)"
	_, err := db.ExecContext(ctx, insert_data, username, password)

	if err != nil {
		panic(err)
	}
	fmt.Println("Insert Data Success")

}

func TestAutoIncrement(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	email := "pambudi@gmail.com"
	comment := "This is comment"

	insert_data := "INSERT INTO comments(email, comment) VALUES (?, ?)"
	result, err := db.ExecContext(ctx, insert_data, email, comment)

	if err != nil {
		panic(err)
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}

	fmt.Println("Insert Data Success with id", lastId)

}

func TestPrepareStatement(test *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	query := "INSERT INTO comments(email, comment) VALUES (?, ?)"
	statement, err := db.PrepareContext(ctx, query)
	if err != nil {
		panic(err)
	}

	defer statement.Close()

	for i := 0; i < 10; i++ {
		email := "eko" + strconv.Itoa(i) + "@gmail.com"
		comment := "comment ke" + strconv.Itoa(i)

		result, err := statement.ExecContext(ctx, email, comment)
		if err != nil {
			panic(err)
		}

		lastInsertId, err := result.LastInsertId()

		if err != nil {
			panic(err)
		}

		fmt.Println("Last Id Ke-", lastInsertId)
	}

}

func TestTransaction(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	query := "INSERT INTO comments(email, comment) VALUES (?, ?)"

	for i := 0; i < 10; i++ {
		email := "eko" + strconv.Itoa(i) + "@gmail.com"
		comment := "comment ke" + strconv.Itoa(i)

		result, err := tx.ExecContext(ctx, query, email, comment)
		if err != nil {
			panic(err)
		}

		lastInsertId, err := result.LastInsertId()

		if err != nil {
			panic(err)
		}

		fmt.Println("Last Id Ke-", lastInsertId)
	}

	// Transaction
	// err = tx.Commit()
	err = tx.Rollback()

	if err != nil {
		panic(err)
	}
}
