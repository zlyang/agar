package core

import (
  "errors"
  "math"
  "math/rand"
  "strconv"
  "time"

  "github.com/busyStone/agar/conn"
)

/*
坐标系
(0,0)-->(CanvasWidth,0)
  |
  |
(0,CanvasHeight)
*/

const (
  ActionLogHandlePeriod = (1000 / 66) * time.Millisecond // 66Hz
  CanvasWidth           = 600                            // 画布的宽度
  CanvasHeight          = 800                            // 画布的高度
  ObjectWidth           = 9                              // 绘制物体的宽度，正方形
  RandString            = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
  RandColorString       = "ABCDEF0123456789"
)

type ActionHandleLog struct {
  ID     string
  Action string
}

var (
  HandleLogicChan chan ActionHandleLog
)

func NewLogicObject() (*conn.Logic, error) {
  r := rand.New(rand.NewSource(time.Now().UnixNano()))

  for i := 0; i < 100; i++ { // 只重试100次，如果没有分配到就返回错误
    x := r.Int31n(CanvasWidth)
    y := r.Int31n(CanvasHeight)

    // 判断有没有碰边
    if x < ObjectWidth/2 || x > (CanvasWidth-ObjectWidth/2) ||
      y < ObjectWidth/2 || y > (CanvasHeight-ObjectWidth/2) {
      break
    }

    // 判断有没有与其它的对象有重叠
    for _, u := range H.Users {
      if math.Abs(float64(u.LogicOb.Position.X-int(x))) < ObjectWidth ||
        math.Abs(float64(u.LogicOb.Position.Y-int(y))) < ObjectWidth {
        break
      }
    }

    return &conn.Logic{
      Position: conn.Coordinate{X: int(x), Y: int(y)},
      Color:    NewColor(), Name: NewName()}, nil
  }

  return nil, errors.New("分配空间失败")
}

func NewName() string { // 不去重，有可能会存在重复的情况
  name := strconv.Itoa(len(H.Users))
  r := rand.New(rand.NewSource(time.Now().UnixNano()))

  for i := 0; i < 5; i++ { // 6个字节长度的名字
    name += string(RandString[r.Int31n(int32(len(RandString)))])
  }

  return name
}

func NewColor() string {
  r := rand.New(rand.NewSource(time.Now().UnixNano()))

  color := "#"
  for i := 0; i < 6; i++ { // 6个字节长度的名字
    color += string(RandColorString[r.Int31n(int32(len(RandColorString)))])
  }

  return color
}

// 周期性处理所有请求，周期小于update clients的周期
// 实时处理，不进行周期。周期处理在这作用不大
func HandleLogicRun() {
  HandleLogicChan = make(chan ActionHandleLog, 1000)

  go func() {
    for {
      select {
      case a := <-HandleLogicChan: // 以进入channel的时间为顺序，不考虑阻塞的情况
        move(a)
      }
    }
  }()
}

func move(a ActionHandleLog) {
  self, ok := H.Users[a.ID]
  if !ok {
    return
  }

  prediction := self.LogicOb.Position
  switch a.Action { // 由于有宽度，且每次都只移动一步，所以不存在有减为负数的情况
  case "R":
    prediction.X += 1
  case "L":
    prediction.X -= 1
  case "U":
    prediction.Y -= 1
  case "D":
    prediction.Y += 1
  case "UL":
    prediction.X -= 1
    prediction.Y -= 1
  case "UR":
    prediction.X += 1
    prediction.Y -= 1
  case "DL":
    prediction.X -= 1
    prediction.Y += 1
  case "DR":
    prediction.X += 1
    prediction.Y += 1
  default:
    return
  }

  // 先查看是否到边沿
  if prediction.X < ObjectWidth/2 || prediction.X > (CanvasWidth-ObjectWidth/2) ||
    prediction.Y < ObjectWidth/2 || prediction.Y > (CanvasHeight-ObjectWidth/2) {
    return
  }

  // 再查看是否与其它玩家有交集
  for n, u := range H.Users {
    if n != a.ID {
      if math.Abs(float64(u.LogicOb.Position.X-prediction.X)) < ObjectWidth ||
        math.Abs(float64(u.LogicOb.Position.Y-prediction.Y)) < ObjectWidth {
        return
      }
    }
  }

  // 有效动作。更新状态
  H.Users[a.ID].LogicOb.Position = prediction
  H.Users[a.ID].Update = true
}
