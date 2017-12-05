package VL53L0XStreamRPI

import (
	"context"
	"io/ioutil"
	//"encoding/json"
	//"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	//"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

func getJsonMetadata() string{
	jsonMetadataBytes, err := ioutil.ReadFile("trigger.json")
	if err != nil{
		panic("No Json Metadata found for trigger.json path")
	}
	return string(jsonMetadataBytes)
}

const testConfig string = `{
  "id": "mytrigger",
  "settings": {
    "delay_ms": "500"
  },
  "handlers": [
    {
      "actionId": "test_action",
      "settings": {
        "delay_ms": "500"
      }
    }
  ]
}`

type TestRunner struct {
}

// Run implements action.Runner.Run
func (tr *TestRunner) Run(context context.Context, action action.Action, uri string, options interface{}) (code int, data interface{}, err error) {
	log.Debugf("Ran Action: %v", uri)
	return 0, nil, nil
}

func (tr *TestRunner) RunAction(ctx context.Context, act action.Action, options map[string]interface{}) (results map[string]*data.Attribute, err error) {
	log.Debugf("Ran Action: %v", act.Config().Id)
	return nil, nil
}
