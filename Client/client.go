package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/theerfan/urlshortener/util"
	"io/ioutil"
	"net/http"
	"os"
	// "strconv"
	"strings"
)

const servAddr = "http://127.0.0.1:3333/shortener"

func sendPOST(client *http.Client, method string, url string) {
	request := util.ClientRequest{Method: method, Url: url}
	body, err := json.Marshal(&request)
	req, err := http.NewRequest("POST", servAddr, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println(err)
	}
    req.Header.Set("Content-Type", "application/json")
	finishRequest(client, req)
}

func sendGET(client *http.Client, url string) {
	// request := util.ClientRequest{Method: method, Url: url}
	req, err := http.NewRequest("POST", servAddr, bytes.NewBuffer([]byte(url)))
	if err != nil {
		fmt.Println(err)
	}
    req.Header.Set("Content-Type", "text/plain")
    finishRequest(client, req)
}

func finishRequest(client *http.Client, req *http.Request) {
	resp, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
    }
    defer resp.Body.Close()

    resBody, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("The Answer is:", string(resBody))
}

func getInput(reader *bufio.Reader) (string, string) {
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading from input: " + err.Error())
	}
	input = strings.Replace(input, " ", "", -1)
	inputSplit := strings.Split(input, ":")	
	return inputSplit[0], strings.TrimSpace(strings.Join(inputSplit[1:], ""))
}


//RunClient starts a new client that reads from input and connects to a server on port 3333
func main() {
	var method string
	var err error
	fmt.Println("wait unitl client starts..")
	client := &http.Client{}
	fmt.Println("Enter 'S:longurl' or 'L:shorturl")
	reader := bufio.NewReader(os.Stdin)
	function, url := getInput(reader)
	for function != "quit" {
		if function == "S" {
			fmt.Println("Enter desired method: hash/counter")
			method, err = reader.ReadString('\n')
			method = strings.TrimSpace(method)
			fmt.Println([]byte(method))
			if err != nil {
				fmt.Println(err)
			}
			go sendPOST(client, method, url)
		} else if function == "L" {
			go sendGET(client, url)
		} else {
			fmt.Println("Wrong function type.")
		}
		fmt.Println("Enter 'S: longurl' or 'L: shorturl")
		function, url = getInput(reader)
	} 
}