package BME280

import (
	"context"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/carlescere/scheduler"
	"github.com/davecheney/i2c"
	"github.com/quinte17/bme280"
)

// log is the default package logger
var log = logger.GetLogger("trigger-tibco-bme280")


// BME280Factory My Trigger factory
type BME280Factory struct{
	metadata *trigger.Metadata
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &BME280Factory{metadata:md}
}

//New Creates a new trigger instance for a given id
func (t *BME280Factory) New(config *trigger.Config) trigger.Trigger {
	return &BME280Trigger{metadata: t.metadata, config:config}
}

// BME280Trigger is a stub for your Trigger implementation
type BME280Trigger struct {
	metadata *trigger.Metadata
	runner   action.Runner
	config   *trigger.Config
	timers map[string]*scheduler.Job
}

// Init implements trigger.Trigger.Init
func (t *BME280Trigger) Init(runner action.Runner) {
	t.runner = runner
	log.Infof("In init, id: '%s', Metadata: '%+v', Config: '%+v'", t.config.Id, t.metadata, t.config)
}

// Metadata implements trigger.Trigger.Metadata
func (t *BME280Trigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Start implements trigger.Trigger.Start
func (t *BME280Trigger) Start() error {
	// start the trigger
	log.Debug("Start Trigger BME280")
	handlers := t.config.Handlers
	t.timers = make(map[string]*scheduler.Job)

	log.Debug("Processing handlers")
	for _, handler := range handlers {
		t.scheduleRepeating(handler)
		log.Debugf("Processing Handler: %s", handler.ActionId)
	}

	return nil
}

// Stop implements trigger.Trigger.Start
func (t *BME280Trigger) Stop() error {
	// stop the trigger
	return nil
}

func (t *BME280Trigger) scheduleRepeating(endpoint *trigger.HandlerConfig){

	interval := 5
	log.Debug("Repeating seconds: ", interval)

	fn2 := func() {
		log.Debug("-- Starting \"Repeating\" (repeat) timer action")

		act := action.Get(endpoint.ActionId)
		log.Debugf("Found action: '%+x'", act)
		log.Debugf("ActionID: '%s'", endpoint.ActionId)
		

		sensorData, err := t.getDataFromSensor(endpoint)
		if err != nil {
			log.Error("Error while reading sensor data. Err: ", err.Error())
		}

		data := make(map[string]interface{})
		data["timestamp"] = t
		data["Temperature"] = sensorData.Temp
		data["Pressure"] = sensorData.Press
		data["Humidity"] = sensorData.Hum

		log.Info("Timestap: ", t, ", Temperature: ", sensorData.Temp, " C, Pressure: ", sensorData.Press, " hPa, Humidity: ", sensorData.Hum, " %")
		startAttrs, err := t.metadata.OutputsToAttrs(data, true)

		if err != nil {
			log.Errorf("After run error' %s'\n", err)
		}

		ctx := trigger.NewContext(context.Background(), startAttrs)
		results, err := t.runner.RunAction(ctx, act, nil)

		if err != nil {
			log.Errorf("An error occured while starting the flow. Err:", err)
		}
		log.Info("Exec: ",results)
	}

	// schedule repeating
	timerJob, err := scheduler.Every(interval).Seconds().Run(fn2)

	if err != nil {
		log.Error("Error scheduleRepeating (first) flow err: ", err.Error())
	}
	if timerJob == nil {
		log.Error("timerJob is nil")
	}

	t.timers["r:"+endpoint.ActionId] = timerJob
}

func (t *BME280Trigger) getDataFromSensor(endpoint *trigger.HandlerConfig) (env bme280.Envdata, err error){
	
	dev, err := i2c.New(0x77, 1)
	if err != nil {
		log.Error(err)
	}
	bme, err := bme280.NewI2CDriver(dev)
	if err != nil {
		log.Error(err)
	}

	return bme.Readenv()
}
