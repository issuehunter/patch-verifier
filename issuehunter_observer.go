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

func _CreateFilter(ethRpcClient *ethrpc.EthRPC, contractAddress string, filter map[string]interface{}) string {
	filterID, err := ethRpcClient.EthNewFilter(filter)
	if err != nil {
		log.Fatal(err)
	}
	return filterID
}

func CreateFilter(ethRpcClient *ethrpc.EthRPC, contractAddress string) string {
	log.Printf("creating a new event filter for the IssueHunter contract deployed at address %v", contractAddress)
	var filter = map[string]interface{}{"address": contractAddress}
	filterID := _CreateFilter(ethRpcClient, contractAddress, filter)
	log.Printf("new event filter for the IssueHunter contract deployed at address %v: %v", contractAddress, filterID)
	return filterID
}

func GetEvents(ethRpcClient *ethrpc.EthRPC, filterID string) []interface{} {
	log.Printf("getting changes for filter %v", filterID)
	logs, err := ethRpcClient.EthGetFilterChanges(filterID)
	if err != nil {
		log.Fatal(err)
	}
	jsonLogs, _ := json.MarshalIndent(logs, "", "  ")
	log.Printf("got changes for filter %v: %v", filterID, string(jsonLogs))
	return logs
}

func UninstallFilter(ethRpcClient *ethrpc.EthRPC, filterID string) bool {
	log.Printf("uninstalling filter %v", filterID)
	uninstall, err := ethRpcClient.EthUninstallFilter(filterID)
	if err != nil {
		log.Fatal(err)
	}
	if !uninstall {
		log.Fatalf("cannot uninstall filter %v", filterID)
	}

	log.Printf("filter %v uninstalled", filterID)
	return uninstall
}

func CreateResolutionProposedFilter(ethRpcClient *ethrpc.EthRPC, contractAddress string) string {
	log.Printf("creating a new filter on ResolutionProposed event for the IssueHunter contract deployed at address %v", contractAddress)
	// > web3.sha3('ResolutionProposed(bytes32,address,bytes32)')
	// "0x3b2c9742ac31922699b361e5e2c3d23b7762a77db8d2ae4b5d8004c904e9a4b7"
	var filter = map[string]interface{}{
		"address": contractAddress,
		"topics":  []interface{}{"0x3b2c9742ac31922699b361e5e2c3d23b7762a77db8d2ae4b5d8004c904e9a4b7"},
	}
	filterID := _CreateFilter(ethRpcClient, contractAddress, filter)
	log.Printf("new event filter for the IssueHunter contract deployed at address %v: %v", contractAddress, filterID)
	return filterID
}

func InteractiveLogObserver(url string, contractAddress string) {
	rpcClient := ethrpc.NewEthRPC(url)
	filterID := CreateFilter(rpcClient, contractAddress)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("quitting...")
		UninstallFilter(rpcClient, filterID)
		os.Exit(1)
	}()
	for {
		GetEvents(rpcClient, filterID)
		time.Sleep(2 * time.Second)
	}
}

func InteractiveResolutionProposedLogObserver(url string, contractAddress string) {
	rpcClient := ethrpc.NewEthRPC(url)
	filterID := CreateResolutionProposedFilter(rpcClient, contractAddress)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("quitting...")
		UninstallFilter(rpcClient, filterID)
		os.Exit(1)
	}()
	for {
		GetEvents(rpcClient, filterID)
		time.Sleep(2 * time.Second)
	}
}
