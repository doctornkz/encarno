package core

import (
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

// WorkloadLevel arrays may be used to specify workload scenario (warmup, ramp-up, steps, steady)
type WorkloadLevel struct {
	Level  float64
	RampUp time.Duration
	Steady time.Duration
}

type WorkloadMode = string

const (
	WorkloadOpen   WorkloadMode = "open"
	WorkloadClosed WorkloadMode = "closed"
)

type WorkerConf struct {
	Mode             WorkloadMode
	WorkloadSchedule []WorkloadLevel
	StartingWorkers  int
	MaxWorkers       int
}

type BaseWorkload struct {
	Workers      []*Worker
	NibMaker     NibMaker
	StartTime    time.Time
	Output       Output
	Status       Status
	InputPayload func() InputChannel
	Scenario     []WorkloadLevel
	cnt          int
}

func (s *BaseWorkload) SpawnWorker(scheduleChan ScheduleChannel) {
	s.cnt++
	name := "worker#" + strconv.Itoa(s.cnt)
	log.Infof("Spawning worker: %s", name)
	abort := make(chan struct{})
	worker := NewBasicWorker(name, abort, s, scheduleChan)
	s.Workers = append(s.Workers, worker)
	go worker.Run()
}

func (s *BaseWorkload) GenerateSchedule() ScheduleChannel {
	ch := make(ScheduleChannel)
	go func() {
		curLevel := float64(0)
		curOffset := time.Duration(0)
		for _, step := range s.Scenario {
			if step.Level > 0 {
			}

			if step.Steady > 0 {

			}

			ch <- offset
		}
		close(ch)
	}()
	return ch
}

func NewBaseWorkload(maker NibMaker, output Output, inputConfig InputConf, mode WorkloadMode) BaseWorkload {
	var payloadGetter func() InputChannel
	if inputConfig.EnableRegexes || mode == WorkloadClosed {
		inputChannel := NewInput(inputConfig)
		payloadGetter = func() InputChannel {
			return inputChannel
		}
	} else {
		payloadGetter = func() InputChannel {
			inputChannel := NewInput(inputConfig)
			return inputChannel
		}
	}

	return BaseWorkload{
		Workers:      make([]*Worker, 0),
		NibMaker:     maker,
		StartTime:    time.Now(),
		Output:       output,
		Status:       &StatusImpl{},
		InputPayload: payloadGetter,
	}
}
