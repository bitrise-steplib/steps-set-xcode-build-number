package step

type Input struct {
	ProjectPath             string `env:"project_path,required"`
	Scheme                  string `env:"scheme,required"`
	Target                  string `env:"target"`
	Configuration           string `env:"configuration"`
	BuildVersion            string `env:"build_version,required"`
	BuildVersionOffset      int64  `env:"build_version_offset"`
	BuildShortVersionString string `env:"build_short_version_string"`
	Verbose                 bool   `env:"verbose,required"`
}

type Config struct {
	ProjectPath             string
	Scheme                  string
	Target                  string
	Configuration           string
	BuildVersion            string
	BuildVersionOffset      int64
	BuildShortVersionString string
}

type Result struct {
	BuildVersion string
}
