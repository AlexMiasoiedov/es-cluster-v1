package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	var esDocUrls []string = []string{"http://localhost:9200/megacorp/employee/1", "http://localhost:9200/megacorp/employee/2"}

	fd, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE, 0755)

	if err != nil {
		log.Fatal(err)
	}

	defer fd.Close()

	for {
		Read(esDocUrls, fd)
		time.Sleep(100 * time.Millisecond)
	}
}

func Read(esDocUrls []string, fd *os.File) {
	ch := make(chan string)

	for _, esDocUrl := range esDocUrls {
		go fetch(esDocUrl, ch)
	}

	for range esDocUrls {
		_, err := fd.WriteString(<-ch)

		if err != nil {
			log.Fatal(err)
		}
	}
}

func fetch(docUrl string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(docUrl)

	if err != nil {
		ch <- fmt.Sprintf("while requesting %s: %v\n", docUrl, err)
		return
	}

	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()

	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v\n", docUrl, err)
		return
	}

	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs %7d %s\n", secs, nbytes, docUrl)
}
