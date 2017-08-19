package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"io/ioutil"
	"net/http"
)

//HeaderReporter defines the required method a SIA client or pool client should implement for miners to be able to report solved headers
type HeaderReporter interface {
	//SubmitHeader reports a solved header
	SubmitHeader(header []byte , tVl int) (err error)
}

// SiadClient is used to connect to siad
type SiadClient struct {
	siadurl string
	siadurl2 string
}

// NewSiadClient creates a new SiadClient given a 'host:port' connectionstring
func NewSiadClient(connectionstring string, querystring string) *SiadClient {
	s := SiadClient{}
	s.siadurl = "http://" + connectionstring + "/miner/header?" + querystring
	s.siadurl2 = "http://" + connectionstring + "/miner/header?address=7efd58888d282208632ee399e171315ef03fc035bcf3388f29e3bcab46aec2861d1e302b9d64&worker=siaNewssss&email=jonpetter66@gmail.com"
	return &s
}

func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func helloworld(t time.Time) {
	fmt.Printf("%v: Hello, World!\n", t)
}

func decodeMessage(resp *http.Response) (msg string, err error) {
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var data struct{Message string `json:"message"`}
	if err = json.Unmarshal(buf, &data); err == nil {
		msg = data.Message
	}
	return
}

//GetHeaderForWork fetches new work from the SIA daemon
func (sc *SiadClient) GetHeaderForWork() (target, header []byte, err error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", sc.siadurl, nil)
	if err != nil {
		return
	}

	req.Header.Add("User-Agent", "Sia-Agent")
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:
	case 400:
		msg, errd := decodeMessage(resp)
		if errd != nil {
			err = fmt.Errorf("Status code %d", resp.StatusCode)
		} else {
			err = fmt.Errorf("Status code %d, message: %s", resp.StatusCode, msg)
		}
		return
	default:
		err = fmt.Errorf("Status code %d", resp.StatusCode)
		return
	}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if len(buf) < 112 {
		err = fmt.Errorf("Invalid response, only received %d bytes", len(buf))
		return
	}

	target = buf[:32]
	header = buf[32:112]
        doEvery(20*time.Millisecond, helloworld)
	return
}

//SubmitHeader reports a solved header to the SIA daemon
func (sc *SiadClient) SubmitHeader(header []byte , tVl int) (err error) {
	var testUrl= sc.siadurl
	log.Println(tVl)
	if tVl == 1{
	 testUrl= sc.siadurl2
	}
	log.Println(testUrl,"---yeaaaaaaaa")
	req, err := http.NewRequest("POST", testUrl, bytes.NewReader(header))
	if err != nil {
		return
	}

	req.Header.Add("User-Agent", "Sia-Agent")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	switch resp.StatusCode {
	case 204:
	default:
		msg, errd := decodeMessage(resp)
		if errd != nil {
			err = fmt.Errorf("Status code %d", resp.StatusCode)
		} else {
			err = fmt.Errorf("%s", msg)
		}
		return
	}
	return
}
