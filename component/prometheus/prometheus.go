package prometheus

import (
	"github.com/prometheus/client_golang/api"
)

type Prom struct {
	Addr string `env:""`

	client api.Client `skip:""`
}

func (p *Prom) Init() error {
	client, err := api.NewClient(api.Config{Address: p.Addr})
	if err != nil {
		return err
	}

	p.client = client
	return nil
}

func (p *Prom) GetClient() api.Client {
	return p.client
}
