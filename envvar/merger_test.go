package envvar

import (
	"github.com/lie-flat-planet/service-init-tool/config_source"
	"github.com/lie-flat-planet/service-init-tool/util"
	"testing"
)

func TestMerger_Action(t *testing.T) {
	envVar := config_source.NewEnvVar()
	filevar := config_source.NewYamlFile("./merger_test.yml")

	m := NewMerger(map[string]struct{}{
		"Mysql_Host":     {},
		"Mysql_Password": {},
		"A_B_C":          {},
	}, envVar, filevar)

	gotStructuralConfigInfo, err := m.Action()
	if err != nil {
		t.Fatal(err)
	}

	util.LogJSON(gotStructuralConfigInfo)
}
