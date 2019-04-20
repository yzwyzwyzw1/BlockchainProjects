[TOC]

---
# 启动网络

```
make env-up
```
启动之后会看到这些容器（使用 docker ps 命令）：
```
CONTAINER ID        IMAGE                        COMMAND                  CREATED             STATUS                  PORTS                                            NAMES
715d1b4d63f7        hyperledger/fabric-tools     "/bin/bash"              1 second ago        Up Less than a second                                                    cli
ed7b07683358        hyperledger/fabric-peer      "peer node start"        7 seconds ago       Up 1 second             0.0.0.0:8051->8051/tcp, 0.0.0.0:8053->8053/tcp   peer0.org2.example.com
df7ae0eb0e41        hyperledger/fabric-peer      "peer node start"        7 seconds ago       Up 2 seconds            0.0.0.0:8151->8051/tcp, 0.0.0.0:8153->8053/tcp   peer1.org2.example.com
157a87a26c9e        hyperledger/fabric-peer      "peer node start"        7 seconds ago       Up 3 seconds            0.0.0.0:7051->7051/tcp, 0.0.0.0:7053->7053/tcp   peer0.org1.example.com
bb688222e61e        hyperledger/fabric-peer      "peer node start"        7 seconds ago       Up 4 seconds            0.0.0.0:7151->7051/tcp, 0.0.0.0:7153->7053/tcp   peer1.org1.example.com
fe8d0f363cb8        hyperledger/fabric-ca        "sh -c 'fabric-ca-se…"   8 seconds ago       Up 5 seconds            0.0.0.0:7054->7054/tcp                           ca.org1.example.com
b9b5aa3c31d2        hyperledger/fabric-ca        "sh -c 'fabric-ca-se…"   8 seconds ago       Up 6 seconds            7054/tcp, 0.0.0.0:7055->7055/tcp                 ca.org2.example.com
6c7c73597c11        hyperledger/fabric-couchdb   "tini -- /docker-ent…"   10 seconds ago      Up 7 seconds            4369/tcp, 9100/tcp, 0.0.0.0:5984->5984/tcp       couchdb
3c787bd51cb5        hyperledger/fabric-orderer   "orderer"                10 seconds ago      Up 8 seconds            0.0.0.0:7050->7050/tcp                           orderer.example.com
```
如果要关闭网络可以使用命令
```
make env-down
```
如果要清除网络可以使用命令
```
make env-clean
```

# 进入couchDB网页版
摘录：建议CORE_LEDGER_STATE_COUCHDBCONFIG_COUCHDBADDRESS变量在与Peer节点相同的服务器上运行CouchDB，而不是将CouchDB容器端口映射到Docker-Compose中的服务器端口。否则在CoucnDB客户端（Peer 节点中）和服务器之间的连接上提供适当的安全性。其默认值为为
127.0.0.1:5984。当CouchDB被成功启动之后，可以通过IP：PORT/_utils的方式进入CouchDB的网页版入口界面。

登录密码需要查看配置文件docker-compose-bash.yaml

# 部署链码

1.根据持票人ID查询账单
```
docker exec -it cli bash
export CHANNEL_NAME=mychannel
peer channel create -o orderer.example.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/mychannel.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
peer channel join -b mychannel.block
peer chaincode install -n mycc -v 1.0 -p github.com/chaincode
peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n mycc -v 1.0 -c'{"Args":["init"]}' -P "OR ('Org1MSP.peer')"
peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n mycc  -c '{"Args":["billInit"]}'
peer chaincode query -C $CHANNEL_NAME -n mycc -c '{"Args":["queryBills","AAAID"]}'
```
输出结果：
```
key:POC101,value{"BillDueDate":"20100110","BillInfoAmt":"1000","BillInfoID":"POC101","BillInfoType":"111","BillIsseDate":"20100101","DocType":"BillObj","HoldrAcct":"AAA","HoldrCmID":"AAAID","WaitEndorseAcct":"","WaitEnorseCmID":""};
key:POC102,value{"BillDueDate":"20100210","BillInfoAmt":"2000","BillInfoID":"POC102","BillInfoType":"111","BillIsseDate":"20100201","DocType":"BillObj","HoldrAcct":"AAA","HoldrCmID":"AAAID","WaitEndorseAcct":"BBB","WaitEnorseCmID":"BBBID"}
```
其他测试：
```
peer chaincode query -C $CHANNEL_NAME -n mycc -c '{"Args":["queryBills","BBBID"]}'
```
输出结果：
```
key:POC103,value{"BillDueDate":"20100310","BillInfoAmt":"2000","BillInfoID":"POC103","BillInfoType":"111","BillIsseDate":"20100301","DocType":"BillObj","HoldrAcct":"BBB","HoldrCmID":"BBBID","WaitEndorseAcct":"CCC","WaitEnorseCmID":"CCCID"}
```
其他测试：
```
peer chaincode query -C $CHANNEL_NAME -n mycc -c '{"Args":["queryBills","CCCID"]}'
```
输出结果：
```
key:POC104,value{"BillDueDate":"20100410","BillInfoAmt":"2000","BillInfoID":"POC104","BillInfoType":"111","BillIsseDate":"20100401","DocType":"BillObj","HoldrAcct":"CCC","HoldrCmID":"CCCID","WaitEndorseAcct":"BBB","WaitEnorseCmID":"BBBID"}
```

