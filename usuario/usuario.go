package usuario

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"pureflix-go/db"
	"pureflix-go/jwt"
	"pureflix-go/mail"
	"pureflix-go/utils"

	"golang.org/x/crypto/bcrypt"
)

func RegistrarNuevoUsuario(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	datos, err := utils.RecibeDatosPost(r, nil)
	if err != 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"mensaje":"Error leyendo datos"}`))
		return
	}

	var metodoPagoID int

	switch datos["metodo_pago"] {
	case "tarjeta":
		metodoPagoID = 1
	case "transferencia":
		metodoPagoID = 2
	case "pago_facil":
		metodoPagoID = 3
	case "rapipago":
		metodoPagoID = 4
	default:
		metodoPagoID = 0
	}

	tokenValidacion, errorTokenValidacion := utils.GenerateToken()
	if errorTokenValidacion != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	mail.EnviarMailDeValidacion(datos["email"], tokenValidacion)

	consulta := `
		INSERT INTO usuario (nombre_usuario, email, pass, token_validacion ,metodo_pago)
		VALUES ($1, $2, $3, $4, $5)`

	hash, errorHash := bcrypt.GenerateFromPassword([]byte(datos["pass"]), bcrypt.DefaultCost)
	if errorHash != nil {
		fmt.Println("Error al generar la clave bcrypt:", err)
		utils.DevolverError(w, http.StatusInternalServerError)
		return
	}

	_, errbd := db.BaseDeDatos.Exec(
		consulta,
		datos["nombre_usuario"],
		datos["email"],
		hash,
		tokenValidacion,
		metodoPagoID,
	)

	if errbd != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"mensaje":"Error al registrar usuario"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"mensaje":"Usuario registrado correctamente"}`))
}
func ActivarCuenta(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	consulta := `
        UPDATE usuario
        SET activa = 1, token_validacion = ''
        WHERE token_validacion = $1;
    `
	res, err := db.BaseDeDatos.Exec(consulta, token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rows, err := res.RowsAffected()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if rows == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func RecuperarContrasena(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	datos, err := utils.RecibeDatosPost(r, nil)
	if err != 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"mensaje":"Error leyendo datosRecibidos"}`))
		return
	}

	consulta := `SELECT nombre_usuario FROM usuario WHERE nombre_usuario = $1 AND email = $2`
	type Usuario struct {
		NombreUsuario string `json:"nombre_usuario"`
		Email         string `json:"email"`
	}

	var u Usuario

	e := db.BaseDeDatos.QueryRow(consulta, datos["nombre_usuario"], datos["email"]).
		Scan(&u.NombreUsuario)

	if e != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"mensaje":"Usuario no encontrado"}`))
		return
	}

	tokenRecuperacion, errTokenRecuperacion := utils.GenerateToken()
	if errTokenRecuperacion != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	mail.EnviarMailDeRecuperacion(datos["email"], tokenRecuperacion)

	consultaDos := `UPDATE usuario SET token_recuperacion = $1 WHERE nombre_usuario = $2`
	db.BaseDeDatos.QueryRow(consultaDos, tokenRecuperacion, datos["nombre_usuario"])
}
func ActualizarContrasena(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	datos, err := utils.RecibeDatosPut(r, nil)
	if err != 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"mensaje":"Error leyendo datosRecibidos"}`))
		return
	}

	if datos["nueva_contrasena"] != datos["nueva_contrasena_rep"] {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"mensaje":"Las contraseñas no coinciden"}`))
		return
	}

	hash, errorHash := bcrypt.GenerateFromPassword([]byte(datos["nueva_contrasena"]), bcrypt.DefaultCost)
	if errorHash != nil {
		fmt.Println("Error al generar la clave bcrypt:", err)
		utils.DevolverError(w, http.StatusInternalServerError)
		return
	}

	consulta := `
        UPDATE usuario
        SET pass = $1, token_recuperacion = ''
        WHERE token_recuperacion = $2;
    `
	db.BaseDeDatos.QueryRow(consulta, hash, datos["token"])
}

