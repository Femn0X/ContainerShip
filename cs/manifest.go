package cs
import(
 "gopkg.in/yaml.v2"
 "os"
)
type HealthConfig struct {
 Test []string `yaml:"test"`
 Interval string `yaml:"interval"`
 Timeout string `yaml:"timeout"`
 Retries int `yaml:"retries"`
 StartPeriod string `yaml:"start_period"`
}
type Service struct{
 Image string `yaml:"image"`
 DependsOn []string `yaml:"depends_on"`
 Ports []string `yaml:"ports,omitempty"`
 Environment map[string]string `yaml:"environment,omitempty"`
 Volumes []string `yaml:"volumes,omitempty"`
 Restart string `yaml:"restart,omitempty"`
 Command []string `yaml:"command,omitempty"`
 EnvFile []string `yaml:"env_file,omitempty"`
 Healthcheck *HealthConfig `yaml:"healthcheck,omitempty"`
 User string `yaml:"user,omitempty"`
 WorkingDir string `yaml:"working_dir,omitempty"`
 Replicas int `yaml:"replicas,omitempty"`
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