2.查询用户BBB的待背书的账单
```
peer chaincode query -C $CHANNEL_NAME -n mycc -c '{"Args":["queryWaitBills","BBBID"]}'
```
输出结果：
```
key:POC102,value{"BillDueDate":"20100210","BillInfoAmt":"2000","BillInfoID":"POC102","BillInfoType":"111","BillIsseDate":"20100201","DocType":"BillObj","HoldrAcct":"AAA","HoldrCmID":"AAAID","WaitEndorseAcct":"BBB","WaitEnorseCmID":"BBBID"};
key:POC104,value{"BillDueDate":"20100410","BillInfoAmt":"2000","BillInfoID":"POC104","BillInfoType":"111","BillIsseDate":"20100401","DocType":"BillObj","HoldrAcct":"CCC","HoldrCmID":"CCCID","WaitEndorseAcct":"BBB","WaitEnorseCmID":"BBBID"}
```
其他测试：
```
peer chaincode query -C $CHANNEL_NAME -n mycc -c '{"Args":["queryWaitBills","CCCID"]}'
```
输出结果：
```
key:POC103,value{"BillDueDate":"20100310","BillInfoAmt":"2000","BillInfoID":"POC103","BillInfoType":"111","BillIsseDate":"20100301","DocType":"BillObj","HoldrAcct":"BBB","HoldrCmID":"BBBID","WaitEndorseAcct":"CCC","WaitEnorseCmID":"CCCID"}
```
其他测试：
```
peer chaincode query -C $CHANNEL_NAME -n mycc -c '{"Args":["queryWaitBills","AAAID"]}'
```
输出结果：
输出为空，无账单信息！

# CouchDB索引应用

1.手动创建文档数据

在web界面中有一个create document的选项，点击输入以下测试数据即可插入新的数据。

```
{
  "_id": "POC105",
  "BillDueDate": "20100410",
  "BillInfoAmt": "2000",
  "BillInfoID": "POC104",
  "BillInfoType": "111",
  "BillIsseDate": "20100401",
  "DocType": "BillObj",
  "HoldrAcct": "CCC",
  "HoldrCmID": "CCCID",
  "WaitEndorseAcct": "BBB",
  "WaitEnorseCmID": "BBBID",
  "~version": "\u0000CgMBAgA="
}
```

中间拦中有一个”十字“按钮，选择下拉菜单中的”Mango Indexes“进入索引界面

修改
```
{
   "index": {
      "fields": [
         "foo",
         "bar"
      ]
   },
   "name": "foo-bar-json-index",
   "type": "json"
}
```
修改结果如下：

TODO

```
---
# 出错及解决方案：

Found map[string]interface{} value for peer.BCCSP
2019-03-26 12:41:50.074 UTC [viperutil] unmarshalJSON -> DEBU 002 Unmarshal JSON: value cannot be unmarshalled: invalid character 'S' looking for beginning of value
 --->(应该是序列化时出的错)仔细检查queryString := fmt.Sprintf("{\"selector\"{\"DocType\":\"BillObj\",\"HoldrCmID\":\"%s\"}}",holdrCmID)