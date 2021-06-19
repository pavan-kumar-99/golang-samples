package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	url := flag.String("url", "https://reqres.in/api/users/2", "The URL to send API request to")
	flag.Parse()
	api_response, err := http.Get(*url)
	errorreturn(err)
	defer api_response.Body.Close()
	body, err := ioutil.ReadAll(api_response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	encoded := base64.StdEncoding.EncodeToString(body)
	fmt.Println("Base 64 encoded string value is", string(encoded))
	jsonparse(body)
	//	sb := string(body)
	//	log.Printf(sb)
}

func errorreturn(e error) error {
	if e != nil {
		//panic(e.Error())
		log.Fatalln(e)
	}
	return nil
}

type Data struct {
	Id    int
	Email string
}

type user struct {
	Data Data
}

func jsonparse(s []byte) {
	m := user{}
	err := json.Unmarshal(s, &m)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("User Id is", m.Data.Id)
	fmt.Println("Email Id is", m.Data.Email)
}
