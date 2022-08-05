package main

import "goreporter/forms"

type Config struct {
	Addr         string         `toml:"addr"`
	Password     string         `toml:"password"`
	Instructions string         `toml:"instructions"`
	Reporter     ReporterConfig `toml:"reporter"`
	Forms        struct {
		Google struct {
			Mapping forms.GoogleFormFieldsMapping
			Params  struct {
				Url      string
				Internal string
			}
		}
		Redmine struct {
			Url string
		}
	}
}

type ReporterConfig struct {
	ProjectTimePrecision int `toml:"project_time_precision"`
	TaskTimePrecision    int `toml:"task_time_precision"`
}
