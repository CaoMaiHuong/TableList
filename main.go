package main

import (
	"database/sql"
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
	Created   string   `json:"created"`
	Modified  string   `json:"modified"`
	Cve       []string `json:"cve"`
	Cvss_base float64  `json:"severity"`
	Qod       int      `json:"qod"`
	Xref      []string `json:"xref"`
	// Tag       NvtInfo  `json:"tag"`
}

type Task struct {
	Id           int            `json:"id"`
	Uuid         string         `json:"uuid"`
	Name         string         `json:"name"`
	Owner        string         `json:"owner"`
	Status       string         `json:"status"`
	ReportNumber sql.NullString `json:"rpnumber"`
	// Reports      string         `json:"report"`
	LastReport             string         `json:"last_report"`
	Severity               sql.NullString `json:"severity"`
	Comment                string         `json:"comment"`
	Target                 string         `json:"target"`
	Alert                  string         `json:"alert"`
	Schedule               sql.NullString `json:"schedule"`
	In_assets              string         `json:"in_assets"`
	Assets_apply_overrides string         `json:"assets_apply_overrides"`
	Assets_min_qod         string         `json:"assets_min_qod"`
	Alterable              int            `json:"alterable"`
	Auto_delete            string         `json:"auto_delete"`
	Auto_delete_data       string         `json:"auto_delete_data"`
	Scanner                int            `json:"scanner"`
	Config                 int            `json:"config"`
	Network                string         `json:"network"`
	Hosts_ordering         string         `json:"hosts_ordering"`
	Max_checks             string         `json:"max_checks"`
	Max_hosts              string         `json:"max_hosts"`
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
	// sort := vars["sort"]
	order := vars["order"]
	// sort = "asc"
	// order = "id"
	var offset int
	limit := 15
	if page == 1 {
		offset = 0
	} else {
		offset = (page - 1) * limit
	}
	rows, err := db.Raw("SELECT id, uuid, name, family, creation_time, modification_time, cve, cvss_base, qod FROM nvts ORDER BY ? ASC LIMIT ? OFFSET ?", order, limit, offset).Rows()
	fmt.Println("SELECT id, uuid, name, family, creation_time, modification_time, cve, cvss_base, qod FROM nvts ORDER BY ? LIMIT ? OFFSET ?", order, limit, offset)
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
		// c := " %"
		// fmt.Println(s + string(c));
		// nvt.Qod = nvt.Qod + string(c)
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
	// paginator.Offset = offset
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

func allTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []Task

	rows, err := db.Raw("SELECT t.id, t.uuid, t.name, t.run_status, t.target, t.config, t.schedule, t.scanner, t.hosts_ordering, t.alterable, tp.max_checks, tp.max_hosts, tp.in_assets, tp.assets_apply_overrides, tp.assets_min_qod, tp.auto_delete, tp.auto_delete_data, r.count as reports, r.max as date, re.severity FROM tasks t LEFT JOIN (select task, count(id), max(date) from reports group by task) as r ON t.id = r.task LEFT JOIN (select task, max(severity) as severity from results group by task) as re ON t.id = re.task LEFT JOIN(select distinct tp.task, mc.value max_checks, mh.value max_hosts, ia.value in_assets, ao.value assets_apply_overrides, mq.value assets_min_qod, ad.value auto_delete, dd.value auto_delete_data from task_preferences tp inner join (select task, value from task_preferences where name='max_checks') as mc on tp.task = mc.task inner join (select task, value from task_preferences where name='max_hosts') as mh on tp.task = mh.task inner join (select task, value from task_preferences where name='in_assets') as ia on tp.task = ia.task inner join (select task, value from task_preferences where name='assets_apply_overrides') as ao on tp.task = ao.task inner join (select task, value from task_preferences where name='assets_min_qod') as mq on tp.task = mq.task inner join (select task, value from task_preferences where name='auto_delete') as ad on tp.task = ad.task inner join (select task, value from task_preferences where name='auto_delete_data') as dd on tp.task = dd.task) as tp on t.id = tp.task WHERE hidden = 0").Rows()
	if err != nil {
		log.Print(err)
		return
	}

	for rows.Next() {
		var task Task
		var last sql.NullString
		err = rows.Scan(&task.Id, &task.Uuid, &task.Name, &task.Status, &task.Target, &task.Config, &task.Schedule, &task.Scanner, &task.Hosts_ordering, &task.Alterable, &task.Max_checks, &task.Max_hosts, &task.In_assets, &task.Assets_apply_overrides, &task.Assets_min_qod, &task.Auto_delete, &task.Auto_delete_data, &task.ReportNumber, &last, &task.Severity)
		if err != nil {
			log.Print(err)
			return
		}

		if last.Valid == false {
			last.String = ""
			task.LastReport = last.String
		} else {
			i, err := strconv.ParseInt(last.String, 10, 64)
			if err != nil {
				panic(err)
			}
			task.LastReport = time.Unix(i, 0).Format(time.RFC850)
		}
		switch task.Status {
		case "0":
			task.Status = "Delete Requested"
		case "1":
			task.Status = "Done"
		case "2":
			task.Status = "New"
		case "3":
			task.Status = "Requested"
		case "4":
			task.Status = "Running"
		case "10":
			task.Status = "Stop Requested"
		case "11":
			task.Status = "Ultimate Delete Requested"
		case "12":
			task.Status = "Stopped"
		case "13":
			task.Status = "Interrupted"
		case "14":
			task.Status = "Ultimate Delete Waiting"
		case "15":
			task.Status = "Stop Request Giveup"
		case "16":
			task.Status = "Deleted Waiting"
		case "17":
			task.Status = "Ultimate Delete Waiting"
		}

		tasks = append(tasks, task)
	}

	json.NewEncoder(w).Encode(tasks)
}

func main() {
	var err error
	db, err = gorm.Open("postgres", "host=112.137.129.225 user=cmhuong dbname=gvmd password=123456 sslmode=disable")
	if err != nil {
		fmt.Println(`Could not connect to db`)
		panic(err)
	}
	myRouter := mux.NewRouter()
	myRouter.HandleFunc("/tasks", allTasks).Methods("GET")
	myRouter.HandleFunc("/nvts/page={page}&_sort={sort}&_order={order}", allNvts).Methods("GET")

	log.Fatal(http.ListenAndServe(":8081", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "DELETE", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(myRouter)))
}
