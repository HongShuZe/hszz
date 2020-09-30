```
error: failed to create resmgmt client due to context error: user not found
您需要检查每个组织中的crypto-config根路径（“ client.cryptoconfig.path”）和加密路径（“ organizations.YourOrganizartion.cryptoPath”）。
```

```
创建应用通道失败: create channel failed: create channel failed: SendEnvelope failed: calling orderer 'localhost:7050' failed: Orderer Server Status Code: (400) BAD_REQUEST. Description: error authorizing update: error validating DeltaSet: policy for [Group]  /Channel/Application not satisfied: Failed to reach implicit threshold of 1 sub-policies, required 1 remaining
config.yaml的 organizations:Org1:mspid: **Org1MSP**这个要对应configtx.yaml文件下的mspID
```

```
创建应用通道失败：create channel failed: create channel failed: SendEnvelope failed: calling orderer 'localhost:7050' failed: Orderer Server Status Code: (403) FORBIDDEN. Description: Failed to reach implicit threshold of 1 sub-policies, required 1 remaining: permission denied
原因：docker-compose.yaml的- ORDERER_GENERAL_LOCALMSPID和configtx的mspid没有对应
```


```
链码交互失败：error getting channel response for channel [education]: Discovery status Code: (11) UNKNOWN. Description: error received from Discovery Server: failed constructing descriptor for chaincodes:<name:"educc"
原因：在ccPolicy := cauthdsl.SignedByAnyMember([]string{"Org1MSP"})中的Org1MSP要和configtx.yaml里面的组织MSPID一致
```