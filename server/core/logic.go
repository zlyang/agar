package core

import (
  "errors"
  "math"
  "math/rand"
  "strconv"
  "time"
)

const (
  CanvasWidth  = 1080 // 画布的宽度
  CanvasHeight = 1920 // 画布的高度
  ObjectWidth  = 9    // 绘制物体的宽度，正方形
  RandString   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type Coordinate struct {
  X int
  Y int
}

type Logic struct {
  Position Coordinate
  Color    string // 显示颜色
  Name     string // 名称
}

func NewLogicObject() (*Logic, error) {
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

    return &Logic{Position: Coordinate{X: int(x), Y: int(y)}, Color: "", Name: NewName()}, nil
  }

  return nil, errors.New("分配空间失败")
}

func NewName() string { // 不去重，有可能会存在重复的情况
  name := strconv.Itoa(len(H.Users))
  r := rand.New(rand.NewSource(time.Now().UnixNano()))

  for i := 0; i < 3; i++ {
    name += string(RandString[r.Int31n(int32(len(RandString)))])
  }

  return name
}

// 周期性处理所有请求，周期小于update clients的周期
func HandleLogicRun() {

}
