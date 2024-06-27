package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/mattn/go-adodb"
)

type Data struct {
	PreviewPic string `json:"PreviewPic"`
	Desc       string `json:"Desc"`
}

func getData(w http.ResponseWriter, r *http.Request) {
	connStr := "Provider=Microsoft.ACE.OLEDB.12.0;Data Source=doorstyle.mdb"
	db, err := sql.Open("adodb", connStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT PreviewPic, [Desc] FROM doorshape")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var data []Data
	for rows.Next() {
		var d Data
		if err := rows.Scan(&d.PreviewPic, &d.Desc); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data = append(data, d)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func main() {
	http.HandleFunc("/api/data", getData)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
