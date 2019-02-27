package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/itfantasy/gonode/utils/io"
	"github.com/itfantasy/gonode/utils/yaml"

	"github.com/itfantasy/grid/utils/args"
)

type clusterYaml struct {
	AppName string `yaml:"appName"`
	Cluster struct {
		Deployments []struct {
			Name    string `yaml:"name"`
			Enable  int    `yaml:"enable"`
			Num     int    `yaml:"num"`
			Proto   string `yaml:"proto"`
			Port    string `yaml:"port"`
			Public  bool   `yaml:"public"`
			NodeIP  string `yaml:"nodeip"`
			Runtime string `yaml:"runtime"`
			Command string `yaml:"command"`
		}
		StateDeployments []struct {
			Name           string `yaml:"name"`
			Enable         int    `yaml:"enable"`
			Num            int    `yaml:"num"`
			Proto          string `yaml:"proto"`
			PortBase       int    `yaml:"portBase"`
			PortStep       int    `yaml:"portStep"`
			StateIndexBase int    `yaml:"stateIndexBase"`
			StateIndexStep int    `yaml:"stateIndexStep"`
			Public         bool   `yaml:"public"`
			NodeIP         string `yaml:"nodeip"`
			Runtime        string `yaml:"runtime"`
			Command        string `yaml:"command"`
		}
	}
}

func run() error {
	parser := args.Parser().
		AddArg("f", ".cluster.yaml", "set the cluster config file")

	file, _ := parser.Get("f")

	cluster, err := clusterConfigParser(file)
	if err != nil {
		return err
	}
	return deployGridClusterForK8s(cluster)
}

func clusterConfigParser(file string) (*clusterYaml, error) {
	fileContent, err := io.LoadFile(io.CurDir() + file)
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

	deploymentYamlConfig, err := io.LoadFile(io.CurDir() + ".deployment.yaml")
	if err != nil {
		return err
	}
	for _, deployment := range cluster.Cluster.Deployments {
		if err := deployDeployment(deploymentYamlConfig,
			cluster.AppName,
			deployment.Name,
			deployment.Enable,
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

	stateDeploymentYamlConfig, err := io.LoadFile(io.CurDir() + ".state-deployment.yaml")
	if err != nil {
		return err
	}
	for _, stateDeployment := range cluster.Cluster.StateDeployments {
		if err := deployStateDeployment(stateDeploymentYamlConfig,
			cluster.AppName,
			stateDeployment.Name,
			stateDeployment.Enable,
			stateDeployment.Num,
			stateDeployment.Proto,
			stateDeployment.PortBase,
			stateDeployment.PortStep,
			stateDeployment.StateIndexBase,
			stateDeployment.StateIndexStep,
			stateDeployment.Public,
			stateDeployment.NodeIP,
			stateDeployment.Runtime,
			stateDeployment.Command); err != nil {
			return err
		}
	}
	return nil
}

func deployStateDeployment(yamlConfig_ string, appName string, name string, enable int, num int, proto string, portBase int, portStep int, stateIndexBase int, stateIndexStep int, public bool, nodeip string, runtime string, command string) error {
	if enable < 0 {
		for i := 0; i < num; i++ {
			stateIndex := strconv.Itoa(stateIndexBase + stateIndexStep*i)
			insName := name + "-" + stateIndex
			var out bytes.Buffer
			var stderr bytes.Buffer
			cmd := exec.Command("kubectl", "delete", "deployment", insName)
			cmd.Stdout = &out
			cmd.Stderr = &stderr
			if err := cmd.Run(); err != nil {
				fmt.Println("[" + insName + "]:" + fmt.Sprint(err) + ": " + stderr.String())
				continue
			} else {
				fmt.Println(insName + " has been deleted!")
			}
		}
		return nil
	} else if enable == 0 {
		return nil
	}

	for i := 0; i < num; i++ {
		port := strconv.Itoa(portBase + portStep*i)
		stateIndex := strconv.Itoa(stateIndexBase + stateIndexStep*i)

		yamlConfig := strings.Replace(yamlConfig_, "##APPNAME##", appName, -1)
		yamlConfig = strings.Replace(yamlConfig, "##NAME##", name, -1)
		yamlConfig = strings.Replace(yamlConfig, "##NUM##", strconv.Itoa(num), -1)
		yamlConfig = strings.Replace(yamlConfig, "##PROTO##", proto2String(proto), -1)
		yamlConfig = strings.Replace(yamlConfig, "##GRIDPROTO##", proto, -1)
		yamlConfig = strings.Replace(yamlConfig, "##PORT##", port, -1)
		yamlConfig = strings.Replace(yamlConfig, "##NODEIP##", nodeip, -1)
		yamlConfig = strings.Replace(yamlConfig, "##STATEINDEX##", stateIndex, -1)
		yamlConfig = strings.Replace(yamlConfig, "##RUNTIME##", runtime, -1)
		if command != "" {
			yamlConfig = strings.Replace(yamlConfig, "##COMMAND##", ", "+command, -1)
		} else {
			yamlConfig = strings.Replace(yamlConfig, "##COMMAND##", "", -1)
		}

		insName := name + "-" + stateIndex
		filePath := io.CurDir() + "." + insName + ".yaml"
		if err := io.SaveFile(filePath, yamlConfig); err != nil {
			return err
		}

		var out bytes.Buffer
		var stderr bytes.Buffer

		cmd := exec.Command("kubectl", "apply", "-f", "."+name+"-"+stateIndex+".yaml")
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("[" + name + "]:" + fmt.Sprint(err) + ": " + stderr.String())
			return err
		}
		io.DeleteFile(filePath)
		fmt.Println(insName + " has been deployed!")
	}
	return nil
}

func deployDeployment(yamlConfig string, appName string, name string, enable int, num int, proto string, port string, public bool, nodeip string, runtime string, command string) error {
	if enable < 0 {
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd := exec.Command("kubectl", "delete", "deployment", name)
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("[" + name + "]:" + fmt.Sprint(err) + ": " + stderr.String())
			return nil
		}
		if public {
			return deleteService(name)
		}
		return nil
	} else if enable == 0 {
		return nil
	}

	yamlConfig = strings.Replace(yamlConfig, "##APPNAME##", appName, -1)
	yamlConfig = strings.Replace(yamlConfig, "##NAME##", name, -1)
	yamlConfig = strings.Replace(yamlConfig, "##NUM##", strconv.Itoa(num), -1)
	yamlConfig = strings.Replace(yamlConfig, "##PROTO##", proto2String(proto), -1)
	yamlConfig = strings.Replace(yamlConfig, "##GRIDPROTO##", proto, -1)
	yamlConfig = strings.Replace(yamlConfig, "##PORT##", port, -1)
	yamlConfig = strings.Replace(yamlConfig, "##RUNTIME##", runtime, -1)
	if command != "" {
		yamlConfig = strings.Replace(yamlConfig, "##COMMAND##", ", "+command, -1)
	} else {
		yamlConfig = strings.Replace(yamlConfig, "##COMMAND##", "", -1)
	}

	filePath := io.CurDir() + "." + name + ".yaml"
	if err := io.SaveFile(filePath, yamlConfig); err != nil {
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
	io.DeleteFile(filePath)

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

	filePath := io.CurDir() + "." + name + "-service.yaml"
	if err := io.SaveFile(filePath, yamlServiceConfig); err != nil {
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
	io.DeleteFile(filePath)

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
		return nil
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
