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
		// do something
		fmt.Fprint(w, "not implemented")
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

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<!doctype html>
<html lang="en">
	<head>
	<!-- Required meta tags -->
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

	<!-- Bootstrap CSS -->
	<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css" integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous">
	<title>Employees App</title>

	</head>
	<body>
		<div class="container-fluid h-100">
			<div class="row bg-veeam-light align-items-center">
				<div class="col-lg-1"></div>
				<div class="col-lg-8   rounded border border-light">
					<div class="m-4">
					<h1>Employees</h1>
					</div>
				</div>
			</div>
			<div class="row" id="employees" >
					<div class="col-lg-1"></div>
                    <div class="col-lg-8">
                      <table class="table">
                        <tr v-for="emp in employees"   >
						  <td><button class="btn btn-danger rounded-circle" v-bind:id="emp.emp_no" onclick="deleterec(this.id);">X</button></td>
						  <td>{{emp.emp_no}}</td>
                          <td>{{emp.first_name}}</td>
						  <td>{{emp.last_name}}</td>
						  <td>{{emp.birth_date}}</td>
						  <td>{{emp.gender}}</td>
                        </tr>
                      </table>
                    </div>
            </div>
		</div>
				
		<!-- Optional JavaScript -->
		<!-- jQuery first, then Popper.js, then Bootstrap JS -->
		<script src="https://code.jquery.com/jquery-3.5.1.min.js" integrity="sha384-ZvpUoO/+PpLXR1lu4jmpXWu80pZlYUAfxl5NsBMWOEPSjUn/6Z/hRTt8+pR6L4N2" crossorigin="anonymous"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.12.9/umd/popper.min.js" integrity="sha384-ApNbgh9B+Y1QKtv3Rn7W3mgPxhU9K/ScQsAP7hUibX39j7fakFPskvXusvfa0b4Q" crossorigin="anonymous"></script>
		<script src="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/js/bootstrap.min.js" integrity="sha384-JZR6Spejh4U02d8jOt6vLEHfe/JQGiRRSQQxSfFWpi1MquVdAyjUar5+76PVCmYl" crossorigin="anonymous"></script>
		<script src="https://cdn.jsdelivr.net/npm/vue@2"></script>
		<script>
			window.resturl = "/api"

			window.employees = new Vue({
				el: "#employees",
				data: {
					employees : [
					]   
				}
			})
			function deleterec(id) {
				fetch(window.resturl, {
					method: 'DELETE', // *GET, POST, PUT, DELETE, etc.
					mode: 'cors', // no-cors, *cors, same-origin
					cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
					credentials: 'same-origin', // include, *same-origin, omit
					redirect: 'follow', // manual, *follow, error
					referrerPolicy: 'no-referrer', // no-referrer, *no-referrer-when-downgrade, origin, origin-when-cross-origin, same-origin, strict-origin, strict-origin-when-cross-origin, unsafe-url
					body: JSON.stringify({deleteid:parseInt(id)})
				}).then(response => response).then(data => {
					console.log("deleted")
				}).catch((error) => {
					console.error('Error:', error);
				}).finally(()=> {
				})
			}

			function refresh() {
				console.log("Refreshing")
				fetch(window.resturl, {
					method: 'GET', // *GET, POST, PUT, DELETE, etc.
					mode: 'cors', // no-cors, *cors, same-origin
					cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
					credentials: 'same-origin', // include, *same-origin, omit
					redirect: 'follow', // manual, *follow, error
					referrerPolicy: 'no-referrer' // no-referrer, *no-referrer-when-downgrade, origin, origin-when-cross-origin, same-origin, strict-origin, strict-origin-when-cross-origin, unsafe-url
				}).then(response => response.json()).then(data => {
					window.employees.employees = []
					data.forEach(emp => {
						window.employees.employees.push(emp)
					})
					console.log("refreshed")
				}).catch((error) => {
					console.error('Error:', error);
				}).finally(()=> {
					setTimeout(refresh, 3000);
				})
			}
			refresh();


		</script>
	</body>
</html>
	`)
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
			stssep := fmt.Sprintf("%s;", statement)
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

func main() {
	initPtr := flag.Bool("init", false, "is this the one time init")
	initfilePtr := flag.String("initfile", "/usr/share/mysql-employees/initmysql.sql", "the file to init from")
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
		mux.HandleFunc("/", handler)
		log.Fatal(http.ListenAndServe(":8080", mux))
	}

}
