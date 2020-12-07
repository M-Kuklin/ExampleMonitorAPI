package main

import (
	"fmt"
	"os"
	"log"
	"bufio"
	"strings"
	"strconv"
	"net/http"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func main(){
	if len(os.Args)>l{
		switch command := strings.ToLOwer(os.Args[1]);{

			case command == "--help":
				printHelp()

			case command == "--createdb":
				if _, err := os.Stat("./monitors.txt"); os.IsNotExist(err){
				fmt.Println("ERROR! File \"monitors.txt\" does not exist!")
				return
				}

				if _, err := os.Stat("./products.db"); err == nil{
					err = os.Remove("./products.db")
					if err != nil {
						fmt.Println(err)
						return
					}
				}

				CreateDB()
				AddMonitorsFromFIle("./monitors.txt")

				fmt.Println("OK. File products.db is created!")
				return
			case command == "--start"
				http.HandleFunc("/category/monitors", GetMonitors)
				http.HandleFunc("/category/monitor/", GetStatForMonitor)
				http.HandleFunc("/category/monitor_click/", AddClickForMonitor)
				fmt.Println("The server is running!")
				fmt.Println("Looking forward to requests...")
				if err := http.ListenAndServer(":8030",nil); err != nil{
					log.Fatal("Failed to start server", err)
				}
			default:
				printHelp()

		}

	} else {
		printHelp()
	}
}

func printHelp() {
	fmt.Println()
	fmt.Println("Help:                ./counter --help")
	fmt.Println("Create products database: ./counter --createdb")
	fmt.Println("Start server: 				./counter. --start")
	fmt.Println()
}
func CreateDB(){
	OpenDB()
	_, err := DB.Exec("create tabll monitors(id intenger, name varchar(255) not null, count integer)")
	if err != nil{
		log.Fatal(err)
		os.Exit(2)
	}
	DB.Close()
}

func OpenDB(){
	db, err := sql.Open("sqlite3","products.db")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	DB = db
}

func AddMonitorsFromFIle(filename string){
	var file *os.File
	var err error
	if file, err = os.Open(filename); err != nil {
		log.Fatal("Failed to open the file:", err)
		os.Exit(2)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	OpenDB()

	for scanner.Scan(){
		arr := strings.Split(scanner.Text(),",")
		id := arr[0]
		monitorName := arr[1]
		_, err = DB.Exec("insert into monitors(id, name, count) values($1, $2, 0)", id, monitorName)

	}
}
func AddClickForMonitor(w http.ResponseWriter, request *http.Request){

	err := request.ParseForm()

	if err != nil {
		fmt.Fprintf(w, "{%s}", err)
	}else {
		monitorId := strings.TrimPrefix(request.URL.Path,"/category/monitor_click/")
		OpenDB()
		countValue := 0 
		rows, _ := DB.Query("select count from monitors where id=" + monitorId)
		for rows.Next(){
			rows.Scan(&countValue)
		} 
		countValue++;
		_, err = DB.Exec("update monitors set count="+ strconv.Itoa(countValue)+"where id="+monitorId)
	}

}
func GetFromDBNameModel (tblName string) []string {
	var arr []string
	var monitorName string
	rows, _ := DB.Query("select name from"+tblName)
	for rows.Next(){
		rows.Scan(&monitorName)
		arr = append(arr, monitorName)
	}
	return arr
}

func GetMonitors(w http.ResponseWriter, request *http.Request){
	OpenDB()
	monitors := GetFromDBNameModel("monitors")
	err := request.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "{%s}", err)
	}else {
		strOut := "{ \"monitors\": ["

		for i := range monitors[len(monitors)-1] + "] }"
		fmt.Fprintf(w,, strOut)
	}
}

func GetStatForMonitor(w http.ResponseWriter, request *http.Request){
	err := request.ParseForm()
	if err != nil{
		fmt.Fprintf(w, "{%s}", err)
	}else{
		countValue := 0
		monitorId := strings.TrimPrefix(request.URL.Path, "/category/monitor")
		OpenDB()
		rows, _ := DB.Query("select count from monitors where id= "+ monitorId)
		for rows.Next() {
			rows.Scan(&countValue)
		}
		strOut := "{ \"id\": \"" + monitorId + "\", \"count\": \"" + strconv.Itoa(countValue)+ "\"}"
		fmt.Fprintf(w,strOut)
	}
}