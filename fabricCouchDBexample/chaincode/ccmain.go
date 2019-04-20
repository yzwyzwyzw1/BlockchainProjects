package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	//"github.com/ethereum/go-ethereum/swarm/api"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	//	"sync/atomic"
)

type CouchDBChaincode struct {

}

func (t *CouchDBChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (t *CouchDBChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fun,args := stub.GetFunctionAndParameters()
	if fun == "billInit" {
		return billInit(stub,args)
	}else if fun == "queryBills" {
		return queryBills(stub,args)
	}else if fun == "queryWaitBills" {
		return queryWaitBills(stub,args)
	}
	return shim.Error("非法操作，指定的函数名无效")
}

func billInit(stub shim.ChaincodeStubInterface ,args []string) peer.Response {
	bill := BillStruct {
		ObjectType: "BillObj",
		BillInfoID: "POC101",
		BillInfoAmt: "1000",
		BillInfoType: "111",
		BillIsseDate:"20100101",
		BillDueDate: "20100110",

		HoldrAcct: "AAA",
		HoldrCmID:"AAAID",

		WaitEndorseAcct: "",
		WaitEndorseCmID: "",
	}
	billByte,_:=json.Marshal(bill)
	err := stub.PutState(bill.BillInfoID,billByte)
	if err != nil {
		return shim.Error("初始化第一个票据失败")
	}

	bill2 := BillStruct {
		ObjectType: "BillObj",
		BillInfoID: "POC102",
		BillInfoAmt: "2000",
		BillInfoType: "111",
		BillIsseDate:"20100201",
		BillDueDate: "20100210",

		HoldrAcct: "AAA",
		HoldrCmID:"AAAID",

		WaitEndorseAcct: "BBB",
		WaitEndorseCmID: "BBBID",
	}
	billByte2,_:=json.Marshal(bill2)
	err = stub.PutState(bill2.BillInfoID,billByte2)
	if err != nil {
		return shim.Error("初始化第二个票据失败")
	}


	bill3 := BillStruct {
		ObjectType: "BillObj",
		BillInfoID: "POC103",
		BillInfoAmt: "2000",
		BillInfoType: "111",
		BillIsseDate:"20100301",
		BillDueDate: "20100310",

		HoldrAcct: "BBB",
		HoldrCmID:"BBBID",

		WaitEndorseAcct: "CCC",
		WaitEndorseCmID: "CCCID",
	}
	billByte3,_:=json.Marshal(bill3)
	err = stub.PutState(bill3.BillInfoID,billByte3)
	if err != nil {
		return shim.Error("初始化第三个票据失败")
	}
	bill4 := BillStruct {
		ObjectType: "BillObj",
		BillInfoID: "POC104",
		BillInfoAmt: "2000",
		BillInfoType: "111",
		BillIsseDate:"20100401",
		BillDueDate: "20100410",

		HoldrAcct: "CCC",
		HoldrCmID:"CCCID",

		WaitEndorseAcct: "BBB",
		WaitEndorseCmID: "BBBID",
	}
	billByte4,_:=json.Marshal(bill4)
	err = stub.PutState(bill4.BillInfoID,billByte4)
	if err != nil {
		return shim.Error("初始化第四个票据失败")
	}


	return shim.Success([]byte(""))
}

func queryBills(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("必须且只能指定持票人的证件号码")
	}
	holdrCmID := args[0]
	//拼装CouchDB所需要的查询字符串是标准的一个JSON字符串
	//“{\"key\":{\"k\":\"v\",\"k\":\"v\"}}“
	queryString := fmt.Sprintf("{\"selector\":{\"DocType\":\"BillObj\",\"HoldrCmID\":\"%s\"}}",holdrCmID)

	result,err := getBillsByQueryString(stub,queryString)
	if err != nil {
		return shim.Error("根据待背书人的证件号码批量查询持票人的持有票据列表")
	}
	return shim.Success(result)
}

func queryWaitBills(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("必须且只能指定持票人的证件号码")
	}
	waitEndorseCmID := args[0]
	queryString := fmt.Sprintf("{\"selector\":{\"DocType\":\"BillObj\",\"WaitEnorseCmID\":\"%s\"}}",waitEndorseCmID)

	result,err := getBillsByQueryString(stub,queryString)
	if err != nil {
		return shim.Error("根据待背书人的证件号码批量查询待背书人的票据列表时发生错误："+ err.Error())
	}
	return shim.Success(result)

}

func getBillsByQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte,error){
	iterator,err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil,err
	}
	defer iterator.Close()
	var isSplit bool
	isSplit = false

	var buffer bytes.Buffer
	for iterator.HasNext() {
		result,err := iterator.Next()
		if err != nil {
			return nil,err
		}

		if isSplit {
			buffer.WriteString(";")
		}

		buffer.WriteString("key:")
		buffer.WriteString(result.Key)
		buffer.WriteString(",value")
		buffer.WriteString(string(result.Value))

		isSplit = true

	}

	return buffer.Bytes(),nil
}


func main() {

	err := shim.Start(new(CouchDBChaincode))
	if err !=nil {
		fmt.Errorf("启动链码失败：%v",err)
	}

}
