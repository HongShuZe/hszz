#!/bin/bash


echo "一、环境清理"
mkdir -p artifacts
mkdir -p crypto-config
rm -fr artifacts/*
rm -fr crypto-config/*
echo "清理完毕"


echo "生成证书和起始区块信息"
cryptogen generate --config=./crypto-config.yaml
configtxgen -profile TwoOrgsOrdererGenesis -outputBlock ./artifacts/genesis.block

echo "生成通道tx文件"
configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./artifacts/education.tx -channelID education

echo "生成锚节点更新配置文件"
configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./artifacts/Org1MSPanchors.tx -channelID channel -asOrg Org1MSP


CA1_PRIVATE_KEY=$(cd crypto-config/peerOrganizations/org1.hsz.education.com/ca && ls *_sk)

echo "生成ca配置文件"
cp docker-compose-ca-template.yaml docker-compose-ca.yaml
sed -it "s/CA1_PRIVATE_KEY/${CA1_PRIVATE_KEY}/g" docker-compose-ca.yaml