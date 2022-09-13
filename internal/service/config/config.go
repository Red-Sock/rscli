package config

import (
	"bufio"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

// root keys names
const (
	dataSource = "data_sources"
)

const (
	sourceNamePg  = "pg"
	sourceNameRds = "rds"
)

func NewConfig(opts map[string][]string) string {
	cfg, err := buildConfig(opts)
	if err != nil {
		return err.Error()
	}

	pth := "./config/local_config.yaml"

	var f *os.File

	_ = os.Mkdir("./config", 0775)
	st, _ := os.Stat(pth)
	if st != nil {
		println("config at" + pth + " already exists. Want to override? (Y)es/(N)o")

		reader := bufio.NewReader(os.Stdin)
		var resp string
		resp, err = reader.ReadString('\n')
		if err != nil {
			return "error reading user response: " + err.Error()
		}
		resp = strings.ToLower(resp)
		if !strings.HasPrefix(resp, "y") {
			return "config creation is aborted by user"
		}
	}
	if f, err = os.Create(pth); err != nil {
		return "error creating file at " + pth + ": " + err.Error()
	}
	defer func() {
		err = f.Close()
		if err != nil {
			println("INTERNAL SYSTEM error when closing file " + pth + " " + err.Error())
		}
	}()

	if err = yaml.NewEncoder(f).Encode(cfg); err != nil {
		return err.Error()
	}
	return "Successfully created config file"
}

func buildConfig(opts map[string][]string) (map[string]interface{}, error) {
	out := map[string]interface{}{}
	ds, err := buildDataSources(opts)
	if err != nil {
		return nil, err
	}
	out[dataSource] = ds
	return out, nil
}

func buildDataSources(opts map[string][]string) (map[string]interface{}, error) {
	cfg := make([]map[string]interface{}, 0, len(opts))

	for f, args := range opts {
		switch strings.Replace(f, "-", "", -1) {
		case sourceNamePg:
			cfg = append(cfg, DefaultPgPattern(args))
		case sourceNameRds:
			cfg = append(cfg, DefaultRdsPattern(args))
		}
	}
	out := make(map[string]interface{})
	for _, item := range cfg {
		for k, v := range item {
			if _, ok := out[k]; ok {
				return nil, fmt.Errorf("colliding names for data sources: %s", k)
			}
			out[k] = v
		}
	}
	return out, nil
}
