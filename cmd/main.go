package main

import "fmt"

const controllerAgentName = "cyberarkcontroller"

var (
	version              string
	kubeconfig           string
	masterURL            string
	logFormat            string
	watchAllNamespaces   bool
	kubeResyncPeriod     int
	cyberarkResyncPeriod int
)

func main() {
	fmt.Println("Do Something")
}
