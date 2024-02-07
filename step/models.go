package step

type Input struct {
	ProjectPath             string `env:"project_path,required"`
	Scheme                  string `env:"scheme,required"`
	Target                  string `env:"target"`
	Configuration           string `env:"configuration"`
	BuildVersion            int    `env:"build_version,required"`
	BuildVersionOffset      int    `env:"build_version_offset"`
	BuildShortVersionString string `env:"build_short_version_string"`
}

type Config struct {
	ProjectPath             string
	Scheme                  string
	Target                  string
	Configuration           string
	BuildVersion            int
	BuildVersionOffset      int
	BuildShortVersionString string
}

type Result struct {
	BuildVersion int
}
