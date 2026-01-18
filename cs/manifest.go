package cs
import(
 "gopkg.in/yaml.v2"
 "os"
)
type Service struct{
 Image string `yaml:"image"`
 DependsOn []string `yaml:"depends_on"`
 Ports []string `yaml:"ports,omitempty"`
 Environment map[string]string `yaml:"environment,omitempty"`
}
type Manifest struct{
 Version string `yaml:"version"`
 Services map[string]Service `yaml:"services"`
}
func LoadManifest(path string)(*Manifest,error){
 data,err:=os.ReadFile(path)
 if err!=nil{
  return nil,err
 }
 var manifest Manifest
 err=yaml.Unmarshal(data,&manifest)
 if err!=nil{
  return nil,err
 }
 return &manifest,nil
}
