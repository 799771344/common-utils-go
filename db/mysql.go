package db

import (
	"database/sql"
	"fmt"
	"log"
)

type Mysql struct {
	db *sql.DB
}

func NewDB(user, password, host, port, dbname string) (*Mysql, error) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbname)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	return &Mysql{db: db}, nil
}

func (d *Mysql) Close() error {
	return d.db.Close()
}

func (d *Mysql) QueryRow(query string, args ...interface{}) *sql.Row {
	return d.db.QueryRow(query, args...)
}

func (d *Mysql) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return d.db.Query(query, args...)
}

func (d *Mysql) Exec(query string, args ...interface{}) (sql.Result, error) {
	result, err := d.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func main() {
	db, err := NewDB("username", "password", "127.0.0.1", "3306", "mydatabase")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 查询单行数据
	row := db.QueryRow("SELECT id, name FROM users WHERE id = ?", 1)
	if err != nil {
		log.Fatal(err)
	}

	// 查询多行数据
	rows, err := db.Query("SELECT * FROM users WHERE age = ?", 18)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("id: %d, name: %s\n", id, name)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// 插入数据
	result, err := db.Exec("INSERT INTO users (name, age) VALUES (?, ?)", "Alice", 20)
	if err != nil {
		log.Fatal(err)
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("LastInsertID: %d, RowsAffected: %d\n", lastInsertID, rowsAffected)
}
