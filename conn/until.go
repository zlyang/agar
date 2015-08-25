package conn

const (
  CDAllClientsType    = "all"    // 所有用户的坐标信息返回，其中第一个是分配给客户端的帐号
  CDUpdateClientsType = "update" // 有更新的用户坐标信息返回
  CDDeleteClientType  = "delete" // 有用户断线
  CDSelfClientType    = "self"   // 连接时返回自己的信息
  CDConnectType       = "connect"

  CUMoveType = "move" // 客户端向服务器发送它的动作
)

// type C2SAction struct {
// 	Type   int
// 	ID     string
// 	Action string
// }

// type S2CSelfInfo struct {
// 	Type         int
// 	CanvasWidth  int
// 	CanvasHeight int
// 	ID           string
// 	Clients      Logic
// }

// type S2CClientInfo struct {
// 	Type    int
// 	Clients []Logic
// }

// type S2CDeleteClient struct {
// 	Type int
// 	Name string
// }

// type Coordinate struct {
// 	X int
// 	Y int
// }

// type Logic struct {
// 	Position Coordinate
// 	Color    string // 显示颜色
// 	Name     string // 名称
// }
