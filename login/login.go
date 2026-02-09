package login

import (
	"encoding/json"
	"net/http"
	"pureflix-go/db"
	"pureflix-go/jwt"
	"pureflix-go/utils"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	datos, errDatos := utils.RecibeDatosPost(r, nil)
	if errDatos != 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"mensaje":"Datos inválidos"}`))
		return
	}

	var passHasheado string

	consulta := `
		SELECT pass 
		FROM usuario 
		WHERE nombre_usuario = $1 AND activa = 1;
	`

	err := db.BaseDeDatos.
		QueryRow(consulta, datos["nombre_usuario"]).
		Scan(&passHasheado)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"mensaje":"Usuario inexistente o cuenta no activa"}`))
		return
	}

	if err = bcrypt.CompareHashAndPassword(
		[]byte(passHasheado),
		[]byte(datos["pass"]),
	); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"mensaje":"Usuario o contraseña incorrectos"}`))
		return
	}

	tokenString, err := jwt.GenerateJWT(datos["nombre_usuario"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	type RespuestaLogin struct {
		Token string `json:"token"`
	}

	resp := RespuestaLogin{
		Token: tokenString,
	}

	json.NewEncoder(w).Encode(resp)
}

func ValidarToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	token := r.URL.Query().Get("token")
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, err := jwt.GetUsernameFromToken(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func ResetDB(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "Application/json; charset=utf-8")
	w.Header().Set("Access-Control-Max-Age", "15")
}
