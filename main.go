package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

//funcion para conectar a la base de datos
func conexionBD() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "root"
	Contraseña := "0801"
	Nombre := "sistema"

	conexion, err := sql.Open(Driver, Usuario+":"+Contraseña+"@tcp(127.0.0.1)/"+Nombre)
	if err != nil {
		panic(err.Error())
	}
	return conexion
}

//este metodo sirve para obtener la informacion de los templates
var plantillas = template.Must(template.ParseGlob("plantillas/*"))

func main() {
	http.HandleFunc("/", Inicio)
	http.HandleFunc("/crear", Crear)
	http.HandleFunc("/insertar", Insertar)
	http.HandleFunc("/borrar", Borrar)
	http.HandleFunc("/editar", Editar)
	http.HandleFunc("/actualizar", Actualizar)

	log.Println("Servidor corriendo")
	//definir que puerto pueda correr el server
	http.ListenAndServe(":8080", nil)
}

func Borrar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")
	//fmt.Println(idEmpleado)

	conexionEstablecida := conexionBD()
	//ejecutar sentencia sql
	borrarRegristros, err := conexionEstablecida.Prepare("DELETE FROM empleados WHERE id=?")

	if err != nil {
		panic(err.Error())
	}
	borrarRegristros.Exec(idEmpleado)

	http.Redirect(w, r, "/", 301)
}

type Empleado struct {
	Id     int
	Nombre string
	Correo string
}

func Inicio(w http.ResponseWriter, r *http.Request) {

	conexionEstablecida := conexionBD()
	//ejecutar sentencia sql

	//insertarRegistros, err := conexionEstablecida.Prepare("INSERT INTO empleados(nombre,correo) VALUES('german', 'german@gmail.com')")
	registros, err := conexionEstablecida.Query("SELECT * FROM empleados")

	if err != nil {
		panic(err.Error())
	}
	//insertarRegistros.Exec()
	empleado := Empleado{}
	arregloEmpleado := []Empleado{}

	for registros.Next() {
		var id int
		var nombre, correo string
		err = registros.Scan(&id, &nombre, &correo)

		//if para verificar cualquier error
		if err != nil {
			panic(err.Error())
		}

		//usar la estructura creada arriba
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo

		//llenar arreglo con los datos
		arregloEmpleado = append(arregloEmpleado, empleado)
	}
	//fmt.Println(arregloEmpleado)

	//mostrar la plantilla html creada en otra carpeta
	plantillas.ExecuteTemplate(w, "inicio", arregloEmpleado)
}

func Editar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")
	fmt.Println(idEmpleado)

	conexionEstablecida := conexionBD()
	registro, err := conexionEstablecida.Query("SELECT * FROM empleados WHERE id=?", idEmpleado)

	empleado := Empleado{}
	for registro.Next() {
		var id int
		var nombre, correo string
		err = registro.Scan(&id, &nombre, &correo)

		//if para verificar cualquier error
		if err != nil {
			panic(err.Error())
		}

		//usar la estructura creada arriba
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo
	}
	fmt.Println(empleado)
	//mostrar la plantilla html creada en otra carpeta
	plantillas.ExecuteTemplate(w, "editar", empleado)

}

func Crear(w http.ResponseWriter, r *http.Request) {
	//mostrar la plantilla html creada en otra carpeta
	plantillas.ExecuteTemplate(w, "crear", nil)
}

func Insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		conexionEstablecida := conexionBD()
		//ejecutar sentencia sql
		insertarRegistros, err := conexionEstablecida.Prepare("INSERT INTO empleados(nombre,correo) VALUES(?,?)")

		if err != nil {
			panic(err.Error())
		}
		insertarRegistros.Exec(nombre, correo)

		http.Redirect(w, r, "/", 301)
	}
}

func Actualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		conexionEstablecida := conexionBD()
		//ejecutar sentencia sql
		ModificarRegistros, err := conexionEstablecida.Prepare("UPDATE empleados SET nombre=?, correo=? WHERE id=?")

		if err != nil {
			panic(err.Error())
		}
		ModificarRegistros.Exec(nombre, correo, id)

		http.Redirect(w, r, "/", 301)
	}
}
