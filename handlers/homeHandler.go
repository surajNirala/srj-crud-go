package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/surajNirala/srj-crud/config"
	"github.com/surajNirala/srj-crud/models"
)

type ViewData struct {
	Items []string
}

func Index(w http.ResponseWriter, r *http.Request) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	absPath := filepath.Join(cwd, "templates", "list.html")
	tmpl, err := template.New("list").ParseFiles(absPath)
	DB := config.DB
	var users []models.User
	DB.Find(&users)
	fmt.Println(users)
	data := struct {
		Users []models.User
	}{
		Users: users,
	}
	// data := ViewData{
	// 	Items: []string{"Item 1", "Item 2", "Item 3"},
	// }
	err = tmpl.ExecuteTemplate(w, "list.html", data)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// responses.ResponseSuccess(w, http.StatusOK, "Fetched User List.", users)
}
