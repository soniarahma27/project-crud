package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type Department struct {
	id   int
	code string
	name string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "mysql"
	dbName := "dbcrud"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

// var tmpl = template.Must(template.ParseGlob("form/*"))
var tmpl *template.Template

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM department")
	if err != nil {
		panic(err.Error())
	}
	dept := Department{}
	res := []Department{}
	for selDB.Next() {
		var id int
		var code, name string
		err = selDB.Scan(&id, &code, &name)
		if err != nil {
			panic(err.Error())
		}
		dept.id = id
		dept.code = code
		dept.name = name
		res = append(res, dept)
	}
	// for _, v := range res {
	// 	log.Println(v)

	// }
	log.Println(res[0])
	log.Println(res[1])
	tmpl.ExecuteTemplate(w, "index.html", res)
	defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM department WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	dept := Department{}
	for selDB.Next() {
		var id int
		var code, name string
		err = selDB.Scan(&id, &code, &name)
		if err != nil {
			panic(err.Error())
		}
		dept.id = id
		dept.code = code
		dept.name = name
	}
	log.Println(dept)
	tmpl.ExecuteTemplate(w, "Show.html", dept)
	defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "Create.html", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM department WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	dept := Department{}
	for selDB.Next() {
		var id int
		var code, name string
		err = selDB.Scan(&code, &name)
		if err != nil {
			panic(err.Error())
		}
		dept.id = id
		dept.code = code
		dept.name = name
	}
	tmpl.ExecuteTemplate(w, "Edit.html", dept)
	defer db.Close()
}

func Create(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		code := r.FormValue("code")
		name := r.FormValue("name")
		insForm, err := db.Prepare("INSERT INTO department(code, name) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(code, name)
		log.Println("INSERT: Code: " + code + " | Name: " + name)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		code := r.FormValue("code")
		name := r.FormValue("name")
		id := r.FormValue("uid")
		insForm, err := db.Prepare("UPDATE department SET code=?, name=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(code, name, id)
		log.Println("UPDATE: Code: " + code + " | Name: " + name)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	dept := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM department WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(dept)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func main() {
	log.Println("Server started on: http://localhost:8080")
	var err error
	tmpl, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Println(err)
	}

	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/create", Create)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	http.ListenAndServe(":8080", nil)
}
