package main

import (
	"fmt"
	//	"log"
	"net/http"
	//	"os"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"
	"time"
	//	calendar "code.google.com/p/google-api-go-client/calendar/v3"
)

type timeS struct {
	dataTime string
}

type insertEvent struct {
	summary  string
	location string
	start    timeS
	end      timeS
}

const calendarID = "d2e8nb7tkmp21gbfl656vqh4j4@group.calendar.google.com"

func createEvent(startime int64, endtime int64, location string, title string) string {
	client, key_api := InitAuth()

	//	sTime := &timeS{dataTime: time_Int2Str(startime)}

	//	eTime := &timeS{dataTime: time_Int2Str(endtime)}

	//	jEvent := insertEvent{summary: title, location: location, start: *sTime, end: *eTime}

	sTime := map[string]interface{}{"dateTime": time_Int2Str(startime)}
	eTime := map[string]interface{}{"dateTime": time_Int2Str(endtime)}

	mapE := map[string]interface{}{"summary": title, "location": location, "start": sTime, "end": eTime}

	event, err := json.Marshal(mapE)
	if err != nil {
		panic(err)
		return "false"
	}
	fmt.Print(string(event))

	resp, err := client.Post("https://www.googleapis.com/calendar/v3/calendars/"+calendarID+"/events?key="+key_api, "application/json", bytes.NewBuffer(event))
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		fmt.Println("response Body:", string(body))
		return "false"
	}
	fmt.Println("response Body:", string(body))
	var f interface{}
	err = json.Unmarshal(body, &f)

	m := f.(map[string]interface{})
	for k, v := range m {
		if k == "id" {
			return v.(string)
		}

	}
	return "false"
}

func subscribeEvent(email string, eventid string) {
	client, key_api := InitAuth()
	body := searchEvent(eventid)
	var f interface{}
	err := json.Unmarshal(body, &f)

	url := "https://www.googleapis.com/calendar/v3/calendars/" + calendarID + "/events/" + eventid + "?key=" + key_api

	var st interface{}
	var en interface{}
	var loc interface{}
	var summ interface{}
	var attend interface{}

	m := f.(map[string]interface{})
	for k, v := range m {
		if k == "start" {
			st = v
		}
		if k == "end" {
			en = v
		}
		if k == "location" {
			loc = v
		}
		if k == "summary" {
			summ = v
		}
		if k == "attendees" {
			attend = v
		}

	}
	fmt.Print("attend: %v\n", attend)

	//	a := map[string]interface{}{"email": email}
	//	attend = append(attend, a)

	mapE := map[string]interface{}{"summary": summ, "location": loc, "start": st, "end": en, "attendees": attend}

	event, err := json.Marshal(mapE)
	if err != nil {
		panic(err)
	}
	fmt.Print("subscribeEvent**********" + string(event))

	request, err := http.NewRequest("PUT", url, bytes.NewBuffer(event))
	request.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	body, _ = ioutil.ReadAll(resp.Body)

	fmt.Print(string(body))

}

func searchEvent(eventid string) []byte {
	client, key_api := InitAuth()
	resp, err := client.Get("https://www.googleapis.com/calendar/v3/calendars/" + calendarID + "/events/" + eventid + "?key=" + key_api)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	return body

}

func time_Int2Str(nsec int64) string {
	timestamp := time.Unix(0, nsec)

	timestampStr := timestamp.Format(time.RFC3339)

	newStr := timestampStr

	if strings.ContainsAny(timestampStr, "Z") {
		str := strings.Split(timestampStr, "Z")
		newStr = str[0] + ".000-" + str[1]
	} else {
		index := strings.LastIndex(timestampStr, "-")
		str0 := timestampStr[0:index]
		str1 := timestampStr[index+1 : len(timestampStr)]
		newStr = str0 + ".000-" + str1
	}

	return newStr
}

func main() {
	createEvent(1416751000000000000, 1416751878000000000, "ll", "lll")
	subscribeEvent("a@b.com", "4cd4l9ok94d4bdm4m2r2tj75u0")

}
