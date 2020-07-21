package main

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/itfantasy/gonode/utils/io"
	"github.com/itfantasy/grid/grid-k8s/api"
	"github.com/itfantasy/grid/utils/args"

	v1 "k8s.io/api/core/v1"
)

func run() error {
	parser := args.Parser().
		AddArg("app", "k8cluster", "set the app name")
	app, b := parser.Get("app")
	if !b || app == "" {
		return errors.New("app name (-app) is necessary!")
	}
	return deployEtcdNodes(app)
}

func deployEtcdNodes(appname string) error {
	yamlConf, err := io.LoadFile(io.CurrentDir() + ".etcd.yaml")
	if err != nil {
		return err
	}

	k8sApi := api.NewKubeApi()
	if err := k8sApi.InitKube(); err != nil {
		return err
	}

	var out bytes.Buffer
	var stderr bytes.Buffer
	pods, err := k8sApi.GetPods("", "grid-etcd")
	if err == nil && len(pods) > 0 {
		cmd := exec.Command("kubectl", "delete", "deployment", "grid-etcd")
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("[grid-etcd]:" + fmt.Sprint(err) + ": " + stderr.String())
		}
	}
	makeSurePodDelete(k8sApi, "grid-etcd")

	conf := strings.Replace(yamlConf, "##APPNAME##", appname, -1)
	filePath := io.CurrentDir() + "grid-etcd.yaml"
	if err := io.SaveFile(filePath, conf); err != nil {
		return err
	}
	cmd := exec.Command("kubectl", "apply", "-f", "grid-etcd.yaml")
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		fmt.Println("[grid-etcd]:" + fmt.Sprint(err) + ": " + stderr.String())
		return err
	}
	if err := cmd.Wait(); err != nil {
		fmt.Println("[grid-etcd]:" + fmt.Sprint(err) + ": " + stderr.String())
		return err
	}
	io.DeleteFile(filePath)
	if _, err := makeSurePodRunning(k8sApi, "grid-etcd"); err != nil {
		return err
	}

	cmd = exec.Command("kubectl", "apply", "-f", ".etcd-service.yaml")
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		fmt.Println("[grid-etcd-service]:" + fmt.Sprint(err) + ": " + stderr.String())
		return err
	}
	if err := cmd.Wait(); err != nil {
		fmt.Println("[grid-etcd-service]:" + fmt.Sprint(err) + ": " + stderr.String())
		return err
	}
	return nil
}

func makeSurePodDelete(k8sApi *api.KubeApi, name string) error {
	fmt.Print(name + " deleting")
	counter := 0
	for {
		pods, err := k8sApi.GetPods("", name)
		if err == nil {
			if len(pods) <= 0 {
				fmt.Println("")
				fmt.Println(name + " has been deleted !!")
				return nil
			}
		}
		<-time.After(time.Millisecond * time.Duration(2000))
		fmt.Print(".")
		counter++
		if counter > 150 {
			break
		}
	}
	fmt.Println("")
	return errors.New(name + " delete time out !!")
}

func makeSurePodRunning(k8sApi *api.KubeApi, name string) (*v1.Pod, error) {
	fmt.Print(name + " creating")
	counter := 0
	for {
		pods, err := k8sApi.GetPods("", name)
		if err == nil {
			for _, pod := range pods {
				if pod.Status.PodIP != "" && pod.Status.Phase == v1.PodRunning {
					fmt.Println("")
					fmt.Println(name + " has been created !!")
					return pod, nil
				}
			}
		}
		<-time.After(time.Millisecond * time.Duration(2000))
		fmt.Print(".")
		counter++
		if counter > 150 {
			break
		}
	}
	fmt.Println("")
	return nil, errors.New(name + " create time out !!")
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}
