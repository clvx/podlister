package config

type Config struct {
	Bucket struct {
		Key       string `env:"BUCKET_KEY" env-required:"Bucket key"`
		Secret    string `env:"BUCKET_SECRET" env-required:"Bucket secret"`
		URL       string `yaml:"url" env:"BUCKET_URL" env-default:"https://nyc3.digitaloceanspaces.com"`
		Name      string `yaml:"name" env:"BUCKET_NAME" env-required:"Bucket name"`
		Region    string `yaml:"region" env:"BUCKET_REGION" env-default:"us-east-1"`
		Privilege string `yaml:"privilege" env:"BUCKET_PRIVILEGE" env-default:"public-read"`
	} `yaml:"bucket"`
	Template struct {
		Name   string `yaml:"name" env:"TEMPLATE_NAME" env-default:"index.template"`
		Output string `yaml:"output" env:"TEMPLATE_OUTPUT" env-default:"index.html"`
	} `yaml:"template"`
	Service struct {
		Name   string `yaml:"name" env:"SERVICE_NAME" env-required:"Service name"`
	} `yaml:"service"`
}
