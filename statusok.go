package main

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
	"github.com/vemonet/statusok/database"
	"github.com/vemonet/statusok/notify"
	"github.com/vemonet/statusok/requests"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

type configParser struct {
	NotifyWhen    NotifyWhen               `json:"notifyWhen"`
	Requests      []requests.RequestConfig `json:"requests"`
	Notifications notify.NotificationTypes `json:"notifications"`
	Database      database.DatabaseTypes   `json:"database"`
	Concurrency   int                      `json:"concurrency"`
	Port          int                      `json:"port"`
}

type NotifyWhen struct {
	MeanResponseCount int `json:"meanResponseCount"`
	ErrorCount        int `json:"errorCount"`
}

func main() {

	//Cli tool setup to get config file path from parameters
	app := cli.NewApp()
	app.Name = "StatusOk"
	app.Usage = "Monitor your website.Get notifications when its down"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Value: "config.json",
			Usage: "location of config file",
		},
		cli.StringFlag{
			Name:  "log",
			Value: "",
			Usage: "file to save logs",
		},
	}

	app.Action = func(c *cli.Context) {

		if fileExists(c.String("config")) {

			if len(c.String("log")) != 0 {
				//log parameter given.Check if file can be created at given path

				if !logFilePathValid(c.String("log")) {
					println("Invalid File Path given for parameter --log")
					os.Exit(3)
				}
			}

			println("Reading File :", c.String("config"))

			//Start monitoring when a valid file path is given
			startMonitoring(c.String("config"), c.String("log"))
		} else {
			println("Config file not present at the given location: ", c.String("config"), "\nPlease give correct file location using --config parameter")
		}

	}

	//Run as cli app
	app.Run(os.Args)
}

func startMonitoring(configFileName string, logFileName string) {

	configFile, err := os.Open(configFileName)

	if err != nil {
		fmt.Println("Error opening config file:\n", err.Error())
	}

	//parse the config file data to configParser struct
	jsonParser := json.NewDecoder(configFile)
	var config configParser
	if err = jsonParser.Decode(&config); err != nil {
		fmt.Println("Error parsing config file .Please check format of the file \nParse Error:", err.Error())
		os.Exit(3)
	}

	//setup different notification clients
	notify.AddNew(config.Notifications)
	//Send test notifications to all the notification clients
	notify.SendTestNotification()

	//Create unique ids for each request date given in config file
	reqs, ids := validateAndCreateIdsForRequests(config.Requests)

	//Set up and initialize databases
	database.AddNew(config.Database)
	database.Initialize(ids, config.NotifyWhen.MeanResponseCount, config.NotifyWhen.ErrorCount)

	for {
		initSuccess, requestConfig := requests.RequestsInit(reqs, config.Concurrency)
		if initSuccess {
			break
		} else {
			// notify.SendErrorNotification(notify.ErrorNotification{
			// 	Url: requestConfig.Url,
			// 	RequestType: requestConfig.RequestType,
			// 	ResponseBody: "Error, one of the given URL could not be reached.",
			// 	Error: "The website is probably down, make sure all the websites you are monitoring are up when starting StatusOK",
			// 	OtherInfo: "Trying again in 1 minute"})
			println("\nTrying again in 1 minute, issue reaching to " + requestConfig.Url)
			time.Sleep(60 * time.Second)
		}
	}

	//Initialize and start monitoring all the apis
	requests.StartMonitoring()

	database.EnableLogging(logFileName)

	//Just to check StatusOk is running or not
	http.HandleFunc("/", statusHandler)

	if config.Port == 0 {
		//Default port
		http.ListenAndServe(":7321", nil)
	} else {
		//if port is mentioned in config file
		http.ListenAndServe(":"+strconv.Itoa(config.Port), nil)
	}
}

//Currently just tells status ok is running
//Planning to display useful information in future
func statusHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "StatusOk is running \n Planning to display useful information in further releases")
}

//Tells whether a file exits or not
func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func logFilePathValid(name string) bool {
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()
	if err != nil {
		return false
	}

	return true
}

//checks whether each request in config file has valid data
//Creates unique ids for each request using math/rand
func validateAndCreateIdsForRequests(reqs []requests.RequestConfig) ([]requests.RequestConfig, map[int]int64) {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	//an array of ids used by database pacakge to calculate mean response time and send notifications
	ids := make(map[int]int64, 0)

	//an array of new requests data after updating the ids
	newreqs := make([]requests.RequestConfig, 0)

	for i, requestConfig := range reqs {
		validateErr := requestConfig.Validate()
		if validateErr != nil {
			println("\nInvalid Request data in config file for Request #", i, " ", requestConfig.Url)
			println("Error:", validateErr.Error())
			os.Exit(3)
		}

		//Set a random value as id
		randInt := random.Intn(1000000)
		ids[randInt] = requestConfig.ResponseTime
		requestConfig.SetId(randInt)
		newreqs = append(newreqs, requestConfig)
	}

	return newreqs, ids
}
