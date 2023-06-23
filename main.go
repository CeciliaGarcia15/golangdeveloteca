package main

import (
	"log"
	"net/http"
	"text/template"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func conexionBD() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "root"
	Contrasenia := ""
	Nombre := "sistemagolang"

	conexion, err := sql.Open(Driver, Usuario+":"+Contrasenia+"@tcp(127.0.0.1)/"+Nombre)
	if err != nil {
		panic(err.Error())
	}
	return conexion
}

var plantillas = template.Must(template.ParseGlob("plantillas/*"))

func main() {
	//RUTAS
	http.HandleFunc("/", Inicio)
	http.HandleFunc("/crear", Crear)
	http.HandleFunc("/insertar", Insertar)
	http.HandleFunc("/borrar", Borrar)
	http.HandleFunc("/editar", Editar)
	http.HandleFunc("/actualizar", Actualizar)

	log.Println("Servidor corriendo...")
	http.ListenAndServe(":8080", nil)
}

func Borrar(w http.ResponseWriter, r *http.Request) {
	//se obtiene el id que viene por la ruta
	idEmpleado := r.URL.Query().Get("id")
	//se guarda dentro de conexionEstablecida la conexion a la db
	conexionEstablecida := conexionBD()
	//le decimos a conexion establecida que guarde la consulta sql que quiero dentro de una variable llamada insertarRegistrops, si ocurre un error entonces se guardara en la variable llamada err
	borrarRegistro, err := conexionEstablecida.Prepare("DELETE FROM empleados WHERE id=?")
	//si err tiene algo significa que ocurrio un error
	if err != nil {
		//se procede a mostrar el error
		panic(err.Error())
	}
	//ejecuta la consulta sql que guardamos en insertar registros
	//le paso las variables por aca,entonces cuando se ejecute va aremplazar los signos ? ? por lo que tenga guardado las variables
	borrarRegistro.Exec(idEmpleado)

	//que redirija hacia la barra y le pasamos el codigo, ese codigo es como el del 404 podes buscar en internet como son los otros codigos
	http.Redirect(w, r, "/", 301)

}

type Empleado struct {
	Id     int
	Nombre string
	Correo string
}

func Inicio(w http.ResponseWriter, r *http.Request) {
	//se guarda dentro de conexionEstablecida la conexion a la db
	conexionEstablecida := conexionBD()
	//le decimos a conexion establecida que guarde la consulta sql que quiero dentro de una variable lo que traiga el query cuando se ejecute, si ocurre un error entonces se guardara en la variable llamada err
	registros, err := conexionEstablecida.Query("SELECT * FROM empleados")
	//si err tiene algo significa que ocurrio un error
	if err != nil {
		//se procede a mostrar el error
		panic(err.Error())
	}
	empleado := Empleado{}
	arregloEmpleado := []Empleado{}

	for registros.Next() {
		var id int
		var nombre, correo string
		err = registros.Scan(&id, &nombre, &correo)
		if err != nil {
			panic(err.Error())
		}
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo
		arregloEmpleado = append(arregloEmpleado, empleado)
	}
	/*fmt.Println(arregloEmpleado)
	fmt.Fprintf(w, "Hola develoteca") */
	plantillas.ExecuteTemplate(w, "inicio", arregloEmpleado)
}

func Crear(w http.ResponseWriter, r *http.Request) {
	plantillas.ExecuteTemplate(w, "crear", nil)
}

func Insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		//se guarda dentro de conexionEstablecida la conexion a la db
		conexionEstablecida := conexionBD()
		//le decimos a conexion establecida que guarde la consulta sql que quiero dentro de una variable llamada insertarRegistrops, si ocurre un error entonces se guardara en la variable llamada err
		insertarRegistros, err := conexionEstablecida.Prepare("INSERT INTO empleados(nombre,correo) VALUES(?,?)")
		//si err tiene algo significa que ocurrio un error
		if err != nil {
			//se procede a mostrar el error
			panic(err.Error())
		}
		//ejecuta la consulta sql que guardamos en insertar registros
		//le paso las variables por aca,entonces cuando se ejecute va aremplazar los signos ? ? por lo que tenga guardado las variables
		insertarRegistros.Exec(nombre, correo)

		//que redirija hacia la barra y le pasamos el codigo, ese codigo es como el del 404 podes buscar en internet como son los otros codigos
		http.Redirect(w, r, "/", 301)
	}
}

func Editar(w http.ResponseWriter, r *http.Request) {
	//se obtiene el id que viene por la ruta
	idEmpleado := r.URL.Query().Get("id")
	//se guarda dentro de conexionEstablecida la conexion a la db
	conexionEstablecida := conexionBD()

	//le decimos a conexion establecida que guarde la consulta sql que quiero dentro de una variable lo que traiga el query cuando se ejecute, si ocurre un error entonces se guardara en la variable llamada err
	registro, err := conexionEstablecida.Query("SELECT * FROM empleados WHERE id=?", idEmpleado)
	empleado := Empleado{}
	for registro.Next() {
		var id int
		var nombre, correo string
		err = registro.Scan(&id, &nombre, &correo)
		if err != nil {
			panic(err.Error())
		}
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo
	}

	/* fmt.Println(empleado) */
	plantillas.ExecuteTemplate(w, "editar", empleado)
}

func Actualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		//se guarda dentro de conexionEstablecida la conexion a la db
		conexionEstablecida := conexionBD()
		//le decimos a conexion establecida que guarde la consulta sql que quiero dentro de una variable llamada insertarRegistrops, si ocurre un error entonces se guardara en la variable llamada err
		actualizarRegistro, err := conexionEstablecida.Prepare("UPDATE empleados SET nombre=?,correo=? WHERE id=?")
		//si err tiene algo significa que ocurrio un error
		if err != nil {
			//se procede a mostrar el error
			panic(err.Error())
		}
		//ejecuta la consulta sql que guardamos en actualizar registros
		//le paso las variables por aca,entonces cuando se ejecute va aremplazar los signos ? ? por lo que tenga guardado las variables
		actualizarRegistro.Exec(nombre, correo, id)

		//que redirija hacia la barra y le pasamos el codigo, ese codigo es como el del 404 podes buscar en internet como son los otros codigos
		http.Redirect(w, r, "/", 301)
	}
}
