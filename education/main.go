package main

import (
	"encoding/json"
	"fmt"
	"github.com/hszz/education/sdkInit"
	"github.com/hszz/education/service"
	"github.com/hszz/education/web"
	"github.com/hszz/education/web/controller"
	"os"
)

const (
	configFile  = "config.yaml"
	initialized = false
	EduCC       = "eduCC"
)

func main() {

	initInfo := &sdkInit.InitInfo{

		ChannelID:     "education",
		ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/hszz/education/fixtures/artifacts/education.tx",

		OrgAdmin:       "Admin",
		OrgName:        "Org1",
		OrdererOrgName: "orderer.hsz.education.com",

		ChaincodeID:     EduCC,
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath:   "github.com/hszz/education/chaincode/",
		UserName:        "User1",
	}

	// 初始化sdk
	sdk, err := sdkInit.SetupSDK(configFile, initialized)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	defer sdk.Close()

	// 创建应用通道
	err = sdkInit.CreateChannel(sdk, initInfo)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// 安装并实例化链码
	channelClient, err := sdkInit.InstallAndInstantiateCC(sdk, initInfo)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(channelClient)

	//+++++++++++++++++++++++

	serviceSetup := service.ServiceSetup{
		ChaincodeID: EduCC,
		Client:      channelClient,
	}

	fmt.Println("hsz")

	edu := service.Education{
		Name:           "张三",
		Gender:         "男",
		Nation:         "汉",
		EntityID:       "101",
		Place:          "潮州",
		BirthDay:       "1999年09月05日",
		EnrollDate:     "2018年九月",
		GraduationDate: "2022年七月",
		SchoolName:     "中国政法大学",
		Major:          "社会学",
		QuaType:        "普通",
		Length:         "四年",
		Mode:           "普通全日制",
		Level:          "本科",
		Graduation:     "未毕业",
		CertNo:         "111",
		Photo:          "/static/photo/11.png",
	}

	edu2 := service.Education{
		Name:           "李四",
		Gender:         "男",
		Nation:         "汉",
		EntityID:       "102",
		Place:          "广州",
		BirthDay:       "1998年09月05日",
		EnrollDate:     "2016年九月",
		GraduationDate: "2020年七月",
		SchoolName:     "中国人民大学",
		Major:          "社会学",
		QuaType:        "普通",
		Length:         "四年",
		Mode:           "普通全日制",
		Level:          "本科",
		Graduation:     "毕业",
		CertNo:         "222",
		Photo:          "/static/photo/22.png",
	}

	msg, err := serviceSetup.SaveEdu(edu)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("信息添加成功， 交易编号为：" + msg)
	}

	msg, err = serviceSetup.SaveEdu(edu2)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("信息添加成功， 交易编号为：" + msg)
	}

	// 根据证书编号和名称查询信息
	eduBytes, err := serviceSetup.FindEduByCertNoAndName("111", "张三")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var edu service.Education
		json.Unmarshal(eduBytes, &edu)
		fmt.Println("根据证书编号和名称查询信息成功")
		fmt.Println(edu)
	}

	// 根据身份号码查询信息
	eduBytes, err = serviceSetup.FindEduInfoByEntityID("102")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var edu service.Education
		json.Unmarshal(eduBytes, &edu)
		fmt.Println("根据身份号码查询信息成功")
		fmt.Println(edu)
	}

	// 修改/添加信息
	info := service.Education{
		Name:           "张三",
		Gender:         "男",
		Nation:         "汉",
		EntityID:       "101",
		Place:          "潮州",
		BirthDay:       "1999年09月05日",
		EnrollDate:     "2016年九月",
		GraduationDate: "2020年七月",
		SchoolName:     "中国政法大学",
		Major:          "社会学",
		QuaType:        "普通",
		Length:         "四年",
		Mode:           "普通全日制",
		Level:          "本科",
		Graduation:     "毕业",
		CertNo:         "333",
		Photo:          "/static/photo/11.png",
	}
	msg, err = serviceSetup.ModifyEdu(info)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("修改信息成功， 交易编号为：" + msg)
	}

	// 根据证书编号和名称查询信息
	eduBytes, err = serviceSetup.FindEduByCertNoAndName("333", "张三")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var edu service.Education
		json.Unmarshal(eduBytes, &edu)
		fmt.Println("根据证书编号和名称查询信息成功")
		fmt.Println(edu)
	}

	// 根据身份证号码查询信息
	eduBytes, err = serviceSetup.FindEduInfoByEntityID("101")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var edu service.Education
		json.Unmarshal(eduBytes, &edu)
		fmt.Println("根据身份号码查询信息成功")
		fmt.Println(edu)
	}

	/*// 删除信息
	msg, err = serviceSetup.DelEdu("101")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("删除信息成功" + msg)
	}

	// 根据身份证号码查询信息
	eduBytes, err = serviceSetup.FindEduInfoByEntityID("101")
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("根据身份号码查询信息失败， 指定的身份证号码信息不存在或被删除")
	} else {
		var edu service.Education
		json.Unmarshal(eduBytes, &edu)
		fmt.Println("根据身份号码查询信息成功")
		fmt.Println(edu)
	}*/

	app := controller.Application{
		Setup: &serviceSetup,
	}
	web.WebStart(app)
}
