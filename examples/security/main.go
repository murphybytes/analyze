package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"

	"github.com/murphybytes/analyze/expression"
	"github.com/osquery/osquery-go"
)

func main() {
	var osQuerySockPath, configPath string

	flag.StringVar(&configPath, "config", "conf.json", "the path to the json configuration defining checks")
	flag.StringVar(&osQuerySockPath, "osquery", "/var/osquery/osquery.em", "the path to the osquery socket")
	flag.Parse()

	client, err := osquery.NewClient(osQuerySockPath, 10*time.Second)
	if err != nil {
		log.Fatalf("could not create osquery client %q", err)
	}
	defer client.Close()

	buff, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("could not read config file %q", err)
	}
	var config AnalysisConfig
	if err := json.NewDecoder(bytes.NewBuffer(buff)).Decode(&config); err != nil {
		log.Fatalf("could not decode config file %q", err)
	}

	datas, err := collect(client, config.Collectors)
	if err != nil {
		log.Fatalf("collection step failed %q", err)
	}

	for _, check := range config.Checks {
		log.Printf("check %s", check.Description)
		data, ok := datas[check.CollectorID]
		if !ok {
			log.Fatalf("no collector %q", check.CollectorID)
		}
		for _, condition := range check.Conditions {
			passed, err := expression.Evaluate(data, condition.Predicate)
			if err != nil {
				log.Fatalf("an error occured %q evaluating predicate %q", err, condition.Predicate)
			}
			if passed {
				log.Printf("%s: %s", condition.Type, condition.Message)
			}
		}

	}
}

type Collector struct {
	ID string `json:"id"`
	// Expression retrieve data from OSQuery
	Expression string `json:"expression"`
}

type Condition struct {
	Type      string `json:"type"`
	Predicate string `json:"predicate"`
	Message   string `json:"message"`
}

type Check struct {
	CollectorID string      `json:"collector-id"`
	Description string      `json:"description"`
	Conditions  []Condition `json:"conditions"`
}

type AnalysisConfig struct {
	Collectors []Collector `json:"collectors"`
	Checks     []Check     `json:"checks"`
}

type collectorMap map[string]interface{}

// convert []map[string]string to something analyze can handle []interface{} where interface{} is map[string]interface{}
func collect(client *osquery.ExtensionManagerClient, collectors []Collector) (collectorMap, error) {
	results := make(collectorMap)
	for _, coll := range collectors {
		rows, err := client.QueryRows(coll.Expression)
		if err != nil {
			return nil, err
		}
		var objects []interface{}
		for _, row := range rows {
			object := make(map[string]interface{})
			for k, v := range row {
				object[k] = v
			}
			objects = append(objects, object)

		}
		results[coll.ID] = objects

	}
	return results, nil
}
