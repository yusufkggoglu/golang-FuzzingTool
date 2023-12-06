package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/mbndr/figlet4go"
)

func worker(jobs <-chan string, wg *sync.WaitGroup) {
	for n := range jobs {
		res, err := http.Get(n)
		if err != nil {
			log.Fatal(err)
			return
		}
		if res.StatusCode == 200 {
			fmt.Println(n)
		}
		wg.Done()
	}
}

func main() {

	ascii := figlet4go.NewAsciiRender()
	options := figlet4go.NewRenderOptions()
	options.FontName = "larry3d"
	ascii.LoadFont("/path/to/fonts/")
	renderStr, _ := ascii.RenderOpts("Fuzzing Tools", options)
	fmt.Print(renderStr)
	fmt.Println("Developed by Yusuf Küçükgökgözoğlu\n\n")
	time.Sleep(3 * time.Second)

	urlPtr := flag.String("u", "https://stackoverflow.com", "Define url to scan ( e.g. http://localhost)")
	speedPtr := flag.Int("s", 72, "Scanning speed")
	txtPtr := flag.String("txt", "wordlist3.txt", "Path of txt which will be using")
	flag.Parse()
	if *urlPtr == "" || *speedPtr == 0 || *txtPtr == "" {
		flag.PrintDefaults()
		return
	}

	if *speedPtr > 100 {
		*speedPtr = 100
	}

	var wg sync.WaitGroup

	file, err := os.Open(*txtPtr)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	defer wg.Wait()

	var words []string

	Scanner := bufio.NewScanner(file)
	Scanner.Split(bufio.ScanWords)

	for Scanner.Scan() {
		words = append(words, Scanner.Text())
	}

	if err := Scanner.Err(); err != nil {
		log.Fatal(err)
	}

	jobs := make(chan string)

	for i := 1; i <= *speedPtr; i++ {
		go worker(jobs, &wg)
	}

	fmt.Println("Preparing ... \n")
	time.Sleep(time.Second)
	fmt.Println("Results :")
	for _, word := range words {
		wg.Add(1)
		url := *urlPtr + "/" + word
		jobs <- url
	}
	close(jobs)
}
