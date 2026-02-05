package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"pureflix-go/db"
	"pureflix-go/jwt"
	"pureflix-go/login"
	"pureflix-go/mail"
	"pureflix-go/utils"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func enableCORS(router *mux.Router) {
	fmt.Printf("Habilitando cors\n")
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}).Methods(http.MethodOptions)
	router.Use(middlewareCors)
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, DELETE, PUT")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			next.ServeHTTP(w, req)
		})
}

func main() {

	// 🔥 Cargar variables de entorno
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error cargando el archivo .env")
	}

	router := mux.NewRouter()

	// Inicializar BD (usa os.Getenv)
	db.InicializarBD()

	enableCORS(router)

	// Endpoints
	router.HandleFunc("/registrarNuevoUsuario", RegistrarNuevoUsuario)
	router.HandleFunc("/activar", ActivarCuenta)

	router.HandleFunc("/login", login.Login)

	router.HandleFunc("/recuperarContrasena", RecuperarContrasena)
	router.HandleFunc("/actualizarContrasena", ActualizarContrasena)

	router.HandleFunc("/getDatosUsuario", GetDatosUsuario)
	router.HandleFunc("/actualizarUsuario", ActualizarUsuario)

	router.HandleFunc("/getContenido", GetContenido)
	router.HandleFunc("/getCategorias", GetCategorias)

	router.HandleFunc("/getImagen", GetImagen)

	fmt.Println("Iniciando servidor")

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "9000"
	}

	log.Fatal(http.ListenAndServe(":"+port, router))
}

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
	if errorUsuario != nil { // no se obtiene usuario a partir del token, debe estar expirado
		utils.DevolverError(w, http.StatusUnauthorized)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	type Usuario struct {
		NombreUsuario string         `json:"nombre_usuario"`
		Email         string         `json:"email"`
		FotoPerfil    sql.NullString `json:"foto_perfil"`
		MetodoPago    int            `json:"metodo_pago"`
	}

	consulta := `SELECT nombre_usuario, email, COALESCE(foto_perfil, ''), metodo_pago FROM usuario WHERE nombre_usuario = $1`
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
	if errorUsuario != nil { // no se obtiene usuario a partir del token, debe estar expirado
		utils.DevolverError(w, http.StatusUnauthorized)
	} else {
		w.WriteHeader(http.StatusOK)
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

func GetContenido(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Max-Age", "15")
	w.WriteHeader(http.StatusOK)

	idTipo := 0
	fmt.Sscanf(r.URL.Query().Get("tipo"), "%d", &idTipo)

	idCategoria := 0
	fmt.Sscanf(r.URL.Query().Get("categoria"), "%d", &idCategoria)

	palabra := ""
	fmt.Sscanf(r.URL.Query().Get("palabra"), "%s", &palabra)

	consulta := `
	SELECT DISTINCT 
		c.id,
		c.link_trailer,
		c.titulo,
		c.imagen,
		c.resumen,
		EXTRACT(YEAR FROM c.anio) AS anio,
		c.tipo_contenido
	FROM contenido c
	`

	filtros := ""

	if idCategoria > 0 {
		consulta += `
		JOIN contenido_genero cg ON c.id = cg.id_contenido
		JOIN genero g ON cg.id_genero = g.id
		`
		filtros += fmt.Sprintf(" AND g.id = %d ", idCategoria)
	}

	if idTipo > 0 {
		filtros += fmt.Sprintf(" AND c.tipo_contenido = %d ", idTipo)
	}

	if palabra != "" {
		filtros += " AND unaccent(c.titulo) ILIKE unaccent('%" + palabra + "%')"
	}

	consulta += `
	WHERE 1=1
	` + filtros

	rows, errbd := db.BaseDeDatos.Query(consulta)
	if errbd != nil {
		fmt.Println("Error en la base de datos:", errbd)
		return
	}
	defer rows.Close()

	type Actor struct {
		Nombre    string `json:"nombre"`
		Apellido  string `json:"apellido"`
		Wikipedia string `json:"wikipedia"`
	}

	type Genero struct {
		Id          int    `json:"id"`
		Descripcion string `json:"descripcion"`
	}

	type Capitulo struct {
		Nro      int    `json:"nro"`
		Titulo   string `json:"titulo"`
		Duracion int    `json:"duracion"`
	}

	type Temporada struct {
		Nro       int        `json:"nro"`
		Capitulos []Capitulo `json:"capitulos"`
	}

	type Contenido struct {
		Id           int         `json:"id"`
		Link_trailer string      `json:"link_trailer"`
		Titulo       string      `json:"titulo"`
		Genero       []Genero    `json:"genero"`
		Imagen       string      `json:"imagen"`
		Actores      []Actor     `json:"actores"`
		Resumen      string      `json:"resumen"`
		Anio         int         `json:"anio"`
		Tipo         int         `json:"tipo"`
		Temporadas   []Temporada `json:"temporadas,omitempty"`
	}

	var Contenidos []Contenido

	for rows.Next() {
		var c Contenido
		rows.Scan(
			&c.Id,
			&c.Link_trailer,
			&c.Titulo,
			&c.Imagen,
			&c.Resumen,
			&c.Anio,
			&c.Tipo,
		)

		// Géneros
		consultaDos := `
		SELECT g.id, g.descripcion
		FROM genero g
		JOIN contenido_genero cg ON g.id = cg.id_genero
		WHERE cg.id_contenido = $1;
		`
		rowsDos, _ := db.BaseDeDatos.Query(consultaDos, c.Id)
		for rowsDos.Next() {
			var aux Genero
			rowsDos.Scan(&aux.Id, &aux.Descripcion)
			c.Genero = append(c.Genero, aux)
		}
		rowsDos.Close()

		// Actores
		consultaTres := `
		SELECT a.nombre, a.apellido, a.wikipedia
		FROM actor a
		JOIN contenido_actor ca ON a.id = ca.id_actor
		WHERE ca.id_contenido = $1;
		`
		rowsTres, _ := db.BaseDeDatos.Query(consultaTres, c.Id)
		for rowsTres.Next() {
			var aux Actor
			rowsTres.Scan(&aux.Nombre, &aux.Apellido, &aux.Wikipedia)
			c.Actores = append(c.Actores, aux)
		}
		rowsTres.Close()

		// Temporadas y capítulos
		if c.Tipo == 2 {
			consultaTemporadas := `
			SELECT id, nro
			FROM temporada
			WHERE id_serie = $1
			ORDER BY nro;
			`
			rowsTemp, _ := db.BaseDeDatos.Query(consultaTemporadas, c.Id)

			for rowsTemp.Next() {
				var temp Temporada
				var tempId int
				rowsTemp.Scan(&tempId, &temp.Nro)

				consultaCapitulos := `
				SELECT nro, titulo, duracion
				FROM capitulo
				WHERE id_temporada = $1
				ORDER BY nro;
				`
				rowsCaps, _ := db.BaseDeDatos.Query(consultaCapitulos, tempId)

				for rowsCaps.Next() {
					var cap Capitulo
					rowsCaps.Scan(&cap.Nro, &cap.Titulo, &cap.Duracion)
					temp.Capitulos = append(temp.Capitulos, cap)
				}
				rowsCaps.Close()

				c.Temporadas = append(c.Temporadas, temp)
			}
			rowsTemp.Close()
		}

		Contenidos = append(Contenidos, c)
	}

	respuestaJson, err := json.Marshal(Contenidos)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Write(respuestaJson)
}
func GetCategorias(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Max-Age", "15")
	w.WriteHeader(http.StatusOK)

	consulta := `SELECT id, descripcion FROM genero;`
	rows, errbd := db.BaseDeDatos.Query(consulta)
	if errbd != nil {
		fmt.Errorf("Error en la base de datos: %d", errbd)
		return
	}

	defer rows.Close()

	type Genero struct {
		Id          string `json:"id"`
		Descripcion string `json:"descripcion"`
	}

	var generos []Genero
	for rows.Next() {
		var g Genero
		rows.Scan(&g.Id, &g.Descripcion)
		generos = append(generos, g)
	}

	respuestaJson, err := json.Marshal(generos)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(respuestaJson)
}

func GetImagen(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Max-Age", "15")

	nombreImagen := ""
	fmt.Sscanf(r.URL.Query().Get("imagen"), "%s", &nombreImagen)

	tipo := ""
	fmt.Sscanf(r.URL.Query().Get("tipo"), "%s", &tipo)

	fileInfo, err := os.Stat("imagenes/" + nombreImagen)
	if err != nil {
		http.Error(w, "Imagen no encontrada", http.StatusNotFound)
		return
	}

	file, err := os.Open("imagenes/" + nombreImagen)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	data := make([]byte, fileInfo.Size())
	count, err := file.Read(data)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(count, nombreImagen)

	w.WriteHeader(http.StatusOK)

	if tipo == "blob" {
		w.Header().Set("Content-Type", "image/jpg; charset=utf-8")
		w.Write(data)
	} else {
		w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
		w.Write([]byte(data))
	}
}
