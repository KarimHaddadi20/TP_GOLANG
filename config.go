package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DefaultFile string
	BaseDir     string
	OutDir      string
	DefaultExt  string
	WikiLang    string
	ProcessTopN int
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

	file, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		switch key {
		case "default_file":
			if val != "" {
				cfg.DefaultFile = val
			}
		case "base_dir":
			if val != "" {
				cfg.BaseDir = val
			}
		case "out_dir":
			if val != "" {
				cfg.OutDir = val
			}
		case "default_ext":
			if val != "" {
				cfg.DefaultExt = val
			}
		case "wiki_lang":
			if val != "" {
				cfg.WikiLang = val
			}
		case "process_top_n":
			if n, err := strconv.Atoi(val); err == nil && n > 0 {
				cfg.ProcessTopN = n
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return cfg, err
	}
	return cfg, nil
}
