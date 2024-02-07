package config

import "fmt"

func Example() {
	cfg, err := NewConfig()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(cfg.GetDBUse())
	fmt.Println(cfg.GetHTTPAddr())
	fmt.Println(cfg.GetNetAddr())

	// Output:
	// false
	// http://localhost:8080
	// localhost:8080
}
