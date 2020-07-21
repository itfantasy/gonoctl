package main

type ClusterInfo struct {
	AppName string
	Cluster struct {
		NameSpace        string
		RegDC            string
		Deployments      []Deployment
		StateDeployments []StateDeployment
	}
}

type Deployment struct {
	Name      string
	Num       int
	Endpoints []string
	Proj      string
	BackEnds  string
	IsPub     int
	Command   string
	Enable    int
}

type StateDeployment struct {
	Name       string
	Num        int
	StartIndex int
	Endpoints  []string
	Proj       string
	BackEnds   string
	IsPub      int
	Command    string
	Enable     int
}
