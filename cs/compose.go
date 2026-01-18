package cs
import(
	"github.com/compose-spec/compose-go/loader"
	"github.com/compose-spec/compose-go/types"
)
func LoadCompose()(*types.Project,error){
	configDetails:=types.ConfigDetails{
		ConfigFiles: []types.ConfigFile{
			{
				Filename:"docker-compose.yml",
			},
		},
	}
	project,err:=loader.Load(configDetails)
	if err!=nil{
		return nil,err
	}
	return project,nil
}
