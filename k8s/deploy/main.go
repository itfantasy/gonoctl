package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/itfantasy/gonode/utils/io"
	"github.com/itfantasy/gonode/utils/yaml"
)

type clusterYaml struct {
	AppName string `yaml:"appName"`
	Cluster struct {
		Deployments []struct {
			Name    string `yaml:"name"`
			Num     int    `yaml:"num"`
			Proto   string `yaml:"proto"`
			Port    string `yaml:"port"`
			Public  bool   `yaml:"public"`
			NodeIP  string `yaml:"nodeip"`
			Runtime string `yaml:"runtime"`
			Command string `yaml:"command"`
		}
	}
}

func run() error {
	cluster, err := clusterConfigParser()
	if err != nil {
		return err
	}
	return deployGridClusterForK8s(cluster)
}

func clusterConfigParser() (*clusterYaml, error) {
	fileContent, err := io.LoadFile(io.CurDir() + ".cluster.yaml")
	if err != nil {
		return nil, err
	}
	cluster := new(clusterYaml)
	if yaml.Decode(fileContent, cluster); err != nil {
		return nil, err
	}
	yaml.Println(cluster)
	return cluster, nil
}

func deployGridClusterForK8s(cluster *clusterYaml) error {
	yamlConfig, err := io.LoadFile(io.CurDir() + ".deployment.yaml")
	if err != nil {
		return err
	}
	for _, deployment := range cluster.Cluster.Deployments {
		if err := deployDeployment(yamlConfig,
			cluster.AppName,
			deployment.Name,
			deployment.Num,
			deployment.Proto,
			deployment.Port,
			deployment.Public,
			deployment.NodeIP,
			deployment.Runtime,
			deployment.Command); err != nil {
			return err
		}
	}
	return nil
}

func deployDeployment(yamlConfig string, appName string, name string, num int, proto string, port string, public bool, nodeip string, runtime string, command string) error {
	if num < 0 {
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd := exec.Command("kubectl", "delete", "deployment", name)
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("[" + name + "]:" + fmt.Sprint(err) + ": " + stderr.String())
			return err
		}
		if public {
			return deleteService(name)
		}
		return nil
	}

	yamlConfig = strings.Replace(yamlConfig, "##APPNAME##", appName, -1)
	yamlConfig = strings.Replace(yamlConfig, "##NAME##", name, -1)
	yamlConfig = strings.Replace(yamlConfig, "##NUM##", strconv.Itoa(num), -1)
	yamlConfig = strings.Replace(yamlConfig, "##PROTO##", proto, -1)
	yamlConfig = strings.Replace(yamlConfig, "##PORT##", port, -1)
	yamlConfig = strings.Replace(yamlConfig, "##RUNTIME##", runtime, -1)
	if command != "" {
		yamlConfig = strings.Replace(yamlConfig, "##COMMAND##", ", "+command, -1)
	} else {
		yamlConfig = strings.Replace(yamlConfig, "##COMMAND##", "", -1)
	}
	if err := io.SaveFile(io.CurDir()+"."+name+".yaml", yamlConfig); err != nil {
		return err
	}

	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("kubectl", "apply", "-f", "."+name+".yaml")
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("[" + name + "]:" + fmt.Sprint(err) + ": " + stderr.String())
		return err
	}

	if public {
		return deployService(appName, name, proto, port, nodeip)
	}

	return nil
}

var _yamlServiceConfig string = ""

func deployService(appName string, name string, proto string, port string, nodeip string) error {

	yamlServiceConfig, err := io.LoadFile(io.CurDir() + ".service.yaml")
	if err != nil {
		return err
	}

	yamlServiceConfig = strings.Replace(yamlServiceConfig, "##APPNAME##", appName, -1)
	yamlServiceConfig = strings.Replace(yamlServiceConfig, "##NAME##", name, -1)
	yamlServiceConfig = strings.Replace(yamlServiceConfig, "##PROTO##", proto2String(proto), -1)
	yamlServiceConfig = strings.Replace(yamlServiceConfig, "##PORT##", port, -1)
	yamlServiceConfig = strings.Replace(yamlServiceConfig, "##NODEIP##", nodeip, -1)
	yamlServiceConfig = strings.Replace(yamlServiceConfig, "##NODEPORT##", port, -1)

	if err := io.SaveFile(io.CurDir()+"."+name+"-service.yaml", yamlServiceConfig); err != nil {
		return err
	}

	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("kubectl", "apply", "-f", "."+name+"-service.yaml")
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("[" + name + "]:" + fmt.Sprint(err) + ": " + stderr.String())
		return err
	}

	return nil
}

func deleteService(name string) error {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("kubectl", "delete", "service", name+"-service")
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("[" + name + "]:" + fmt.Sprint(err) + ": " + stderr.String())
		return err
	}
	return nil
}

func proto2String(proto string) string {
	switch proto {
	case "kcp":
		return "UDP"
	case "tcp":
		return "TCP"
	}
	return ""
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}
