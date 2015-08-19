package logic

import (
  "math"
  "math/rand"
  "time"

  "github.com/zlyang/agar/server/hub"
)

const (
  CanvasWidth  = 1080 // 画布的宽度
  CanvasHeight = 1920 // 画布的高度
  ObjectWidth  = 9    // 绘制物体的宽度，正方形
)

type Logic struct {
  Position struct { // 当前作标位置
    X int
    Y int
  }
  Color string // 显示颜色
  Name  string // 名称
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
    for c := range hub.H.Users {
      if math.Abs(c.LogicOb.Position.X-x) < ObjectWidth ||
        math.Abs(c.LogicOb.Position.Y-y) < ObjectWidth {
        break
      }
    }

    return &Logic{{x, y}, ""}, nil
  }

  return nil, nil
}
