package main

import (
	//"io/ioutil"
	"os"
	"strings"
	"log"
	"net/http"
	"bufio"
	"html/template"
	"fmt"
	"strconv"
	"crypto/rand"
	"html"
	"math/big"
	//"time"
)



func randInt(min,max int64) int {
	//rand.Seed(time.Now().UnixNano())
	byte, err := rand.Int(rand.Reader,big.NewInt(max))
	if err != nil{
		panic(err)
	}
	n := int(byte.Int64())
	return ((n - int(min)) + int(min))
}

func getHostname() string {
	hostname, err := os.Hostname()

	if err != nil {
		panic(err)
	}

	return hostname
}

func PWGen(min int, emoji bool) string {
	file,err := os.Open("wordlist.txt")
	emojilist := []int{128513, 128514, 128169, 128126, 9924, 127383, 128586, 127943}
	max := 1000
	if err != nil {
		fmt.Print(err)
	}
	var list []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		list = append(list, scanner.Text())
	}
	var password string = ""
	for{
		if len(password) > min || len(password) > max {
			break
		}
		if password == "" {
			password += strings.Title(list[randInt(0,int64(len(list)))]) //choose capitalized random word from list
		} else {
			if emoji {
				password += (html.UnescapeString("&#" + strconv.Itoa(emojilist[randInt(0, int64(len(emojilist)))]) +";"))
			} else {
				password += strconv.Itoa(randInt(10,99))
			}
			password += strings.Title(list[randInt(0,int64(len(list)))])//add number, and then capitalized random word
		}
	}
	return password
}

func main() {
	http.HandleFunc("/",func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("static/layout.html")
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		t.Execute(w, map[string] string{ "hostname": getHostname() })
	})
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		boolEmoji := false
		minimum := 12
		if minint,_ := strconv.Atoi(strings.Join(r.Form["minimum"],"")); minint != 0 {
			minimum = minint
		}
		if strings.Join(r.Form["emoji"],"") == "on" {
			boolEmoji = true
		}

		fmt.Fprintf(w,PWGen(minimum,boolEmoji))
	})
	log.Fatal(http.ListenAndServe(":8081",nil))
}
