package config

import "fmt"

func Example() {
	exampleCfg, err := NewConfig()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(exampleCfg.GetDBUse())
	fmt.Println(exampleCfg.GetHTTPAddr())
	fmt.Println(exampleCfg.GetNetAddr())

	// Output:
	// false
	// http://localhost:8080
	// localhost:8080
}
