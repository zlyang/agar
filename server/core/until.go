package core

const (
  CDAllClientsType    = 1 // 所有用户的坐标信息返回，其中第一个是分配给客户端的帐号
  CDUpdateClientsType = 2 // 有更新的用户坐标信息返回
  CDDeleteClientType  = 3 // 有用户断线
  CDSelfClientType    = 4 // 连接时返回自己的信息

  CUActionType = 1 // 客户端向服务器发送它的动作
)

type C2SAction struct {
  Type   int
  ID     string
  Action string
}

type S2CSelfInfo struct {
  Type         int
  CanvasWidth  int
  CanvasHeight int
  ID           string
  Clients      Logic
}

type S2CClientInfo struct {
  Type    int
  Clients []Logic
}

type S2CDeleteClient struct {
  Type int
  Name string
}
