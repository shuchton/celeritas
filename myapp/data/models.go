package data

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/upper/db/v4/adapter/mysql"
	"github.com/upper/db/v4/adapter/postgresql"

	db2 "github.com/upper/db/v4"
)

var db *sql.DB
var upper db2.Session

type Models struct {
	// any models inserted here (and in the New function)
	// are easily accessible throught the entire application
	Users  User
	Tokens Token
}

func New(databasePool *sql.DB) Models {
	db = databasePool

	switch os.Getenv("DATABASE_TYPE") {
	case "mysql", "mariadb":
		upper, _ = mysql.New(databasePool)
	case "postgres", "postgresql":
		upper, _ = postgresql.New(databasePool)
	default:
		// do nothing
	}

	return Models{
		Users:  User{},
		Tokens: Token{},
	}
}

func getInsertID(i db2.ID) int {
	idType := fmt.Sprintf("%T", i)
	if idType == "int64" {
		return int(i.(int64))
	}
	return i.(int)
}