func GetDatosUsuario(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Max-Age", "15")

	token := ""
	fmt.Sscanf(r.URL.Query().Get("token"), "%s", &token)

	usuario, errorUsuario := jwt.GetUsernameFromToken(token)
	if errorUsuario != nil {
		utils.DevolverError(w, http.StatusUnauthorized)
		return
	}

	type Usuario struct {
		NombreUsuario string `json:"nombre_usuario"`
		Email         string `json:"email"`
		FotoPerfil    string `json:"foto_perfil"`
		MetodoPago    int    `json:"metodo_pago"`
		//Favoritos  []Favorito `json:"favoritos"`
	}

	consulta := `SELECT nombre_usuario, email, foto_perfil, metodo_pago FROM usuario WHERE nombre_usuario = $1`
	row := db.BaseDeDatos.QueryRow(consulta, usuario)

	var u Usuario
	err := row.Scan(
		&u.NombreUsuario,
		&u.Email,
		&u.FotoPerfil,
		&u.MetodoPago,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.DevolverError(w, http.StatusNotFound)
			return
		}
		utils.DevolverError(w, http.StatusInternalServerError)
		return
	}

	respuestaJson, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(respuestaJson)
}
func ActualizarUsuario(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	datos, err := utils.RecibeDatosPut(r, nil)
	if err != 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"mensaje":"Error leyendo datos recibidos"}`))
		return
	}
	usuario, errorUsuario := jwt.GetUsernameFromToken(datos["token"])
	if errorUsuario != nil {
		utils.DevolverError(w, http.StatusUnauthorized)
		return
	}

	var metodoPagoID int

	switch datos["metodo_pago"] {
	case "tarjeta":
		metodoPagoID = 1
	case "transferencia":
		metodoPagoID = 2
	case "pago_facil":
		metodoPagoID = 3
	case "rapipago":
		metodoPagoID = 4
	default:
		metodoPagoID = 0
	}

	if datos["cambio_pass"] == "true" {
		hash, errorHash := bcrypt.GenerateFromPassword([]byte(datos["new_pass"]), bcrypt.DefaultCost)
		if errorHash != nil {
			fmt.Println("Error al generar la clave bcrypt:", err)
			utils.DevolverError(w, http.StatusInternalServerError)
			return
		}
		consulta := `UPDATE usuario SET pass = $1, metodo_pago = $2 WHERE nombre_usuario = $3;`
		_, errorConsulta := db.BaseDeDatos.Exec(consulta, hash, metodoPagoID, usuario)

		if errorConsulta != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"mensaje":"Error al actualizar los nuevos datos"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"mensaje":"Datos actualizados correctamente"}`))
	} else {
		consulta := `UPDATE usuario SET metodo_pago = $1 WHERE nombre_usuario = $2;`
		_, errorConsulta := db.BaseDeDatos.Exec(consulta, metodoPagoID, usuario)

		if errorConsulta != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"mensaje":"Error al actualizar los nuevos datos"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"mensaje":"Datos actualizados correctamente"}`))
	}
}
func ActualizarFotoPerfil(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	datos, err := utils.RecibeDatosPut(r, nil)
	if err != 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"mensaje":"Error leyendo datos recibidos"}`))
		return
	}

	usuario, errorUsuario := jwt.GetUsernameFromToken(datos["token"])
	if errorUsuario != nil {
		utils.DevolverError(w, http.StatusUnauthorized)
		return
	}

	consulta := `UPDATE usuario SET foto_perfil = $1 WHERE nombre_usuario = $2;`
	_, errorConsulta := db.BaseDeDatos.Exec(consulta, datos["foto_perfil"], usuario)

	if errorConsulta != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"mensaje":"Error al actualizar la foto de perfil"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"mensaje":"Foto de perfil actualizada correctamente"}`))
}
