package main

//all made to work even if there is no db

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Employee struct {
	Emp_no     int       `json:"emp_no"`
	Birth_date time.Time `json:"birth_date"`
	First_name string    `json:"first_name"`
	Last_name  string    `json:"last_name"`
	Gender     string    `json:"gender"`
	Hire_date  time.Time `json:"hire_date"`
}
type DeleteRec struct {
	DeleteId int `json:"deleteid"`
}

type APIHandler struct {
	db   *sql.DB
	mu   sync.Mutex
	dbok bool
}

func (h *APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("content-type", "application/json")

	if r.Method == "GET" {
		employees := []Employee{
			{1, time.Now(), "Could not", "Connect", "M", time.Now()},
			{2, time.Now(), "Fake Data", "For Testing", "M", time.Now()},
		}

		if h.dbok {
			h.mu.Lock()
			defer h.mu.Unlock()

			// Execute the query
			results, err := h.db.Query("SELECT emp_no, birth_date, first_name, last_name, gender, hire_date FROM employees")
			if err != nil {
				log.Println("Whoops, probaly was not able to connect")
				log.Println(err.Error()) // proper error handling instead of panic in your app
			} else {
				employees = []Employee{}
				for results.Next() {
					var employee Employee
					// for each row, scan the result into our tag composite object
					err = results.Scan(&employee.Emp_no, &employee.Birth_date, &employee.First_name, &employee.Last_name, &employee.Gender, &employee.Hire_date)
					if err != nil {
						log.Println(err.Error()) // proper error handling instead of panic in your app
					} else {
						employees = append(employees, employee)
					}
				}
			}
		}

		b, err := json.Marshal(employees)
		if err != nil {
			fmt.Fprintf(w, "[]")
		} else {
			w.Write(b)
		}
	} else if r.Method == "POST" {
		b, err := ioutil.ReadAll(r.Body)
		if err == nil {
			emp := Employee{}
			fmt.Println(string(b))
			err = json.Unmarshal(b, &emp)
			if err == nil {
				fmt.Println("Adding ", emp)
				h.mu.Lock()
				defer h.mu.Unlock()
				//INSERT INTO `employees` VALUES (10001,'1953-09-02','Georgi','Facello','M','1986-06-26')
				_, err = h.db.Exec("INSERT INTO `employees` VALUES (?,?,?,?,?,?)", emp.Emp_no, emp.Birth_date, emp.First_name, emp.Last_name, emp.Gender, emp.Hire_date)
			}
		}

		if err != nil {
			fmt.Println("Somewhere an error occured", err)
		}
	} else if r.Method == "DELETE" {
		b, err := ioutil.ReadAll(r.Body)
		if err == nil {
			r := DeleteRec{}

			err = json.Unmarshal(b, &r)
			if err == nil {
				fmt.Println("deleting", r.DeleteId)
				h.mu.Lock()
				defer h.mu.Unlock()
				_, err = h.db.Exec("DELETE FROM employees WHERE emp_no = ?", r.DeleteId)
			}
		}

		if err != nil {
			fmt.Println("Somewhere an error occured", err)
		}
	}

}

type HTTPHandler struct {
	staticcontent string
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, h.staticcontent)
}

func initdb(db *sql.DB, filen *string) {
	fmt.Println("Reading ", *filen)
	content, err := ioutil.ReadFile(*filen)
	if err != nil {
		fmt.Println("Couldnt read file, doing nothing")
		fmt.Println(err)
	} else {
		statements := strings.Split(string(content), ";")
		for _, statement := range statements {
			ststep := strings.TrimSpace(statement)
			if ststep != "" {
				stssep := fmt.Sprintf("%s;", ststep)
				_, err = db.Exec(stssep)
				if err != nil {
					log.Println("Error on ", stssep)
					log.Println(err.Error())
					break
				} else {
					log.Println("Executed ", stssep)
				}
			}
		}

	}
}

func main() {
	initPtr := flag.Bool("init", false, "is this the one time init")
	initfilePtr := flag.String("initfile", "/usr/share/mysql-employees/initmysql.sql", "the file to init from")
	htmlPtr := flag.String("htmlfile", "/usr/share/mysql-employees/index.html", "the file to serve as main point")

	flag.Parse()

	server := os.Getenv("MYSQL_SERVER")
	username := os.Getenv("MYSQL_USERNAME")
	password := os.Getenv("MYSQL_PASSWORD")
	dbname := os.Getenv("MYSQL_DB")
	connstring := fmt.Sprintf("%s:%s@%s/?parseTime=true", username, password, server)

	if !*initPtr {
		connstring = fmt.Sprintf("%s:%s@%s/%s?parseTime=true", username, password, server, dbname)
	}

	db, err := sql.Open("mysql", connstring)
	dbok := true

	if err != nil {
		log.Println(err)
		dbok = false
	} else {
		db.SetConnMaxLifetime(time.Minute * 3)
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(10)
	}

	defer db.Close()

	if *initPtr {
		if dbok {
			fmt.Println("Init mode")
			initdb(db, initfilePtr)
		} else {
			panic("Was not able to connect, not going to continue")
		}
	} else {
		mux := http.NewServeMux()

		apihandler := &APIHandler{db, sync.Mutex{}, dbok}
		mux.Handle("/api", apihandler)

		staticcontent := html
		content, err := ioutil.ReadFile(*htmlPtr)
		if err == nil {
			staticcontent = string(content)
		} else {
			fmt.Println("Defaulting to built-in html file, error reading html file ", *htmlPtr)
		}

		htmlhandler := &HTTPHandler{staticcontent}
		mux.Handle("/", htmlhandler)
		log.Fatal(http.ListenAndServe(":8080", mux))
	}

}
