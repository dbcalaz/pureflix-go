package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"pureflix-go/contenido"
	"pureflix-go/db"
	"pureflix-go/login"
	"pureflix-go/usuario"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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
	router.HandleFunc("/registrarNuevoUsuario", usuario.RegistrarNuevoUsuario)
	router.HandleFunc("/activar", usuario.ActivarCuenta)

	router.HandleFunc("/login", login.Login)
	router.HandleFunc("/validarToken", login.ValidarToken)

	router.HandleFunc("/recuperarContrasena", usuario.RecuperarContrasena)
	router.HandleFunc("/actualizarContrasena", usuario.ActualizarContrasena)
	router.HandleFunc("/actualizarFotoPerfil", usuario.ActualizarFotoPerfil)

	router.HandleFunc("/getDatosUsuario", usuario.GetDatosUsuario)
	router.HandleFunc("/actualizarUsuario", usuario.ActualizarUsuario)

	router.HandleFunc("/getContenido", contenido.GetContenido)
	router.HandleFunc("/getCategorias", contenido.GetCategorias)
	router.HandleFunc("/getImagen", contenido.GetImagen)

	router.HandleFunc("/marcarFavorito", contenido.MarcarFavorito)
	router.HandleFunc("/eliminarFavorito", contenido.EliminarFavorito)

	router.HandleFunc("/marcarNotificacion", contenido.MarcarNotificacion)
	router.HandleFunc("/eliminarNotificacion", contenido.EliminarNotificacion)

	fmt.Println("Iniciando servidor")

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "9000"
	}

	log.Fatal(http.ListenAndServe(":"+port, router))
}
