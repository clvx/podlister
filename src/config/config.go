package config

/*
podlister:
  namespace:
  - foo
    bar
	baz
  labels:
  - foobar
    foobaz
*/

type Config struct {
	Lister `mapstructure:"podlister"`
}
type Lister struct {
	Namespace []string `mapstructure:"namespace"`
	Labels    []string `mapstructure:"labels"`
}
