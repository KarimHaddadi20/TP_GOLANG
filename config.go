package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	DefaultFile string `json:"default_file"`
	BaseDir     string `json:"base_dir"`
	OutDir      string `json:"out_dir"`
	DefaultExt  string `json:"default_ext"`
	WikiLang    string `json:"wiki_lang"`
	ProcessTopN int    `json:"process_top_n"`
}

func defaultConfig() Config {
	return Config{
		DefaultFile: "data/input.txt",
		BaseDir:     "data",
		OutDir:      "out",
		DefaultExt:  ".txt",
		WikiLang:    "fr",
		ProcessTopN: 10,
	}
}

func loadConfig(path string) (Config, error) {
	cfg := defaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	var raw Config
	if err := json.Unmarshal(data, &raw); err != nil {
		return cfg, err
	}

	if raw.DefaultFile != "" {
		cfg.DefaultFile = raw.DefaultFile
	}
	if raw.BaseDir != "" {
		cfg.BaseDir = raw.BaseDir
	}
	if raw.OutDir != "" {
		cfg.OutDir = raw.OutDir
	}
	if raw.DefaultExt != "" {
		cfg.DefaultExt = raw.DefaultExt
	}
	if raw.WikiLang != "" {
		cfg.WikiLang = raw.WikiLang
	}
	if raw.ProcessTopN > 0 {
		cfg.ProcessTopN = raw.ProcessTopN
	}
	return cfg, nil
}
