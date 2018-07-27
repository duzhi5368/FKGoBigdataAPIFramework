package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"utils"
)

var (
	Config Configure
)

func InitConfig(path string) {
	// 使用默认config
	json.Unmarshal(defaultConfig, &Config)

	if path == "" {
		fmt.Printf("config path is empty, will use default config.\n")
	} else {
		err := load(path)
		if err != nil {
			fmt.Printf("load keys failed: %v\n", err)
			return
		}
	}

	fmt.Println(Config.DumpInfo())
}
func load(path string) error {
	if utils.PathExist(path) == false {
		return fmt.Errorf("%s not exist\n", path)
	}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf(" read file %s content error %v", path, err)
	}
	err = json.Unmarshal(content, &Config)
	if err != nil {
		return fmt.Errorf(" unmarshal file %s content error %v", path, err)
	}
	return nil
}
