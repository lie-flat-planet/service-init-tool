package config_source

type Source interface {
	GetFlattenedConfigInfo() (map[string]any, error)
	parse() error
}
