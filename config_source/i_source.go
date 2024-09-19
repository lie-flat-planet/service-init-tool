package config_source

type ISource interface {
	GetFlattenedConfigInfo() (map[string]any, error)
	parse() error
}
