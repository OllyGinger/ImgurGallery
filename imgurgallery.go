package main

import (
	"database/sql"
	"fmt"
	"github.com/bmizerany/pq"
	"html/template"
	"net/http"
	"os"
	"time"
)

type ImgRow struct {
	Hash       string
	Date       int
	DeleteHash string
	Orig       string
	Thumb      string
}

var galleryTemplate = template.Must(template.ParseFiles("templates/gallery.html"))

func main() {
	http.HandleFunc("/", webIndex)
	http.HandleFunc("/uploaded", webUpload)
	http.Handle("/js/", http.StripPrefix("/js", http.FileServer(http.Dir("js"))))

	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		fmt.Printf("Error listning: %s", err)
	}
}

func dbConnect() *sql.DB {
	connStr, err := pq.ParseURL(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Printf("Error connecting: %s", err)
	}

	return db
}

func dbStoreImage(db *sql.DB, hash string, delHash string, orig string, thumb string) {
	defer db.Close()

	_, err := db.Exec("INSERT INTO img (hash, date, deletehash, orig, thumb) VALUES ($1, $2, $3, $4, $5);", hash, time.Now().Unix(), delHash, orig, thumb)

	if err != nil {
		fmt.Printf("Error storing image. %s", err)
	}
}


func dbGalleryView(db *sql.DB) *[]ImgRow {
	defer db.Close()

	rows, err := db.Query("SELECT * FROM img ORDER BY date DESC;")
	if err != nil {
		fmt.Printf("Error getting images:. %s", err)
		return nil
	}

	var images []ImgRow
	for rows.Next() {
		var img ImgRow
		err := rows.Scan(&img.Hash, &img.Date, &img.DeleteHash, &img.Orig, &img.Thumb)
		if err != nil {
			fmt.Printf("Error getting images:. %s", err)
			return nil
		}

		images = append(images, img)
	}

	return &images
}

type galleryView struct {
	Images *[]ImgRow
	ImgurAPIKey string
}

func webIndex(w http.ResponseWriter, req *http.Request) {
	var gallery galleryView
	gallery.Images = dbGalleryView(dbConnect())
	gallery.ImgurAPIKey = os.Getenv( "IMGUR_API_KEY" )

	if err := galleryTemplate.Execute(w, gallery); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func webUpload(w http.ResponseWriter, req *http.Request) {
	hash := req.FormValue("hash")
	deleteHash := req.FormValue("deletehash")
	orig := req.FormValue("orig")
	thumb := req.FormValue("thumb")


	fmt.Fprintf( w, hash )
	dbStoreImage(dbConnect(), hash, deleteHash, orig, thumb)
}
