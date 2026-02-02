package login

import (
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"pureflix-go/db"
	"pureflix-go/jwt"
	"pureflix-go/utils"
)

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "Application/json; charset=utf-8")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	datos, error := utils.RecibeDatosPost(r, nil)
	if error != 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"mensaje":"Datos inválidos"}`))
		return
	}

	var passHaseado string

	consulta := `SELECT pass FROM usuario WHERE nombre_usuario = $1 AND activa = $2;`

	err := db.BaseDeDatos.QueryRow(consulta, datos["nombre_usuario"], 1).Scan(&passHaseado)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"No se ha activado la cuenta todavía. Revisa tu email."}`))
		return
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	if err = bcrypt.CompareHashAndPassword([]byte(passHaseado), []byte(datos["pass"])); err != nil {
		// If the two passwords don't match, return a 401 status
		fmt.Printf("Clave incorrecta\n")
		//db.Auditoria("no requiere token", "Error Login", "Usuario: ["+creds.Username+"]. Contraseña incorrecta.", 0, "", false)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// If we reach this point, that means the users password was correct, and that they are authorized
	tokenString, _ := jwt.GenerateJWT(datos["nombre_usuario"])

	// Envía el token JWT al cliente
	w.Header().Set("Content-Type", "application/json")

	type Usuario struct {
		NombreUsuario string `json:"nombre_usuario"`
		Email         string `json:"email"`
		FotoPerfil    string `json:"foto_perfil"`
		Token         string `json:"token"`
		Activa        int    `json:"activa"`
		MetodoPago    int32  `json:"metodo_pago"`
	}

	var u Usuario
	u.Token = tokenString
	u.NombreUsuario = datos["nombre_usuario"]

	consulta = `SELECT email, metodo_pago, activa, foto_perfil FROM usuario WHERE nombre_usuario=$1;`
	db.BaseDeDatos.QueryRow(consulta, u.NombreUsuario).Scan(&u.Email, &u.MetodoPago, &u.Activa, &u.FotoPerfil)

	usuJson, err := json.Marshal(u)
	if err != nil {
		fmt.Println("Error:", err)
	}
	w.Write(usuJson)
	//	fmt.Printf("Clave correcta, generando token: \n[%v]\n", tokenString)
}

func ValidaToken(w http.ResponseWriter, r *http.Request) {
	datos, error := utils.RecibeDatosPost(r, nil)
	if error != 0 {
		utils.DevolverError(w, error)
		return
	}
	_, errorUsuario := jwt.GetUsernameFromToken(datos["token"])
	if errorUsuario != nil { // no se obtiene usuario a partir del token, debe estar expirado
		utils.DevolverError(w, http.StatusUnauthorized)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func ResetDB(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "Application/json; charset=utf-8")
	w.Header().Set("Access-Control-Max-Age", "15")
}
