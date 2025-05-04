package machinery

import "github.com/RichardKnop/machinery/v2/config"

type MachineryOptions struct {
	Broker       string
	Backend      string
	DefaultQueue string
}

func NewMachineryConfig(opts *MachineryOptions) (*config.Config, error) {
	return &config.Config{
		Broker:        opts.Broker,
		ResultBackend: opts.Backend,
		DefaultQueue:  opts.DefaultQueue,
	}, nil
}
