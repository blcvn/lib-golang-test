package types

type DockerTask struct {
	Cmd           string   `json: cmd`
	Module        string   `json: module`
	ContainerId   string   `json: containerId`
	ContainerName string   `json: containerName`
	SrcPath       string   `json: srcPath`
	Limit         string   `json: limit`
	DockerImage   string   `json: dockerImage`
	Ports         []string `json: ports`
	Timeout       string   `json: timeout`
	Src           string   `json: src`
	Dst           string   `json: dst`
	VolumeMaps    []string `json: volumeMaps`
}
