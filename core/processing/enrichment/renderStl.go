package enrichment

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/Maker-Management-Platform/fauxgl"
	"github.com/eduardooliveira/stLib/core/runtime"
	"github.com/eduardooliveira/stLib/core/utils"
	"github.com/nfnt/resize"
)

type stlRenderer struct {
	scale  int
	width  int
	height int
	fovy   float64
	near   float64
	far    float64

	eye    fauxgl.Vector
	center fauxgl.Vector
	up     fauxgl.Vector
	light  fauxgl.Vector
	color  fauxgl.Color
}

func NewSTLRenderer() *stlRenderer {
	return &stlRenderer{
		scale:  1,    // optional supersampling
		width:  1920, // output width in pixels
		height: 1080, // output height in pixels
		fovy:   30,   // vertical field of view in degrees
		near:   1,    // near clipping plane
		far:    10,   // far clipping plane

		eye:    fauxgl.V(-3, -3, -0.75),                        // camera position
		center: fauxgl.V(0, -0.07, 0),                          // view center position
		up:     fauxgl.V(0, 0, 1),                              // up vector
		light:  fauxgl.V(-0.75, -5, 0.25).Normalize(),          // light direction
		color:  fauxgl.HexColor(runtime.Cfg.Render.ModelColor), // object color
	}
}

func (s *stlRenderer) Render(job Enrichable) (string, error) {
	renderName := fmt.Sprintf("%s.r.png", job.GetAsset().ID)
	renderSavePath := utils.ToAssetsPath(job.GetAsset().ProjectUUID, renderName)

	if _, err := os.Stat(renderSavePath); err == nil {
		return renderName, errors.New("already exists")
	}

	mesh, err := fauxgl.LoadSTL(utils.ToLibPath(path.Join(job.GetProject().FullPath(), job.GetAsset().Name)))
	if err != nil {
		log.Println(err)
		return "", err
	}

	// fit mesh in a bi-unit cube centered at the origin
	mesh.BiUnitCube()

	// smooth the normals
	mesh.SmoothNormalsThreshold(fauxgl.Radians(30))

	// create a rendering context
	context := fauxgl.NewContext(s.width*s.scale, s.height*s.scale)
	context.ClearColorBufferWith(fauxgl.HexColor(runtime.Cfg.Render.BackgroundColor))

	// create transformation matrix and light direction
	aspect := float64(s.width) / float64(s.height)
	matrix := fauxgl.LookAt(s.eye, s.center, s.up).Perspective(s.fovy, aspect, s.near, s.far)

	// use builtin phong shader
	shader := fauxgl.NewPhongShader(matrix, s.light, s.eye)
	shader.ObjectColor = s.color
	context.Shader = shader

	// render
	context.DrawMesh(mesh)

	// downsample image for antialiasing
	image := context.Image()
	image = resize.Resize(uint(s.width), uint(s.height), image, resize.Bilinear)

	return renderName, fauxgl.SavePNG(renderSavePath, image)
}
