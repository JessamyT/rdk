package movementsensor

import (
	"context"

	"google.golang.org/protobuf/types/known/anypb"

	"go.viam.com/rdk/data"
)

type method int64

const (
	readLocation method = iota
	readAltitude method = iota
	readSpeed    method = iota
)

func (m method) String() string {
	switch m {
	case readLocation:
		return "ReadLocation"
	case readAltitude:
		return "ReadAltitude"
	case readSpeed:
		return "ReadSpeed"
	}
	return "Unknown"
}

// TODO: add tests for this file.

func assertMovementSensor(resource interface{}) (MovementSensor, error) {
	ms, ok := resource.(MovementSensor)
	if !ok {
		return nil, data.InvalidInterfaceErr(SubtypeName)
	}
	return ms, nil
}

type lowLevelCollector func(ctx context.Context, ms MovementSensor) (interface{}, error)

func registerCollector(name string, f lowLevelCollector) {
	data.RegisterCollector(data.MethodMetadata{
		Subtype:    SubtypeName,
		MethodName: name,
	}, func(resource interface{}, params data.CollectorParams) (data.Collector, error) {
		ms, err := assertMovementSensor(resource)
		if err != nil {
			return nil, err
		}

		cFunc := data.CaptureFunc(func(ctx context.Context, _ map[string]*anypb.Any) (interface{}, error) {
			v, err := f(ctx, ms)
			if err != nil {
				return nil, data.FailedToReadErr(params.ComponentName, name, err)
			}
			return v, nil
		})
		return data.NewCollector(cFunc, params)
	},
	)
}