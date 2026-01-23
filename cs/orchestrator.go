package cs
import (
 "context"
 "fmt"
 "os"
 "strings"
 "time"
 "path/filepath"
 "github.com/docker/docker/api/types"
 "github.com/docker/docker/api/types/container"
 "github.com/docker/docker/client"
 "github.com/docker/go-connections/nat"
)
func Ship() error{
 manifest,err:=LoadManifest("containership.yaml")
 if err!=nil{
  return err
 }
 cli,err:=client.NewClientWithOpts(client.FromEnv,client.WithAPIVersionNegotiation())
 if err!=nil{
  return err
 }
 startedServices:=map[string]bool{}
 for len(startedServices)<len(manifest.Services){
  for name,svc:=range manifest.Services{
   if startedServices[name]{
    continue
   }
   ready:=true
   for _,dep:=range svc.DependsOn{
    if !startedServices[dep]{
     ready=false
     break
    }
   }
   if ready{
    fmt.Printf("Starting service: %s\n", name)
    containerNames := getContainerNames(name, svc)
    for i, containerName := range containerNames {
     err:=startContainer(cli,containerName,svc, i)
     if err!=nil{
      return err
     }
    }
    startedServices[name]=true
   }
  }
 }
 fmt.Println("All services started!")
 return nil
}
func getContainerNames(serviceName string, svc Service) []string {
 replicas := svc.Replicas
 if replicas <= 0 {
  replicas = 1
 }
 names := make([]string, replicas)
 for i := 0; i < replicas; i++ {
  if replicas == 1 {
   names[i] = serviceName
  } else {
   names[i] = fmt.Sprintf("%s_%d", serviceName, i+1)
  }
 }
 return names
}
func Stop() error{
	// Implementation for stopping all services
	manifest, err := LoadManifest("containership.yaml")
	if err != nil {
		return err
	}
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	ctx := context.Background()
	for name, svc := range manifest.Services {
		containerNames := getContainerNames(name, svc)
		for _, containerName := range containerNames {
			fmt.Printf("Stopping container: %s\n", containerName)
			if err := cli.ContainerStop(ctx, containerName, container.StopOptions{}); err != nil {
				// Ignore error if container not found or already stopped
				fmt.Printf("Warning: Could not stop %s: %v\n", containerName, err)
			}
		}
	}
	fmt.Println("All services stopped!")
	return nil
}
func Down() error {
	fmt.Println("Stopping and removing services...")
	err := Stop()
	if err != nil {
		return err
	}
	return removeContainers()
}
func Restart() error {
	fmt.Println("Restarting services...")
	err := Stop()
	if err != nil {
		return err
	}
	err = removeContainers()
	if err != nil {
		return err
	}
	return Ship()
}
func removeContainers() error {
	manifest, err := LoadManifest("containership.yaml")
	if err != nil {
		return err
	}
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	ctx := context.Background()
	for name, svc := range manifest.Services {
		containerNames := getContainerNames(name, svc)
		for _, containerName := range containerNames {
			fmt.Printf("Removing container: %s\n", containerName)
			if err := cli.ContainerRemove(ctx, containerName, container.RemoveOptions{Force: true}); err != nil {
				// Ignore error if container not found
				fmt.Printf("Warning: Could not remove %s: %v\n", containerName, err)
			}
		}
	}
	fmt.Println("All containers removed!")
	return nil
}
func PrintLogs() error {
	manifest, err := LoadManifest("containership.yaml")
	if err != nil {
		return err
	}
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	ctx := context.Background()
	for name, svc := range manifest.Services {
		containerNames := getContainerNames(name, svc)
		for _, containerName := range containerNames {
			fmt.Printf("Logs for %s:\n", containerName)
			options := container.LogsOptions{ShowStdout: true, ShowStderr: true, Follow: false, Tail: "50"}
			logs, err := cli.ContainerLogs(ctx, containerName, options)
			if err != nil {
				fmt.Printf("Error getting logs for %s: %v\n", containerName, err)
				continue
			}
			// Read and print logs
			buf := make([]byte, 1024)
			for {
				n, err := logs.Read(buf)
				if n > 0 {
					fmt.Print(string(buf[:n]))
				}
				if err != nil {
					break
				}
			}
			logs.Close()
			fmt.Println()
		}
	}
	return nil
}
func ExecInContainer(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: cs exec <container> <command...>")
	}
	containerName := args[0]
	cmd := args[1:]
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	ctx := context.Background()
	execConfig := types.ExecConfig{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
	}
	execID, err := cli.ContainerExecCreate(ctx, containerName, execConfig)
	if err != nil {
		return err
	}
	resp, err := cli.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{})
	if err != nil {
		return err
	}
	defer resp.Close()
	// For simplicity, just start and wait; full interactive exec would need more work
	return cli.ContainerExecStart(ctx, execID.ID, types.ExecStartCheck{})
}
func ShellInContainer() error {
	args := os.Args[2:] // Assuming called as cs shell <container>
	if len(args) < 1 {
		return fmt.Errorf("usage: cs shell <container>")
	}
	containerName := args[0]
	return ExecInContainer([]string{containerName, "/bin/sh"})
}
func Deploy() error {
	fmt.Println("Deploying services...")
	// For now, just ship
	return Ship()
}
func listContainers() (error, []string) {
	manifest, err := LoadManifest("containership.yaml")
	if err != nil {
		return err, nil
	}
	containers := []string{}
	for name, svc := range manifest.Services {
		containerNames := getContainerNames(name, svc)
		containers = append(containers, containerNames...)
	}
	return nil, containers
}
func loadEnvFromFiles(files []string) (map[string]string, error) {
    env := make(map[string]string)
    for _, file := range files {
        data, err := os.ReadFile(file)
        if err != nil {
            fmt.Printf("Warning: Could not read env file %s: %v\n", file, err)
            continue
        }
        lines := strings.Split(string(data), "\n")
        for _, line := range lines {
            line = strings.TrimSpace(line)
            if line == "" || strings.HasPrefix(line, "#") {
                continue
            }
            parts := strings.SplitN(line, "=", 2)
            if len(parts) == 2 {
                env[parts[0]] = parts[1]
            }
        }
    }
    return env, nil
}
func startContainer(cli *client.Client,name string,svc Service, replicaIndex int) error{
    ctx:=context.Background()
    
    // Load env from files
    fileEnv, err := loadEnvFromFiles(svc.EnvFile)
    if err != nil {
        return err
    }
    
    // Build Env
    envMap := make(map[string]string)
    for k, v := range fileEnv {
        envMap[k] = v
    }
    for k, v := range svc.Environment {
        envMap[k] = v
    }
    var env []string
    for k, v := range envMap {
        env = append(env, k+"="+v)
    }
    
    // Build PortBindings
    portBindings := nat.PortMap{}
    if replicaIndex == 0 {
        for _, port := range svc.Ports {
            parts := strings.Split(port, ":")
            if len(parts) == 2 {
                hostPort := parts[0]
                containerPort := parts[1]
                portBindings[nat.Port(containerPort+"/tcp")] = []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: hostPort}}
            }
        }
    }
    
    // Build Volume mounts
    var binds []string
    for _, vol := range svc.Volumes {
        if strings.Contains(vol, ":") {
            parts := strings.SplitN(vol, ":", 2)
            hostPath := parts[0]
            if !filepath.IsAbs(hostPath) {
                absPath, err := filepath.Abs(hostPath)
                if err != nil {
                    return fmt.Errorf("failed to get absolute path for %s: %v", hostPath, err)
                }
                hostPath = absPath
            }
            binds = append(binds, hostPath+":"+parts[1])
        } else {
            binds = append(binds, vol)
        }
    }
    
    // Build Restart Policy
    restartPolicy := container.RestartPolicy{}
    switch svc.Restart {
    case "always":
        restartPolicy.Name = "always"
    case "unless-stopped":
        restartPolicy.Name = "unless-stopped"
    case "on-failure":
        restartPolicy.Name = "on-failure"
    default:
        restartPolicy.Name = "no"
    }
    
    // Build Healthcheck
    var healthConfig *container.HealthConfig
    if svc.Healthcheck != nil {
        interval, _ := time.ParseDuration(svc.Healthcheck.Interval)
        timeout, _ := time.ParseDuration(svc.Healthcheck.Timeout)
        startPeriod, _ := time.ParseDuration(svc.Healthcheck.StartPeriod)
        healthConfig = &container.HealthConfig{
            Test: svc.Healthcheck.Test,
            Interval: interval,
            Timeout: timeout,
            Retries: svc.Healthcheck.Retries,
            StartPeriod: startPeriod,
        }
    }
    
    resp,err:=cli.ContainerCreate(ctx,&container.Config{
        Image:svc.Image,
        Env: env,
        Cmd: svc.Command,
        Healthcheck: healthConfig,
        User: svc.User,
        WorkingDir: svc.WorkingDir,
    },&container.HostConfig{
        PortBindings: portBindings,
        Binds: binds,
        RestartPolicy: restartPolicy,
    },nil,nil,name)
    if err!=nil{
        return err
    }
    if err:=cli.ContainerStart(ctx,resp.ID,container.StartOptions{}); err!=nil{
        return err
    }
    fmt.Println("Container started:",name)
    return nil
}