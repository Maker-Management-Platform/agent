package render

import (
	"fmt"
	"log"
	"path"

	"github.com/Maker-Management-Platform/fauxgl"
	"github.com/eduardooliveira/stLib/core/runtime"
	"github.com/eduardooliveira/stLib/core/utils"
	"github.com/nfnt/resize"
)

const (
	scale  = 1    // optional supersampling
	width  = 1920 // output width in pixels
	height = 1080 // output height in pixels
	fovy   = 30   // vertical field of view in degrees
	near   = 1    // near clipping plane
	far    = 10   // far clipping plane
)

var (
	eye    = fauxgl.V(-3, -3, -0.75)                        // camera position
	center = fauxgl.V(0, -0.07, 0)                          // view center position
	up     = fauxgl.V(0, 0, 1)                              // up vector
	light  = fauxgl.V(-0.75, -5, 0.25).Normalize()          // light direction
	color  = fauxgl.HexColor(runtime.Cfg.Render.ModelColor) // object color
)

func renderStl(job RenderJob) (string, error) {
	mesh, err := fauxgl.LoadSTL(utils.ToLibPath(path.Join(job.Project().FullPath(), job.Asset().Name)))
	if err != nil {
		log.Println(err)
		return "", err
	}

	// fit mesh in a bi-unit cube centered at the origin
	mesh.BiUnitCube()

	// smooth the normals
	mesh.SmoothNormalsThreshold(fauxgl.Radians(30))

	// create a rendering context
	context := fauxgl.NewContext(width*scale, height*scale)
	context.ClearColorBufferWith(fauxgl.HexColor(runtime.Cfg.Render.BackgroundColor))

	// create transformation matrix and light direction
	aspect := float64(width) / float64(height)
	matrix := fauxgl.LookAt(eye, center, up).Perspective(fovy, aspect, near, far)

	// use builtin phong shader
	shader := fauxgl.NewPhongShader(matrix, light, eye)
	shader.ObjectColor = color
	context.Shader = shader

	// render
	context.DrawMesh(mesh)

	// downsample image for antialiasing
	image := context.Image()
	image = resize.Resize(width, height, image, resize.Bilinear)

	renderName := fmt.Sprintf("%s.render.png", job.Asset().Name)
	renderSavePath := utils.ToLibPath(path.Join(job.Project().FullPath(), renderName))
	return renderName, fauxgl.SavePNG(renderSavePath, image)
}
