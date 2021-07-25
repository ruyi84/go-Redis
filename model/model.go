package model

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

var M map[string]Value

func init() {
	M = make(map[string]Value)

	_, err := os.Stat("./saveData.json")
	if err == nil {
		saveData, err := ioutil.ReadFile("./saveData.json")
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(saveData, &M)
		if err != nil {
			panic(err)
		}
	}
}

type Value struct {
	Value      interface{}
	ExpireDate int
}

func Set(key string, value interface{}, expireDate int) {
	M[key] = Value{
		Value:      value,
		ExpireDate: expireDate,
	}
}

func Get(key string) interface{} {
	now := time.Now().UnixNano() / 1e6
	value := M[key]
	if value.ExpireDate == 0 || int64(value.ExpireDate) >= now {
		return value.Value
	}
	delete(M, key)
	return nil
}

func Del(key string) {
	delete(M, key)
}

func SaveToFile() {
	nowMap := M

	marshal, err := json.Marshal(nowMap)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("./saveData.json", marshal, 0644)
	if err != nil {
		panic(err)
	}
}
