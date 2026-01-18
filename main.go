package main
import(
 "containership/cs"
 "fmt"
 "os"
)
func main(){
 err := cs.Run()
 if err != nil {
  fmt.Println("Error:", err)
  os.Exit(1)
 }
}