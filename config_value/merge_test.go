package config_value

import (
	"encoding/json"
	"fmt"
	"github.com/lie-flat-planet/service-init-tool/config_source"
	"testing"
)

func TestMergeConfigValue(t *testing.T) {
	envVar := config_source.NewEnvVar()
	filevar := config_source.NewYamlFile("./test.yml")

	gotStructuralConfigInfo, err := MergeConfigValue(map[string]struct{}{
		"Mysql_Host":     {},
		"Mysql_Password": {},
		"A_B_C":          {},
	}, filevar, envVar)
	if err != nil {
		t.Fatal(err)
	}

	LogJSON(gotStructuralConfigInfo)

}

func LogJSON(v interface{}) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}
