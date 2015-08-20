package ui

import (
	"encoding/binary"
	"log"

	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/gl"
)

var (
	program  gl.Program
	position gl.Attrib
	offset   gl.Uniform
	color    gl.Uniform
	buf      gl.Buffer
)

// Player xxx
type Player struct {
	Pos struct {
		X float32
		Y float32
	}
	Color struct {
		R float32
		G float32
		B float32
	}
}

// OnStart app 获得运行机会是调用
func OnStart() {
	var err error
	program, err = glutil.CreateProgram(vertexShader, fragmentShader)
	if err != nil {
		log.Printf("error creating GL program: %v", err)
		return
	}

	buf = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, buf)
	gl.BufferData(gl.ARRAY_BUFFER, triangleData, gl.STATIC_DRAW)

	position = gl.GetAttribLocation(program, "position")
	color = gl.GetUniformLocation(program, "color")
	offset = gl.GetUniformLocation(program, "offset")

	// TODO(crawshaw): the debug package needs to put GL state init here
	// Can this be an app.RegisterFilter call now??
}

// OnStop app 进入后台时调用
func OnStop() {
	gl.DeleteProgram(program)
	gl.DeleteBuffer(buf)
}

// OnPaint app 帧绘制函数
func OnPaint(players []Player) {
	// gl.ClearColor(1, 0, 0, 1)
	gl.ClearColor(1, 1, 1, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.UseProgram(program)

	for i := 0; i < len(players); i++ {
		gl.Uniform4f(color,
			players[i].Color.R, players[i].Color.G, players[i].Color.B,
			1)

		gl.Uniform2f(offset, players[i].Pos.X, players[i].Pos.Y)

		gl.BindBuffer(gl.ARRAY_BUFFER, buf)
		gl.EnableVertexAttribArray(position)
		gl.VertexAttribPointer(position, coordsPerVertex, gl.FLOAT, false, 0, 0)
		gl.DrawArrays(gl.TRIANGLES, 0, vertexCount)
		gl.DisableVertexAttribArray(position)
	}

}

var triangleData = f32.Bytes(binary.LittleEndian,
	0.0, 0.4, 0.0, // top left
	0.0, 0.0, 0.0, // bottom left
	0.4, 0.0, 0.0, // bottom right
)

const (
	coordsPerVertex = 3
	vertexCount     = 3
)

const vertexShader = `#version 100
uniform vec2 offset;

attribute vec4 position;
void main() {
	// offset comes in with x/y values between 0 and 1.
	// position bounds are -1 to 1.
	vec4 offset4 = vec4(2.0*offset.x-1.0, 1.0-2.0*offset.y, 0, 0);
	gl_Position = position + offset4;
}`

const fragmentShader = `#version 100
precision mediump float;
uniform vec4 color;
void main() {
	gl_FragColor = color;
}`
