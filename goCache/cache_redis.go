package goCache

import (
	"encoding/json"
	"fmt"
)

type DataPage struct {
	Name     string   `json:"name"`
	PageKeys []string `json:"pageKeys"`
}

func GetPageKey(name string) string {
	return fmt.Sprintf("PAG/%s", name)
}

func (dp *DataPage) AddPageKey(key string) {
	dp.PageKeys = append(dp.PageKeys, key)
}

/* Serialize and Deserialize DataPage */
func (dp *DataPage) ToJSON() interface{} {
	json, _ := json.Marshal(dp)
	return json
}

func (dp *DataPage) FromJSON(data interface{}) {
	if err := json.Unmarshal([]byte(data.(string)), dp); err != nil {
		fmt.Println("failed to unmarshal:", err)
	}
}
