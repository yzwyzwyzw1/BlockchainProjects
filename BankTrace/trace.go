package BankTrace

import (
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"

	"strconv"
)


type TraceChainCode struct {

}

func (t *TraceChainCode)Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (t *TraceChainCode)Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fun, args := stub.GetFunctionAndParameters()
	if fun == "loan" {
		return t.Loan(stub,args)
	} else if fun == "repayment" {
		return t.Repayment(stub,args)
	} else if fun == "queryAccountByCardNo" {
		return t.QueryAccountByCardNo(stub,args)
	}
	return shim.Error("指定的操作为非法操作")
}

func saveAccount(stub shim.ChaincodeStubInterface,account Account) bool {

	acc,err := json.Marshal(account)
	if err != nil {
		return  false
	}

	err = stub.PutState(account.CardNo ,acc)
	if err != nil {
		return false
	}

	return true

}


func (t *TraceChainCode)Loan(stub shim.ChaincodeStubInterface,args []string) peer.Response {

	am,err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("填写的贷款金额错误 loan error")
	}
	bank :=Bank {
		BankName:args[1],
		Flag:Bank_Flag_loan,
		Amount:am,
		StartDate:"20100901",
		EndDate:"20101201",
	}

	account := Account{
		CardNo:args[0],
		Aname:"jack",
		Age:29,
		Gender:"男",
		Mobile:"13866667777",
		Bank:bank,
	}

    bl := saveAccount(stub,account)
	if !bl {
		return shim.Error("保存贷款记录失败 loan error")
	}
    return shim.Success([]byte("保存贷款记录成功 loan success"))

}

func (t *TraceChainCode)Repayment(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	am,err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("还款的贷款金额错误")
	}
	bank :=Bank {
		BankName:args[1],
		Flag:Bank_Flag_Repayment,
		Amount:am,
		StartDate:"20100902",
		EndDate:"20101202",
	}

	account := Account{
		CardNo:args[0],
		Aname:"jack",
		Age:29,
		Gender:"男",
		Mobile:"13866667777",
		Bank:bank,
	}
	bl := saveAccount(stub,account)
	if !bl {
		return shim.Error("保存还款记录失败 repayment error")
	}
	return shim.Success([]byte("保存贷还款记录成功 repayment success"))

}

func (t *TraceChainCode)GetAccountByNo(stub shim.ChaincodeStubInterface,cardNo string) (Account,bool,peer.Response ) {
	var account Account
	result,err := stub.GetState(cardNo)
	if err != nil {

		return account,false,shim.Error("Get cardNo error")
	}
	err =json.Unmarshal(result,&account)
	if err != nil {

		return account,false,shim.Error("Get cardNo error")
	}

	return account,true,shim.Success([]byte("Get cardNo success"))
}



// -c '{"Args":["QueryAccountBycardNo","身份证号码"]}'
func (t *TraceChainCode)QueryAccountByCardNo(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	if len(args) !=1 {
		return shim.Error("必须且只能指定要查询的账户信息的身份证号码")
	}

	account,bl,_ := t.GetAccountByNo(stub,args[0])
	if !bl {
		return shim.Error("根据指定的身份证号码查询信息时发生错误 query ID error")
	}

	//查询历史记录信息
	accIterator,err := stub.GetHistoryForKey(account.CardNo)
	if err != nil {
		return shim.Error("查询历史记录信息时发生错误 query history error ")
	}
	defer accIterator.Close()


	//处理查询到的历史信息迭代器对象
	var historys []HistoryItem
	var acc Account
	for accIterator.HasNext() {
		hisData,err := accIterator.Next()
		if err != nil {
			return shim.Error("处理迭代器对象是发生错误 accIterator error")
		}

		var HisItem HistoryItem
		HisItem.TxId = hisData.TxId //获取此条交易数据编号

		err = json.Unmarshal(hisData.Value,&acc)

		if err != nil {
			return shim.Error("反序列化历史状态时发生错误 history error")
		}

		//处理当前记录为nil的情况
		if hisData.Value == nil {
			var empty Account
			HisItem.Account = empty
		}else {
			HisItem.Account = acc
		}
		//将当前处理完毕的历史状态保存到数组中
		historys = append(historys,HisItem)
    }

	account.Historys = historys
	accByte,err := json.Marshal(account)
	if err != nil {
		return shim.Error("将账户信息序列化时发生错误 account Marshal error ")
	}

	return shim.Success(accByte)
}

