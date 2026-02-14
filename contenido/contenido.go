package contenido

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"pureflix-go/db"
	"pureflix-go/jwt"
	"pureflix-go/utils"
)

func GetContenido(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Max-Age", "15")

	idTipo := 0
	fmt.Sscanf(r.URL.Query().Get("tipo"), "%d", &idTipo)
	fmt.Println("idTipo", idTipo)

	idCategoria := 0
	fmt.Sscanf(r.URL.Query().Get("categoria"), "%d", &idCategoria)

	palabra := ""
	fmt.Sscanf(r.URL.Query().Get("palabra"), "%s", &palabra)

	favorito := "false"
	fmt.Sscanf(r.URL.Query().Get("favorito"), "%s", &favorito)

	token := ""
	fmt.Sscanf(r.URL.Query().Get("token"), "%s", &token)

	proximo := "false"
	fmt.Sscanf(r.URL.Query().Get("proximo"), "%s", &proximo)

	notificacion := "false"
	fmt.Sscanf(r.URL.Query().Get("notificacion"), "%s", &notificacion)

	var usuario string

	if favorito == "true" || notificacion == "true" {
		usuarioExtraido, errToken := jwt.GetUsernameFromToken(token)
		if errToken != nil {
			utils.DevolverError(w, http.StatusUnauthorized)
			return
		}
		usuario = usuarioExtraido
	}

	consulta := `
	SELECT DISTINCT 
		c.id,
		c.link_trailer,
		c.titulo,
		c.imagen,
		c.resumen,
		EXTRACT(YEAR FROM c.anio) AS anio,
		c.tipo_contenido,
		es_proximo
	FROM contenido c
	`

	filtros := ""
	tieneJoinFavorito := false
	tieneJoinNotificacion := false

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

	if favorito == "true" {
		consulta += `JOIN favorito f ON f.id_contenido = c.id
					JOIN usuario u ON u.id = f.id_usuario`
		filtros += fmt.Sprintf(" AND u.nombre_usuario = $1")
		tieneJoinFavorito = true
		//TODO: ver inyecciones
	}

	if notificacion == "true" {
		consulta += `JOIN notificacion n ON n.id_contenido = c.id
					JOIN usuario u ON u.id = n.id_usuario`
		filtros += fmt.Sprintf(" AND u.nombre_usuario = $1")
		tieneJoinNotificacion = true
	}

	if proximo == "true" {
		filtros += " AND c.es_proximo = true "
	} else {
		filtros += " AND c.es_proximo = false "
	}

	consulta += `
	WHERE 1=1
	` + filtros

	var rows *sql.Rows
	var errbd error

	if tieneJoinFavorito || tieneJoinNotificacion {
		rows, errbd = db.BaseDeDatos.Query(consulta, usuario)
	} else {
		rows, errbd = db.BaseDeDatos.Query(consulta)
	}

	if errbd != nil {
		utils.DevolverError(w, http.StatusInternalServerError)
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
		Es_proximo   bool        `json:"es_proximo"`
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
			&c.Es_proximo,
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

	w.WriteHeader(http.StatusOK)
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
		fmt.Println("Error al abrir imagen", nombreImagen, err)
		return
	}
	defer file.Close()
	data := make([]byte, fileInfo.Size())
	_, err = file.Read(data)
	if err != nil {
		log.Println("Error devolviendo imagen", nombreImagen, err)
	}

	w.WriteHeader(http.StatusOK)

	if tipo == "blob" {
		w.Header().Set("Content-Type", "image/jpg; charset=utf-8")
		fmt.Println("blob:", nombreImagen)
		w.Write(data)
	} else {
		w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
		w.Write([]byte(data))
	}
}

func MarcarFavorito(w http.ResponseWriter, r *http.Request) {
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

	_, id, errToken := jwt.GetUsernameAndIdFromToken(datos["token"])
	if errToken != nil {
		utils.DevolverError(w, http.StatusUnauthorized)
		return
	}

	consulta := `INSERT INTO favorito (id_usuario, id_contenido) VALUES ($1, $2)`
	_, errbd := db.BaseDeDatos.Exec(consulta, id, datos["id_contenido"])

	if errbd != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"mensaje":"Error al marcar como favorito"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"mensaje":"Agergado como favorito correctamente"}`))
}
func EliminarFavorito(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	idContenido := 0
	fmt.Sscanf(r.URL.Query().Get("id_contenido"), "%d", &idContenido)

	token := ""
	fmt.Sscanf(r.URL.Query().Get("token"), "%s", &token)

	usuario, errorUsuario := jwt.GetUsernameFromToken(token)
	if errorUsuario != nil {
		utils.DevolverError(w, http.StatusUnauthorized)
		return
	}

	consulta := `DELETE FROM favorito f
				USING usuario u
				WHERE u.id = f.id_usuario
				  AND u.nombre_usuario = $1
				  AND f.id_contenido = $2;
				`
	_, errbd := db.BaseDeDatos.Exec(consulta, usuario, idContenido)

	if errbd != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"mensaje":"Error"}`))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"mensaje":"Se quitó de favorito correctamente"}`))
}

func MarcarNotificacion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	datos, err := utils.RecibeDatosPost(r, nil)
	if err != 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"mensaje":"Error leyendo datos"}`))
		return
	}

	_, id, errToken := jwt.GetUsernameAndIdFromToken(datos["token"])
	if errToken != nil {
		utils.DevolverError(w, http.StatusUnauthorized)
		return
	}

	consulta := `INSERT INTO notificacion (id_usuario, id_contenido) VALUES ($1, $2)`
	_, errbd := db.BaseDeDatos.Exec(consulta, id, datos["id_contenido"])

	if errbd != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"mensaje":"Error al marcar notifiación"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"mensaje":"Se le notificará cuando el contenido se estrene"}`))
}
func EliminarNotificacion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	idContenido := 0
	fmt.Sscanf(r.URL.Query().Get("id_contenido"), "%d", &idContenido)

	token := ""
	fmt.Sscanf(r.URL.Query().Get("token"), "%s", &token)

	usuario, errorUsuario := jwt.GetUsernameFromToken(token)
	if errorUsuario != nil {
		utils.DevolverError(w, http.StatusUnauthorized)
		return
	}

	consulta := `DELETE FROM notificacion n
				USING usuario u
				WHERE u.id = n.id_usuario
				  AND u.nombre_usuario = $1
				  AND n.id_contenido = $2;
				`
	_, errbd := db.BaseDeDatos.Exec(consulta, usuario, idContenido)

	if errbd != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"mensaje":"Error"}`))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"mensaje":"Ya no se le notificará cuando el contenido se estrene"}`))
}
