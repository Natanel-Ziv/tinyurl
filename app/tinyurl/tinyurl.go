package tinyurl

type TinyURL interface {
	Ping() bool
}

type Broker struct {
	cfg Config
}

func New(cfg Config) (*Broker, error) {
	b := &Broker{}

	if err := validateConfigs(cfg); err != nil {
		return nil, err
	}
	b.cfg = cfg

	return b, nil
}

func (broker *Broker) Ping() bool {
	return true
}