package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

const DOC_TYPE = "eduObj"

// 保存edu
// args： education
func PutEdu(stub shim.ChaincodeStubInterface, edu Education) ([]byte, bool) {
	edu.ObjectType = DOC_TYPE

	//序列化edu
	eduBytes, err := json.Marshal(edu)
	if err != nil {
		return nil, false
	}

	//保存edu
	err = stub.PutState(edu.EntityID, eduBytes)
	if err != nil {
		return nil, false
	}

	return eduBytes, true
}

// 根据身份证查询信息状态
// args： entityID
func GetEduInfo(stub shim.ChaincodeStubInterface, entityID string) (Education, bool) {

	var edu Education
	// 根据身份证号码查询信息
	eduBytes, err := stub.GetState(entityID)
	if err != nil {
		return edu, false
	}
	if eduBytes == nil {
		return edu, false
	}

	//反序列化查询内容
	err = json.Unmarshal(eduBytes, &edu)
	if err != nil {
		return edu, false
	}

	return edu, true
}

// 根据指定的查询字符串实现富查询
func getEduByQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {
	//富查询
	result, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	// 遍历查询到的内容
	var buffer bytes.Buffer //声明字节缓冲区
	comma := false          //buffer是否要加逗号的判断条件
	for result.HasNext() {
		eduBytes, err := result.Next()
		if err != nil {
			return nil, err
		}

		if comma {
			buffer.WriteString(",")
		}

		buffer.WriteString(string(eduBytes.Value))
		comma = true
	}

	return buffer.Bytes(), nil
}

// 添加信息
// args: educationObject， eventID
// 身份证号为 key, Education 为 value
func (t *EducationChaincode) addEdu(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	//判断参数正确性
	if len(args) != 2 {
		return shim.Error("参数个数有误")
	}
	if args[0] == "" || args[1] == "" {
		return shim.Error("参数值不能为空")
	}

	// 反序列化要添加的edu信息
	var edu Education
	err := json.Unmarshal([]byte(args[0]), &edu)
	if err != nil {
		return shim.Error("反序列化edu数据失败")
	}

	// 查询该身份证是否被使用
	_, bl := GetEduInfo(stub, edu.EntityID)
	if bl {
		return shim.Error("该身份证已经被使用了")
	}

	// 提交信息
	_, bl = PutEdu(stub, edu)
	if !bl {
		return shim.Error("提交信息失败")
	}

	// 添加新增事件
	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("添加信息成功"))
}

// 根据证书编号及姓名查询信息
// args: CertNo, name
func (t *EducationChaincode) queryEduByCertNoAndName(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	//判断参数正确性
	if len(args) != 2 {
		return shim.Error("参数个数有误")
	}

	CertNo := args[0]
	Name := args[1]
	if CertNo == "" || Name == "" {
		return shim.Error("参数值不能为空")
	}

	//拼装CouchDB需要的查询字符串（JSON串）
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\", \"CertNo\":\"%s\", \"Name\":\"%s\"}}", DOC_TYPE, CertNo, Name)
	// 查询edu信息
	eduBytes, err := getEduByQueryString(stub, queryString)
	if err != nil {
		return shim.Error("根据字符串查询信息失败")
	} else if eduBytes == nil {
		return shim.Error("根据字符串查询信息为空")
	}

	return shim.Success(eduBytes)
}

// 根据身份证号码查询详情（溯源）
// args: entityID
func (t *EducationChaincode) queryEduInfoByEntityID(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	//判断参数正确性
	if len(args) != 1 {
		return shim.Error("参数个数有误")
	}

	if args[0] == "" {
		return shim.Error("参数值不能为空")
	}

	// 根据身份证号码查询edu状态
	eduBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("根据身份证查询信息失败")
	} else if eduBytes == nil {
		return shim.Error("根据身份证查询信息为空")
	}

	// 反序列化edu信息
	var edu Education
	err = json.Unmarshal(eduBytes, &edu)
	if err != nil {
		return shim.Error("反序列化edu失败")
	}

	// 根据身份证查询对应历史记录
	result, err := stub.GetHistoryForKey(edu.EntityID)
	if err != nil {
		return shim.Error("根据身份证查询对应历史记录失败")
	}
	defer result.Close()

	//迭代处理历史信息
	var historys []HistoryItem
	var hisEdu Education
	for result.HasNext() {
		byte, err := result.Next()
		if err != nil {
			return shim.Error("获取edu历史数据失败")
		}

		var history HistoryItem
		history.TxId = byte.TxId

		if byte.Value == nil {
			var emtry Education
			history.Education = emtry
		} else {
			json.Unmarshal(byte.Value, &hisEdu)
			history.Education = hisEdu
		}

		historys = append(historys, history)
	}

	edu.Historys = historys

	// 序列化edu信息
	eduBytes, err = json.Marshal(edu)
	if err != nil {
		return shim.Error("序列化edu失败")
	}

	return shim.Success(eduBytes)
}

// 根据身份证号更新信息
// args: educationObject， eventID
func (t *EducationChaincode) updateEdu(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	//判断参数正确性
	if len(args) != 2 {
		return shim.Error("参数个数有误")
	}

	if args[0] == "" || args[1] == "" {
		return shim.Error("参数值不能为空")
	}

	// 反序列化信息
	var info Education
	err := json.Unmarshal([]byte(args[0]), &info)
	if err != nil {
		return shim.Error("反序列化失败")
	}

	// 根据身份证号码查询信息
	edu, bl := GetEduInfo(stub, info.EntityID)
	if !bl {
		return shim.Error("该身份证对应的信息不存在")
	}

	// 修改信息
	edu.Name = info.Name
	edu.Gender = info.Gender
	edu.Nation = info.Nation
	edu.EntityID = info.EntityID
	edu.Place = info.Place
	edu.BirthDay = info.BirthDay

	edu.EnrollDate = info.EnrollDate
	edu.GraduationDate = info.GraduationDate
	edu.SchoolName = info.SchoolName
	edu.Major = info.Major
	edu.QuaType = info.QuaType
	edu.Length = info.Length
	edu.Mode = info.Mode
	edu.Level = info.Level
	edu.Graduation = info.Graduation
	edu.CertNo = info.CertNo
	edu.Photo = info.Photo

	// 提交edu信息
	_, bl = PutEdu(stub, edu)
	if !bl {
		return shim.Error("提交edu数据失败")
	}

	// 添加修改事件
	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("更改信息成功"))
}

// 根据身份证号删除信息
// args: entityID
func (t *EducationChaincode) delEdu(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	//判断参数正确性
	if len(args) != 2 {
		return shim.Error("参数个数有误")
	}

	if args[0] == "" || args[1] == "" {
		return shim.Error("参数值不能为空")
	}

	// 查询该身份证对应的信息是否存在
	_, bl := GetEduInfo(stub, args[0])
	if !bl {
		return shim.Error("该身份证对应的信息不存在")
	}

	// 删除身份证对应的信息
	err := stub.DelState(args[0])
	if err != nil {
		shim.Error("删除信息失败")
	}

	// 添加删除事件
	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		shim.Error(err.Error())
	}

	return shim.Success([]byte("删除信息成功"))
}
