// Copyright 2022 Praetorian Security, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package runner

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/praetorian-inc/fingerprintx/pkg/plugins"
)

type outputFormat string

const (
	JSON    outputFormat = "JSON"
	CSV     outputFormat = "CSV"
	DEFAULT outputFormat = "DEFAULT"
)

func Report(services []plugins.Service) ([]string, error) {
	var writeFile *os.File
	//mysql://127.0.0.1:3306  default
	var outputFormat = JSON
	var csvWriter *csv.Writer
	var err error

	log.SetFlags(0)

	if len(config.outputFile) > 0 {
		var fileErr error
		writeFile, fileErr = os.Create(config.outputFile)
		if fileErr != nil {
			return nil, fileErr
		}
		log.SetOutput(writeFile)
	} else {
		log.SetOutput(os.Stdout)
	}
	defer writeFile.Close()

	if config.outputJSON {
		outputFormat = JSON
	} else if config.outputCSV {
		outputFormat = CSV
		csvWriter = csv.NewWriter(writeFile)
		if config.showErrors {
			err = csvWriter.Write([]string{"Host", "Port", "Service", "Metadata", "Error"})
		} else {
			err = csvWriter.Write([]string{"Host", "Port", "Service", "Data"})
		}
		if err != nil {
			return nil, err
		}
	}
	datas := make([]string, 0)
	for _, service := range services {

		switch outputFormat {
		case JSON:
			data, jerr := json.Marshal(service)
			datas = append(datas, string(data))
			if jerr != nil {
				return nil, err
			}
			//return string(data), nil
			//log.Println(string(data))
		case CSV:
			portStr := strconv.FormatInt(int64(service.Port), 10)
			err = csvWriter.Write([]string{service.Host, service.IP, portStr, service.Protocol,
				strconv.FormatBool(service.TLS), string(service.Raw)})
			if err != nil {
				return nil, err
			}
			csvWriter.Flush()
		default:
			if len(service.Host) > 0 {
				log.Printf("%s://%s:%d (%s)\n", strings.ToLower(service.Protocol), service.Host, service.Port, service.IP)
			} else {
				log.Printf("%s://%s:%d\n", strings.ToLower(service.Protocol), service.IP, service.Port)
			}
		}
	}

	return datas, nil
}
