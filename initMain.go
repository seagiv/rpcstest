package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var (
	flagURL      = flag.String("url", "http://127.0.0.1:8545/", "(option) url to test")
	flagNum      = flag.Int("n", 1, "(option) number requests to test")
	flagTimeout  = flag.Int("t", 100, "(option) time in milliseconds")
	flagRoutines = flag.Int("r", 1, "(option) use parallels routines")
	flagLogging  = flag.Bool("l", false, "(option) use logging")
)

var rMap map[int]*int

func isDone(thr int) (bool, int) {
	p := 0

	for _, v := range rMap {
		p = p + (*v)
	}

	/*fmt.Println("=========================")

	for k, v := range rMap {
		fmt.Printf("%d-%d\n", k, (*v))
	}

	fmt.Println("=========================")*/

	for _, v := range rMap {
		if thr > (*v) {
			return false, p
		}
	}

	return true, p
}

const hexh256 = "H256H"
const hexh64 = "H64H"
const hexh40 = "H40H"

var requests = []string{
	`{"jsonrpc":"2.0","id":1,"method":"eth_getTransactionCount","params":["0xH40H","latest"]}`,
	`{"jsonrpc":"2.0","method":"eth_getTransactionByHash","params":["0xH64H"],"id":1}`,
	`{"jsonrpc":"2.0","method":"eth_sendRawTransaction","params":["0xH256H"],"id":1}`,
	`{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}`,
}

func getRequest() []byte {
	format := requests[rand.Intn(len(requests))]

	if strings.Contains(format, hexh256) {
		buffer := make([]byte, 256)

		rand.Read(buffer)

		format = strings.Replace(format, hexh256, "%s", -1)

		return []byte(fmt.Sprintf(format, hex.EncodeToString(buffer)))
	}

	if strings.Contains(format, hexh64) {
		buffer := make([]byte, 32)

		rand.Read(buffer)

		format = strings.Replace(format, hexh64, "%s", -1)

		return []byte(fmt.Sprintf(format, hex.EncodeToString(buffer)))
	}

	if strings.Contains(format, hexh40) {
		buffer := make([]byte, 20)

		rand.Read(buffer)

		format = strings.Replace(format, hexh40, "%s", -1)

		return []byte(fmt.Sprintf(format, hex.EncodeToString(buffer)))
	}

	return []byte(format)
}

func doRequest(url string, jsonStr []byte, logging bool) error {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	//req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	if logging {
		fmt.Printf("response Body: %s\n", string(body))
	}

	return nil
}

func doRoutine(url string, t time.Duration, n int, progress *int, logging bool) {
	for i := 0; i < n; i++ {
		doRequest(
			url,
			getRequest(),
			logging,
		)

		*progress = i + 1

		time.Sleep(t)
	}
}

func doTest(url string, t time.Duration, n, r int, logging bool) error {
	rMap = make(map[int]*int, r)

	pr := int(n / r)

	for i := 0; i < r; i++ {
		p := 0

		rMap[i] = &p

		go doRoutine(url, t, pr, rMap[i], logging)
	}

	pp := 0

	for {
		time.Sleep(10 * time.Second)

		d, p := isDone(pr)

		fmt.Printf("progress: %d of %d %d r/s\n", p, n, (p-pp)/10)

		if d {
			return nil
		}

		pp = p
	}
}

func main() {
	flag.Parse()

	fmt.Printf("Testing: %s\n", *flagURL)

	err := doTest(*flagURL, time.Duration(*flagTimeout)*time.Millisecond, *flagNum, *flagRoutines, *flagLogging)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)

		return
	}

	fmt.Printf("Done.\n")
}

