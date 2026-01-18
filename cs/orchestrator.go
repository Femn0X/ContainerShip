package cs
import (
 "context"
 "fmt"
 "os"
 "strings"
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
 started:=map[string]bool{}
 for len(started)<len(manifest.Services){
  for name,svc:=range manifest.Services{
   if started[name]{
    continue
   }
   ready:=true
   for _,dep:=range svc.DependsOn{
    if !started[dep]{
     ready=false
     continue
    }
   }
   if ready{
    fmt.Println("Started service:",name)
    err:=startContainer(cli,name,svc)
    if err!=nil{
        return err
    }
    started[name]=true
   }
  }
 }
 fmt.Println("All services started!")
 return nil
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
	for name := range manifest.Services {
		fmt.Println("Stopping container:", name)
		if err := cli.ContainerStop(ctx, name, container.StopOptions{}); err != nil {
			// Ignore error if container not found or already stopped
			fmt.Printf("Warning: Could not stop %s: %v\n", name, err)
		}
	}
	fmt.Println("All services stopped!")
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
	for name := range manifest.Services {
		fmt.Printf("Logs for %s:\n", name)
		options := container.LogsOptions{ShowStdout: true, ShowStderr: true, Follow: false, Tail: "50"}
		logs, err := cli.ContainerLogs(ctx, name, options)
		if err != nil {
			fmt.Printf("Error getting logs for %s: %v\n", name, err)
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
	for name := range manifest.Services {
		containers = append(containers, name)
	}
	return nil, containers
}
func startContainer(cli *client.Client,name string,svc Service) error{
    ctx:=context.Background()
    
    // Build Env
    var env []string
    for k, v := range svc.Environment {
        env = append(env, k+"="+v)
    }
    
    // Build PortBindings
    portBindings := nat.PortMap{}
    for _, port := range svc.Ports {
        parts := strings.Split(port, ":")
        if len(parts) == 2 {
            hostPort := parts[0]
            containerPort := parts[1]
            portBindings[nat.Port(containerPort+"/tcp")] = []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: hostPort}}
        }
    }
    
    resp,err:=cli.ContainerCreate(ctx,&container.Config{
        Image:svc.Image,
        Env: env,
    },&container.HostConfig{
        PortBindings: portBindings,
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