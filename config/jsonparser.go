package config

import (
	"encoding/json"
	"fmt"
)

func (c *Config) parseJSON(data []byte) error {
	tmp := &tmpConfig{}
	err := json.Unmarshal(data, tmp)
	if err != nil {
		return fmt.Errorf("error parse json config: %w", err)
	}
	fmt.Println("111", tmp)
	c.SetFromTmpConfig(tmp)
	return nil
}
