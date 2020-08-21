/**
    @Author : Loneyers
    @Date : 2020/8/16
    @FileName : main
**/

package main

import (
	"bufio"
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)
var (
	result []string
	q	string
	page	int
)
type Config struct {
	Fofa struct {
		Email        string   `yaml:"email"`
		Key      string   `yaml:"key"`
	}
}

func init() {
	flag.IntVar(&page, "page", 1, "page,default 1")
	flag.StringVar(&q, "q", "", "example:app=\"Solr\"")
}
func writeLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}
func fofa(q string,page int){
	config := new(Config)
	file,err:=ioutil.ReadFile("config.yaml")
	if err!=nil{
		log.Println(err)
	}
	err =yaml.Unmarshal(file,&config)
	if err!=nil{
		log.Println(err)
	}
	if config.Fofa.Email == "" ||config.Fofa.Key == ""{
		fmt.Println("email or key is empty.")
		os.Exit(0)
	}
	base64q := base64.StdEncoding.EncodeToString([]byte(q))
	url := fmt.Sprintf("https://fofa.so/api/v1/search/all?email=%s&key=%s&qbase64=%s&size=100&page=%d&full=true",config.Fofa.Email,config.Fofa.Key,base64q,page)
	resp,err := http.Get(url)
	if err!=nil{
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err!=nil{
		log.Fatal(err)
	}
	for _, v := range gjson.Get(string(body), `results.#.0`).Array() {
		result = append(result,v.String())
		fmt.Println(v.String())
	}
}
func getPwd() string{
	dir,err :=os.Getwd();
	if err!=nil{
		log.Fatal(err)
	}
	if runtime.GOOS =="windows"{
		return dir+"\\"
	}
	return dir+"/"
}
func main(){
	flag.Parse()
	if q==""{
		fmt.Println(`exmaple: ./fofa -q domain="exmaple.com" -page 2`)
		return
	}
	for i:=1;i<= page ;i++{
		fofa(q,i)
	}
	if err := writeLines(result, time.Now().Format("2006-01-02-15-04-05")+".txt"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("result file:",getPwd()+time.Now().Format("2006-01-02-15-04-05")+".txt")
}
