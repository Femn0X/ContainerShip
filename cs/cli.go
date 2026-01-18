package cs
import(
 "context"
 "fmt"
 "os"
 "github.com/docker/docker/client"
)

func Run() error {
 if len(os.Args)<2{
  printUsage()
 }
 cmd:=os.Args[1]
 switch cmd{
  case "ship":
   err:=Ship()
   if err!=nil{
    return err
   }
  case "stop":
	err:=Stop()
	if err!=nil{
		return err
	}
  case "status":
	err:=getStatus()
	if err!=nil{
		return err
	}
  case "help":
	printUsage()
  case "list":
	err,containers:=listContainers()
	if err!=nil{
		return err
	}
	fmt.Println("Defined containers:")
	for _,c:=range containers{
		fmt.Println(" -",c)
	}
  case "logs":
	err:=PrintLogs()
	if err!=nil{
		return err
	}
   case "exec":
	err:=ExecInContainer(os.Args[2:])
	if err!=nil{
		return err
	}
	case "shell":
	err:=ShellInContainer()
	if err!=nil{
		return err
	}
	case "deploy":
		err:=Deploy()
		if err!=nil{
			return err
		}
	// Add more commands as needed
 default:
  fmt.Println("Command not implemented yet")
 }
 return nil
}
func getStatus() error {
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
		containerJSON, err := cli.ContainerInspect(ctx, name)
		if err != nil {
			fmt.Printf("Container %s: Not found or error: %v\n", name, err)
			continue
		}
		status := containerJSON.State.Status
		fmt.Printf("Container %s: %s\n", name, status)
	}
	return nil
}
func printUsage(){
	fmt.Println(`ContainerShip (cs)

Usage:
	cs ship       Start all services
	cs stop       Stop all services
	cs status     Show service status
	cs list       List defined services
	cs logs       Show logs for all services
	cs exec <svc> <cmd>  Execute command in service
	cs shell <svc>       Start shell in service
	cs deploy     Deploy services (alias for ship)
	cs help       Show this help`)
}