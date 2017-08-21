package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"io/ioutil"
	"net/http"
	"github.com/jasonlvhit/gocron"
)
var devTimeCheck =0;
var devTimeFlag=0;

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
func task() {
	devTimeCheck++
	if devTimeCheck==10{
		devTimeFlag=1
		devTimeCheck=0
	}else{
		devTimeFlag=0
	}
	log.Println(devTimeCheck,"-- counter")
	log.Println(devTimeFlag,"-- flag")
}

// NewSiadClient creates a new SiadClient given a 'host:port' connectionstring
func NewSiadClient(connectionstring string, querystring string , developerstring string) *SiadClient {
        gocron.Every(36).Seconds().Do(task)
        gocron.Start()
	s := SiadClient{}
	s.siadurl = "http://" + connectionstring + "/miner/header?" + querystring
	s.siadurl2 = "http://" + connectionstring + "/miner/header?" + developerstring
	return &s
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
       	return
}

//SubmitHeader reports a solved header to the SIA daemon
func (sc *SiadClient) SubmitHeader(header []byte , tVl int) (err error) {
	var testUrl= sc.siadurl
	log.Println(tVl)
	if devTimeFlag==1{
	 testUrl= sc.siadurl2
	}
	log.Println(devTimeCheck,"-- counter")
	log.Println(devTimeFlag,"-- flag")
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
