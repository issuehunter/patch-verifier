package patchverifier

import (
    "encoding/json"
    "fmt"
    "log"
	"os"
	"os/signal"
	"syscall"
	"time"

    "github.com/issuehunter/ethrpc"
)

func CreateFilter(ethRpcClient *ethrpc.EthRPC, contractAddress string) string {
	log.Printf("creating a new event filter for the IssueHunter contract deployed at address %v", contractAddress)
	var filter = map[string]interface{}{"address": contractAddress}
	filterId, err := ethRpcClient.EthNewFilter(filter)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("new event filter for the IssueHunter contract deployed at address %v: %v", contractAddress, filterId)
	return filterId
}

func GetEvents(ethRpcClient *ethrpc.EthRPC, filterId string) []interface{} {
	log.Printf("getting changes for filter %v", filterId)
	logs, err := ethRpcClient.EthGetFilterChanges(filterId)
	if err != nil {
		log.Fatal(err)
	}
	jsonLogs, _ := json.MarshalIndent(logs, "", "  ")
	log.Printf("got changes for filter %v: %v", filterId, string(jsonLogs))
	return logs
}


func UninstallFilter(ethRpcClient *ethrpc.EthRPC, filterId string) bool {
	log.Printf("uninstalling filter %v", filterId)
	uninstall, err := ethRpcClient.EthUninstallFilter(filterId)
	if err != nil {
		log.Fatal(err)
	}
	if !uninstall {
		log.Fatalf("cannot uninstall filter %v", filterId)
	}

	log.Printf("filter %v uninstalled", filterId)
	return uninstall
}

func InteractiveLogObserver(url string, contractAddress string) {
	rpcClient := ethrpc.NewEthRPC(url)
	filterId := CreateFilter(rpcClient, contractAddress)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<- c
		fmt.Println("quitting...")
		UninstallFilter(rpcClient, filterId)
		os.Exit(1)
	}()
	for {
		GetEvents(rpcClient, filterId)
		time.Sleep(2 * time.Second)
	}
}
