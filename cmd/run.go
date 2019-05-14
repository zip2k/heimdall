package cmd

import (
	"encoding/json"
	"fmt"
	"git01.bravofly.com/N7/heimdall/pkg/kubernetes"
	"git01.bravofly.com/N7/heimdall/pkg/logging"
	"git01.bravofly.com/N7/heimdall/pkg/model"
	"git01.bravofly.com/N7/heimdall/pkg/scheduler"
	"io/ioutil"
	"os"
)

func Run() {
	log.Init()
	kubernetes.Readiness()
	kubernetes.Liveness()

	config := readConfig(os.Getenv("CONFIG_PATH"))

	scheduler.Scheduler{
		Config: config,
	}.Start(Orchestration())

}

func readConfig(filePath string) *model.Config {
	file, err := os.Open(filePath)
	if err != nil {
		log.Error(fmt.Sprintf("error reading configuration. %v", err), nil)
		return model.DefautConfig()
	}
	defer file.Close()
	byteValue, _ := ioutil.ReadAll(file)
	var config *model.Config
	json.Unmarshal([]byte(byteValue), &config)
	return config
}
