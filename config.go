package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"
)

const defaultDir = ".route2cloud"

type conf struct {
	Doc string `json: doc`
	MainUser string `json:"Main_user"`
	MainPass string `json:"Main_password"`
}

func readConfigFiles(confFlag *string) (err error) {
	err = nil
	if len(*confFlag) < 1 {
		u, _ := user.Current()
		*confFlag = u.HomeDir + "/" + defaultDir
	}
	fmt.Println("Read config files in " + *confFlag + "/")

	files, err := ioutil.ReadDir(*confFlag)
	if err != nil {
		fmt.Println(err)
		os.Exit(1);
	}

	for _, file := range files {
		name := file.Name()
		if !strings.HasSuffix(name, ".conf") {
			continue
		}
		fmt.Println("Read config " + name)

		var config map[string]interface{}

		//		fmt.Println("FILE: " +  file.Name())
		raw, err := ioutil.ReadFile(*confFlag + "/" + name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err := json.Unmarshal(raw, &config); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		tmp, exist := config["address"]
		if exist {
			address = tmp.(string)
		}
		tmp, exist = config["main_username"]
		if exist {
			main_username = tmp.(string)
		}
		tmp, exist = config["main_password"]
		if exist {
			main_password = tmp.(string)
		}
		_, globalDebug = config["debug"] // Optional. Debug print

		tmp, exist = config["routes"]
		if exist {
			r := tmp.([]interface{})
			for _, v := range r {
				// var route map[string]interface{}
				route := v.(map[string]interface{})
				route["file"] = name
				fmt.Printf("CONF: %+v\n" , route)
				routes = append(routes, route)
				// routes = append(routes, v.(map[string]interface{}))
			}
		}
	}
	return //err
}

func readConfig() {
	// TODO: Support folder with many json files in.
	// Read input parameters
	confFlag := flag.String("conf", "", "Configuration directory, containing *.conf files. Default: ~/.route2cloud")
	flag.Parse()
	fmt.Println("Generate main config");

	readConfigFiles(confFlag)
}
