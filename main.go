package main

import (
	"http"
	"fmt"
	"bufio"
	"os"
	"time"
	"flag"
	"strconv"
)
/*
 * Reasonable DELAY_TIME values, tested by using curl on the service:
 *     (Firehose, Peak)     8000 tweets pr. s (~18 mb/s):      1000000000/140000
 *     (Firehose)           3000 tweets pr. s (~6.9 mb/s):     1000000000/4250
 *     (Gardenhose)          300 tweets pr. s (~0.69 mb/s):    1000000000/350
 *     (Sprinkler)            30 tweets pr. s (~0.07 mb/s):    1000000000/30
 */
var delay int64 = 0

var fileHandle *os.File

func main() {
	flag.Int64Var(&delay, "delay", 0, "Delay per tweet in ns.")
	port := flag.Int("port", 3000, "The port to bind to.")
	flag.Usage = usage
	flag.Parse()
	
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(2)
	}
	
	var err os.Error
	fileHandle, err = os.Open(flag.Arg(0))
	if err != nil {
		fmt.Printf("Could not open file: %s\n", flag.Arg(0))
		os.Exit(2)
	}
	defer fileHandle.Close()
	fmt.Printf("Starting up server on port %d.\n", *port)
	if (delay > 0) {
		fmt.Printf("Delaying tweets by %d ns.\n", delay)
	}
	
	http.HandleFunc("/", sample)
	http.ListenAndServe(":"+strconv.Itoa(*port), nil)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [options] inputfile\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n");
	fmt.Fprintf(os.Stderr, "Reasonable DELAY_TIME values, tested by using curl on the service:\n");
	fmt.Fprintf(os.Stderr, "    (Firehose, Peak)     8000 tweets pr. s (~18 mb/s):      7142\n");
	fmt.Fprintf(os.Stderr, "    (Firehose)           3000 tweets pr. s (~6.9 mb/s):     235294\n");
	fmt.Fprintf(os.Stderr, "    (Gardenhose)          300 tweets pr. s (~0.69 mb/s):    2857142\n");
	fmt.Fprintf(os.Stderr, "    (Sprinkler)            30 tweets pr. s (~0.07 mb/s):    33333333\n");
}

func sample(w http.ResponseWriter, r *http.Request) {
	/*_, err := fileHandle.Seek(0, 0) // Seek to start of file on new request.
	if err != nil {
		fmt.Printf("Got seek err: %s\n", err)
		os.Exit(2)
	}*/
	
	var reader *bufio.Reader
	reader = bufio.NewReader(fileHandle)
	for {
		line, readErr := reader.ReadBytes('\n')
		if readErr != nil {
			if readErr != os.EOF  {
				fmt.Printf("Got readErr: %s\n", readErr)
			}
			return
		}
		fmt.Fprintf(w, "%s", line)
		if delay > 0 {
			time.Sleep(delay)
		}
	}
}