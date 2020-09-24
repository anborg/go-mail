package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

//For test json file.
func readEftJSON(path string, myarray *EftInfos) error {
	jsonFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	//log.Println(byteValue)
	json.Unmarshal(byteValue, &myarray)
	//log.Println(myarray)
	return nil
}
