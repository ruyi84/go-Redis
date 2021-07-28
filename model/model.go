package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"time"
)

type Disk map[string]Value

var disk Disk

func init() {
	disk = make(Disk)

	_, err := os.Stat("./saveData.json")
	if err == nil {
		saveData, err := ioutil.ReadFile("./saveData.json")
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(saveData, &disk)
		if err != nil {
			panic(err)
		}
	}
}

func GetDisk() Disk {
	return disk
}

type Value struct {
	Value      interface{}
	ExpireDate int
}

func (disk Disk) Set(key string, value interface{}, expireDate int) {
	disk[key] = Value{
		Value:      value,
		ExpireDate: expireDate,
	}
}

func (disk Disk) Get(key string) interface{} {
	now := time.Now().UnixNano() / 1e6
	value := disk[key]
	if value.ExpireDate == 0 || int64(value.ExpireDate) >= now {
		return value.Value
	}
	delete(disk, key)
	return nil
}

func (disk Disk) Del(key string) {
	delete(disk, key)
}

func (disk Disk) Keys(search string) []string {
	var keyList []string
	for key := range disk {
		matchString, err := regexp.MatchString(search, key)
		if err != nil || !matchString {
			fmt.Println(search, key)
			continue
		}
		keyList = append(keyList, key)
	}

	return keyList
}

func SaveToFile() {
	nowMap := disk

	marshal, err := json.Marshal(nowMap)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("./saveData.json", marshal, 0644)
	if err != nil {
		panic(err)
	}
}
