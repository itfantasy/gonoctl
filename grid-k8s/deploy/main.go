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

func run() error {
	parser := args.Parser().
		AddArg("f", ".cluster", "set the cluster config file")
	file, _ := parser.Get("f")
	cluster, err := clusterConfigParser(file)
	if err != nil {
		return err
	}
	return deployGridClusterForK8s(cluster)
}

func clusterConfigParser(file string) (*ClusterInfo, error) {
	fileContent, err := io.LoadFile(io.CurrentDir() + file)
	if err != nil {
		return nil, err
	}
	cluster := new(ClusterInfo)
	if yaml.Unmarshal(fileContent, cluster); err != nil {
		return nil, err
	}
	yaml.Println(cluster)
	return cluster, nil
}

func deployGridClusterForK8s(cluster *ClusterInfo) error {
	yamlConf, err := io.LoadFile(io.CurrentDir() + ".deployment.yaml")
	if err != nil {
		return err
	}
	for _, deployment := range cluster.Cluster.Deployments {
		if err := deployDeployment(yamlConf,
			cluster.AppName,
			cluster.Cluster.NameSpace,
			&deployment); err != nil {
			return err
		}
	}
	yamlConf, err = io.LoadFile(io.CurrentDir() + ".state-deployment.yaml")
	if err != nil {
		return err
	}
	for _, stateDeployment := range cluster.Cluster.StateDeployments {
		if err := deployStateDeployment(yamlConf,
			cluster.AppName,
			cluster.Cluster.NameSpace,
			&stateDeployment); err != nil {
			return err
		}
	}
	return nil
}

func parseStateDeploymentPorts(endpoints []string) string {
	var ret = ""
	var itemFormate string = "        - containerPort: %v\r\n          hostPort: %v\r\n          protocol: %v\r\n"
	for _, item := range endpoints {
		info := parseEndpoint(item)
		ret += fmt.Sprint(itemFormate, info[1], info[1], info[0])
	}
	return ret
}

func parseDeploymentPorts(endpoints []string) string {
	var ret = ""
	var itemFormate string = "        - containerPort: %v\r\n          protocol: %v\r\n"
	for _, item := range endpoints {
		info := parseEndpoint(item)
		ret += fmt.Sprint(itemFormate, info[1], info[0])
	}
	return ret
}

func parseServicePorts(endpoints []string) string {
	var ret = ""
	var itemFormate string = "  - port: %v\r\n    targetPort: %v\r\n     nodePort: %v\r\n     protocol: %v\r\n"
	for _, item := range endpoints {
		info := parseEndpoint(item)
		ret += fmt.Sprint(itemFormate, info[1], info[1], info[1], info[0])
	}
	return ret
}

func deployStateDeployment(yamlConf string, appName string, namespace string, deploy *StateDeployment) error {
	if deploy.Enable < 0 {
		for i := 0; i < deploy.Num; i++ {
			stateIndex := strconv.Itoa(deploy.StartIndex + i)
			insName := deploy.Name + "-" + stateIndex
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
	} else if deploy.Enable == 0 {
		return nil
	}

	for i := 0; i < deploy.Num; i++ {
		stateIndex := strconv.Itoa(deploy.StartIndex + i)
		conf := strings.Replace(yamlConf, "##APPNAME##", appName, -1)
		conf = strings.Replace(conf, "##NAMESPACE##", namespace, -1)
		conf = strings.Replace(conf, "##NAME##", deploy.Name, -1)
		conf = strings.Replace(conf, "##NUM##", strconv.Itoa(deploy.Num), -1)
		conf = strings.Replace(conf, "##PORTS##", parseStateDeploymentPorts(deploy.Endpoints), -1)
		conf = strings.Replace(conf, "##STATEINDEX##", stateIndex, -1)
		conf = strings.Replace(conf, "##RUNTIME##", deploy.Runtime, -1)
		conf = strings.Replace(conf, "##COMMAND##", parseCommand(deploy.Command), -1)

		insName := deploy.Name + "-" + stateIndex
		filePath := io.CurrentDir() + "." + insName + ".yaml"
		if err := io.SaveFile(filePath, conf); err != nil {
			return err
		}
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd := exec.Command("kubectl", "apply", "-f", "."+deploy.Name+"-"+stateIndex+".yaml")
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("[" + insName + "]:" + fmt.Sprint(err) + ": " + stderr.String())
			return err
		}
		io.DeleteFile(filePath)
		fmt.Println(insName + " has been deployed!")
	}
	return nil
}

func deployDeployment(yamlConf string, appName string, namespace string, deploy *Deployment) error {
	if deploy.Enable < 0 {
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd := exec.Command("kubectl", "delete", "deployment", deploy.Name)
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("[" + deploy.Name + "]:" + fmt.Sprint(err) + ": " + stderr.String())
			return nil
		}
		return deleteService(deploy.Name)
	} else if deploy.Enable == 0 {
		return nil
	}

	conf := strings.Replace(yamlConf, "##APPNAME##", appName, -1)
	conf = strings.Replace(conf, "##NAMESPACE##", namespace, -1)
	conf = strings.Replace(conf, "##NAME##", deploy.Name, -1)
	conf = strings.Replace(conf, "##NUM##", strconv.Itoa(deploy.Num), -1)
	conf = strings.Replace(conf, "##PORTS##", parseDeploymentPorts(deploy.Endpoints), -1)
	conf = strings.Replace(conf, "##RUNTIME##", deploy.Command, -1)
	conf = strings.Replace(conf, "##COMMAND##", parseCommand(deploy.Command), -1)

	filePath := io.CurrentDir() + "." + deploy.Name + ".yaml"
	if err := io.SaveFile(filePath, conf); err != nil {
		return err
	}

	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("kubectl", "apply", "-f", "."+deploy.Name+".yaml")
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("[" + deploy.Name + "]:" + fmt.Sprint(err) + ": " + stderr.String())
		return err
	}
	io.DeleteFile(filePath)
	return deployService(appName, deploy.Name, deploy.Endpoints)

}

func deployService(appName string, name string, endPoints []string) error {
	yamlConf, err := io.LoadFile(io.CurrentDir() + ".service.yaml")
	if err != nil {
		return err
	}
	conf := strings.Replace(yamlConf, "##APPNAME##", appName, -1)
	conf = strings.Replace(conf, "##NAME##", name, -1)
	conf = strings.Replace(conf, "##PORTS##", parseServicePorts(endPoints), -1)
	filePath := io.CurrentDir() + "." + name + "-service.yaml"
	if err := io.SaveFile(filePath, conf); err != nil {
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

func parseEndpoint(endpoint string) []string {
	ret := make([]string, 0, 2)
	info := strings.Split(endpoint, "://:")
	if info[0] == "kcp" {
		ret = append(ret, "UDP")
	} else {
		ret = append(ret, "TCP")
	}
	ret = append(ret, info[1])
	return ret
}

func parseCommand(command string) string {
	ret := ""
	if command != "" {
		info := strings.Split(command, ",")
		for _, item := range info {
			item = strings.Trim(item, " ")
			ret += ", \"" + item + "\""
		}
	}
	return ret
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}
