# censys
Simple library for censys.io

##Usage
```go
import (
  "censys"
 )

func main() {
	auth := [2]string{"API_ID", "Secret"}
	cs, _ := censys.SearchIPv4(auth, "google.com", 1)
	fmt.Println(cs.Results)
}
```
