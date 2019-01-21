package cmd

import (
	"encoding/json"
	"fmt"
	"git01.bravofly.com/golang/appfw/pkg/http/mgmt"
	"git01.bravofly.com/n7/heimdall/cmd/client"
	"git01.bravofly.com/n7/heimdall/cmd/data_collector"
	"git01.bravofly.com/n7/heimdall/cmd/metric"
	"git01.bravofly.com/n7/heimdall/cmd/model"
	"gopkg.in/robfig/cron.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
)

var logger = log.New(os.Stdout, "[HEIMDALL] ", log.LstdFlags)

func Run(filePath string) {
	mgmt.ConfigureLiveness(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		return
	})
	mgmt.ConfigureReadiness(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		return
	})

	//	config := readConfig(filePath)
	config := readConfig("/appfw/config/config.json")

	cronExpression := fmt.Sprintf("*/%s * * * *", config.CollectEveryMinutes)

	logger.Printf("start collecting data at every %sth minute", config.CollectEveryMinutes)

	c := cron.New()
	c.AddFunc(cronExpression, orchestrator(config))

	go c.Start()
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	s := <-sig
	c.Stop()

	fmt.Println("Got signal:", s)

	//orchestrator(config)
}

func readConfig(filePath string) *model.Config {
	file, err := os.Open(filePath)
	if err != nil {
		logger.Printf("error reading configuration. %v", err)
		return model.DefautConfig()
	}
	defer file.Close()
	byteValue, _ := ioutil.ReadAll(file)
	var config *model.Config
	json.Unmarshal([]byte(byteValue), &config)
	return config
}

func orchestrator(config *model.Config) func() {
	return func() {
		aggregate := dataCollector(config)
		metric.PushMetrics(aggregate, config)
	}
}

func dataCollector(config *model.Config) []*model.Aggregate {
	aggregate, _ := client.GetZonesId()
	aggregate, _ = data_collector.GetColocationTotals(aggregate, config)
	//aggregate, _ = data_collector.GetWafTotals(aggregate)

	return aggregate
}
