package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	//_ "golang.org/x/crypto/bcrypt"
	_ "strings"
)

type ErrorJson struct {
	Mensaje string
}

func devolverMensaje(w http.ResponseWriter, mensaje string) {
	var ej ErrorJson
	ej.Mensaje = mensaje
	j, _ := json.Marshal(ej)
	w.Write(j)
}

var BaseDeDatos *sql.DB

func InicializarBD() {
	var err error

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		dbname,
	)

	BaseDeDatos, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("Error iniciando la BD")
		panic(err)
	}
}
