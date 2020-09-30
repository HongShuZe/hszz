package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type EducationChaincode struct {
}

func (t *EducationChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (t *EducationChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	funcName, args := stub.GetFunctionAndParameters()
	switch funcName {
	case "addEdu":
		return t.addEdu(stub, args)
	case "queryEduInfoByEntityID":
		return t.queryEduInfoByEntityID(stub, args)
	case "queryEduByCertNoAndName":
		return t.queryEduByCertNoAndName(stub, args)
	case "updateEdu":
		return t.updateEdu(stub, args)
	case "delEdu":
		return t.delEdu(stub, args)
	default:
		return shim.Error("函数名错误")
	}
}

func main() {
	err := shim.Start(new(EducationChaincode))
	if err != nil {
		fmt.Printf("启动链码失败")
	}
}
