package main

import (
	"encoding/json"
	"fmt"
	"log"

	"math"
	"net/http"
	"os"
	"strconv"

	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
	// uuid "github.com/satori/go.uuid"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

var db *gorm.DB

type Nvt struct {
	Id        int      `json:"id"`
	Uuid      string   `json:"uuid"`
	Name      string   `json:"name"`
	Family    string   `json:"family"`
	Created   string   `json:"creation_time"`
	Modified  string   `json:"modification_time"`
	Cve       []string `json:"cve"`
	Cvss_base float64  `json:"cvss_base"`
	Qod       int      `json:"qod"`
	Xref      []string `json:"xref"`
	// Tag       NvtInfo  `json:"tag"`
}

type Paginator struct {
	Records     interface{} `json:"datas"`
	TotalRecord int         `json:"total"`
	TotalPage   int         `json:"total_pages"`
	Limit       int         `json:"per_page"`
	Page        int         `json:"current_page"`
	PrevPage    int         `json:"prev_page"`
	NextPage    int         `json:"next_page"`
}

func allNvts(w http.ResponseWriter, r *http.Request) {
	var nvts []Nvt
	vars := mux.Vars(r)
	page, err := strconv.Atoi(vars["page"])

	orders := vars["order"]
	var offset int
	limit := 15
	if page == 1 {
		offset = 0
	} else {
		offset = (page - 1) * limit
	}
	rows, err := db.Raw("SELECT id, uuid, name, family, creation_time, modification_time, cve, cvss_base, qod FROM nvts ORDER BY "+orders+" LIMIT  ? OFFSET ?", limit, offset).Rows()

	if err != nil {
		log.Print(err)
		return
	}

	for rows.Next() {
		var nvt Nvt
		var cve string
		err = rows.Scan(&nvt.Id, &nvt.Uuid, &nvt.Name, &nvt.Family, &nvt.Created, &nvt.Modified, &cve, &nvt.Cvss_base, &nvt.Qod)
		if err != nil {
			log.Print(err)
			return
		}
		i, err := strconv.ParseInt(nvt.Created, 10, 64)
		j, err := strconv.ParseInt(nvt.Modified, 10, 64)
		if err != nil {
			panic(err)
		}
		nvt.Created = time.Unix(i, 0).Format(time.RFC850)
		nvt.Modified = time.Unix(j, 0).Format(time.RFC850)
		nvt.Cve = strings.Split(cve, ",")
		for i := range nvt.Cve {
			nvt.Cve[i] = strings.TrimSpace(nvt.Cve[i])
			if nvt.Cve[i] == "NOCVE" {
				nvt.Cve[i] = ""
			}
		}
		nvts = append(nvts, nvt)
	}
	var count int
	row := db.Raw("Select count(*) from nvts").Row() // (*sql.Row)
	row.Scan(&count)
	paginator := Paging(page, limit, count, &nvts)
	json.NewEncoder(w).Encode(paginator)
}

func Paging(page int, limit int, count int, result interface{}) *Paginator {
	var paginator Paginator

	paginator.TotalRecord = count
	paginator.Records = result
	paginator.Page = page
	paginator.Limit = limit
	if page >= 1 {
		paginator.PrevPage = page - 1
	} else {
		paginator.PrevPage = page
	}

	if page == paginator.TotalPage {
		paginator.NextPage = page
	} else {
		paginator.NextPage = page + 1
	}

	paginator.TotalPage = int(math.Ceil(float64(paginator.TotalRecord) / float64(paginator.Limit)))
	return &paginator
}

func main() {
	var err error
	db, err = gorm.Open("postgres", "host=112.137.129.225 user=cmhuong dbname=gvmd password=123456 sslmode=disable")
	if err != nil {
		fmt.Println(`Could not connect to db`)
		panic(err)
	}
	myRouter := mux.NewRouter()
	myRouter.HandleFunc("/nvts/page={page}&_order={order}", allNvts).Methods("GET")

	log.Fatal(http.ListenAndServe(":8081", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "DELETE", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(myRouter)))
}
