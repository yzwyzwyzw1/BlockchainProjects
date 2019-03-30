package main


import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
    "github.com/hyperledger/fabric/scripts/fabric-samples/chaincode/BankTrace"


)




func main() {
	err := shim.Start(new(BankTrace.TraceChainCode))

	if err != nil {
		fmt.Println("启动链码失败")
	}
}