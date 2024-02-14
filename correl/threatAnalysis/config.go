package threatAnalysis

import "context"

type IThreatAnalysis interface {
	AnalysisEvent() error
}

type FilterConfig struct {
	ctx context.Context
	msgInput <-chan map[string]interface{}
	msgOutput chan <- map[string]interface{}
}

func InitThreatAnalysisConfig(ctx context.Context,chanInput  <-chan map[string]interface{}, chanOutput  chan <- map[string]interface{}) (IThreatAnalysis, error) {
	return &FilterConfig{
		ctx: ctx,
		msgInput : chanInput,
		msgOutput: chanOutput,
	}, nil
}

