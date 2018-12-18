package main

import (
	"fmt"
	"io/ioutil"

	"github.com/itfantasy/gonode/utils/io"
	"github.com/itfantasy/gonode/utils/json"
	"github.com/itfantasy/gonode/utils/yaml"
)

type clusterYaml struct {
	appName string `yaml:appName`
	cluster struct {
		deployments []struct {
			name     string `yaml:name`
			replicas int    `yaml:replicas`
			port     int    `yaml:port`
			runtime  string `yaml:runtime`
			command  string `yaml:command`
		}
	}
}

func run() error {
	dict, err := configParser()
	if err != nil {
		return err
	}
	return deployGridForK8s(dict)
}

func configParser() (*clusterYaml, error) {
	fileContent, err := load(io.CurDir() + "cluster.yaml")
	if err != nil {
		return nil, err
	}
	cluster := new(clusterYaml)
	dict, err := yaml.Decode(fileContent, cluster)
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

func deployGridForK8s(cluster *clusterYaml) error {
	_, err := load(io.CurDir() + ".grid-deployment.yaml")
	if err != nil {
		return err
	}
	json.Println(cluster)
	return nil
}

func load(file string) (string, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func save(file string, content string) error {
	data := []byte(content)
	return ioutil.WriteFile(file, data, 0644)
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}
