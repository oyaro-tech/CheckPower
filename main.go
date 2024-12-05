package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"os"
	"time"
	"strconv"
	"fmt"
)

type HourlyData struct {
	Hour           string `json:"hour"`
	Electricity    int    `json:"electricity"`
	Description    string `json:"description"`
	PeriodLimitValue int `json:"periodLimitValue"`
}

type Today struct {
	EventDate     string       `json:"eventDate"`
	HoursList     []HourlyData `json:"hoursList"`
	ScheduleSince string       `json:"scheduleApprovedSince"`
}

type Graphs struct {
	Today Today `json:"today"`
}

type Current struct {
	Note     string `json:"note"`
	HasQueue string `json:"hasQueue"`
	Subqueue int    `json:"subqueue"`
	Queue    int    `json:"queue"`
}

type Response struct {
	Current Current `json:"current"`
	Graphs  Graphs  `json:"graphs"`
}

func main() {
	if (len(os.Args[1:]) < 1 || len(os.Args[1:]) > 2) {
		log.Fatalf("Usage: %s <accountNumber> <timeAhead?>\n", os.Args[0])
	}

	accountNumber, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Usage: %s <accountNumber> <timeAhead?>\n", os.Args[0])
	}

	timeAhead, err := time.ParseDuration("15m")
	if err != nil {
		log.Fatalf("[!] Error parsing duration: %v", err)
	}

	if (len(os.Args[1:]) == 2) {
		timeAhead, err = time.ParseDuration(os.Args[2])
		if err != nil {
			log.Fatalf("Usage: %s <accountNumber> <timeInMinutes?>\n", os.Args[0])
		}
	}

	// Prepare POST request
	url := "https://be-svitlo.oe.if.ua/GavGroupByAccountNumber"
	payload := []byte(fmt.Sprintf("accountNumber=%d&userSearchChoice=pob&address=", accountNumber))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		log.Fatalf("[!] Error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Origin", "https://svitlo.oe.if.ua")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.1.1 Safari/605.1.15")
	req.Header.Set("Referer", "https://svitlo.oe.if.ua/")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("[!] Error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("[!] Error reading response body: %v", err)
	}

	// Parse the response
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatalf("[!] Error parsing response: %v", err)
	}

	// Get the current time and add 15 minutes
	currentTime := time.Now()
	newTime := currentTime.Add(timeAhead)
	prev := 3

	// Check if the electricity is off during the specified hours
	for _, hourData := range response.Graphs.Today.HoursList {
		h := hourData.PeriodLimitValue - 1 // Shift to right hour

		if ((h == newTime.Hour()) && (prev == 0 || prev == 2)) {
			switch(hourData.Electricity) {
			case 0:
				break;
			case 1:
				cmd := exec.Command("osascript", "-e", `display notification "Очікується відключення електроенергії" with title "CheckPower"`)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				err := cmd.Run()
				if err != nil {
					log.Fatalf("Error running command: %v", err)
				}

				cmd = exec.Command("osascript", "-e", `say "Очікується відключення електроенергії" using "Lesya"`)
				err = cmd.Run()
				if err != nil {
					log.Fatalf("Error running command: %v", err)
				}

				os.Exit(0);
				break;
			case 2:
				if (prev == 2) {
					break;
				}

				cmd := exec.Command("osascript", "-e", `display notification "Очукіється можливе відключення електроенергії" with title "CheckPower"`)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				err := cmd.Run()
				if err != nil {
					log.Fatalf("Error running command: %v", err)
				}

				cmd = exec.Command("osascript", "-e", `say "Очукіється можливе відключення електроенергії" using "Lesya"`)
				err = cmd.Run()
				if err != nil {
					log.Fatalf("Error running command: %v", err)
				}

				os.Exit(0);
				break;
			default:
				log.Fatalf("Bad value of response")
				break;
			}
		}

		prev = hourData.Electricity
	}
}
