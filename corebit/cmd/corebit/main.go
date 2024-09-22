package main

import (
	"fmt"
	"github.com/example/corebit/pkg/workflow"
)

func main() {
	fmt.Println("Starting CoreBit System...")
	workflow.RunWorkflow("config.yaml")
}
