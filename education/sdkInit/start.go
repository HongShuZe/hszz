package sdkInit

import (
	"fmt"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
)

const ChaincodeVersion = "1.0"

//初始化fabric SDK
func SetupSDK(ConfigFile string, initialized bool) (*fabsdk.FabricSDK, error) {
	if initialized {
		return nil, fmt.Errorf("fabric SDK 已经被实例化")
	}

	sdk, err := fabsdk.New(config.FromFile(ConfigFile))
	if err != nil {
		return nil, fmt.Errorf("实例化fabric SDK失败：%v", err)
	}

	fmt.Println("fabric SDK初始化成功")
	return sdk, nil
}

// 创建通道并将指定的Peers加入
func CreateChannel(sdk *fabsdk.FabricSDK, info *InitInfo) error {
	clientContext := sdk.Context(fabsdk.WithUser(info.OrgAdmin), fabsdk.WithOrg(info.OrgName))
	if clientContext == nil {
		return fmt.Errorf("根据指定的组织名称与管理员创建资源管理客户端Context失败")
	}

	// 返回资源客户端实例
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		return fmt.Errorf("根据指定的资源管理客户端Context创建通道客户端失败：%v", err)
	}

	mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(info.OrgName))
	if err != nil {
		return fmt.Errorf("根据指定的 OrgName 创建Org MSP 客户端失败：%v", err)
	}

	adminIdentity, err := mspClient.GetSigningIdentity(info.OrgAdmin)
	if err != nil {
		return fmt.Errorf("获取指定id的签名标识失败：%v", err)
	}

	channelReq := resmgmt.SaveChannelRequest{
		ChannelID:         info.ChannelID,
		ChannelConfigPath: info.ChannelConfig,
		SigningIdentities: []msp.SigningIdentity{adminIdentity},
	}
	_, err = resMgmtClient.SaveChannel(channelReq,
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
		resmgmt.WithOrdererEndpoint(info.OrdererOrgName))
	if err != nil {
		return fmt.Errorf("创建应用通道失败：%v", err)
	}
	fmt.Println("通道创建成功")

	info.OrgResMgmt = resMgmtClient

	err = info.OrgResMgmt.JoinChannel(
		info.ChannelID,
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
		resmgmt.WithOrdererEndpoint(info.OrdererOrgName))
	if err != nil {
		return fmt.Errorf("peers 加入应用通道失败：%v", err)
	}

	fmt.Println("peers 已经成功加入通道")
	return nil
}

func InstallAndInstantiateCC(sdk *fabsdk.FabricSDK, info *InitInfo) (*channel.Client, error) {
	fmt.Println("开始安装链码......")
	ccPkg, err := gopackager.NewCCPackage(info.ChaincodePath, info.ChaincodeGoPath)
	if err != nil {
		return nil, fmt.Errorf("创建链码包失败：%v", err)
	}

	installCCReq := resmgmt.InstallCCRequest{
		Name:    info.ChaincodeID,
		Path:    info.ChaincodePath,
		Version: ChaincodeVersion,
		Package: ccPkg,
	}
	_, err = info.OrgResMgmt.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return nil, fmt.Errorf("安装链码失败：%v", err)
	}

	fmt.Println("指定链码安装成功")
	fmt.Println("开始实例化链码")

	ccPolicy := cauthdsl.SignedByAnyMember([]string{"Org1MSP"})
	instantiateCCReq := resmgmt.InstantiateCCRequest{
		Name:    info.ChaincodeID,
		Path:    info.ChaincodePath,
		Version: ChaincodeVersion,
		Args:    [][]byte{[]byte("init")},
		Policy:  ccPolicy,
	}
	_, err = info.OrgResMgmt.InstantiateCC(
		info.ChannelID,
		instantiateCCReq,
		resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return nil, fmt.Errorf("实例化链码失败：%v", err)
	}
	fmt.Println("实例化链码成功")

	clientChannelContext := sdk.ChannelContext(
		info.ChannelID,
		fabsdk.WithUser(info.UserName),
		fabsdk.WithOrg(info.OrgName),
	)
	channelClient, err := channel.New(clientChannelContext)
	if err != nil {
		return nil, fmt.Errorf("创建应用通道客户端失败：%v", err)
	}

	fmt.Println("创建应用通道客户端成功，可以用此客户端调用链码进行查询或执行事务.")
	return channelClient, nil
}
