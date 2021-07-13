package main

import (
    "fmt"
    "os"
    "net/http"
    "encoding/json"
    "path/filepath"
    "log"
    "io/ioutil"
   _ "github.com/lib/pq"
    "ab_log_bmx280/models"
    "ab_log_bmx280/configer"
     "time"
)



//config reader
func Config_reader(cfg_file string) configer.Configuration {

	file, err := os.Open(cfg_file)
	if err != nil {
		fmt.Println("can't open config file: ", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	Config := configer.Configuration{}
	err = decoder.Decode(&Config)
	if err != nil {
		fmt.Println("can't decode config JSON: ", err)
	}

	return Config
}




func inserter(press_val string) {
        dt := time.Now()
        var lastInsertId int
        err := models.Db.QueryRow("INSERT INTO pressure (p_val,p_date) VALUES($1,$2) returning p_id;", press_val, dt).Scan(&lastInsertId)
        checkErr(err)
        fmt.Println("last inserted id =", lastInsertId)

}



func getBmx280(host string )(press_val string){
    resp, err := http.Get(host) 
    if err != nil { 
        fmt.Println(err) 
        return
    } 
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
          fmt.Println(err)
          return
    }
    //fmt.Println(string(body))
    press_val = string(body)
    return press_val
}



func checkErr(err error) {
        if err != nil {
            panic(err)
        }
    }


func main() {

     version := "0.0.1"
     fmt.Println("ab-log inserter bmx280  version:"+version)
//************************* read config ******************************************//
     dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
     if err != nil {
            log.Fatal(err)
     }

     //log.Println(dir)
     cfg := Config_reader(filepath.Join(dir,"bmx280.conf"))
     models.Initdb(cfg)

//*********************** parse config **********************************//
   //logging
   log_dir := "./log"
   if _, err := os.Stat(log_dir); os.IsNotExist(err) {
		os.Mkdir(log_dir, 0644)
   }
   file, err := os.OpenFile(filepath.Join(log_dir,cfg.Log_file_name), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
   if err != nil {
		log.Fatal(err)
   }
   defer file.Close()
   log.SetOutput(file)
   log.Println("Logging to a file bmx280!")

   pressure :=getBmx280("http://192.168.71.74/sec/?pt=33&scl=34&i2c_dev=bmx280")
   //fmt.Println(pressure)
   inserter(pressure)

}
