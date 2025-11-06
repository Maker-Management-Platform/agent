package renderers

import (
	"context"
	"fmt"
	"image/png"
	"io/fs"
	"log"
	"log/slog"

	"github.com/Maker-Management-Platform/fauxgl"
	"github.com/eduardooliveira/stLib/v2/config"
	"github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/library/libfs"
	"github.com/eduardooliveira/stLib/v2/utils"
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

		eye:    fauxgl.V(-3, -3, -0.75),                       // camera position
		center: fauxgl.V(0, -0.07, 0),                         // view center position
		up:     fauxgl.V(0, 0, 1),                             // up vector
		light:  fauxgl.V(-0.75, -5, 0.25).Normalize(),         // light direction
		color:  fauxgl.HexColor(config.Cfg.Render.ModelColor), // object color
	}
}

func (s *stlRenderer) Render(ctx context.Context, asset *entities.Asset) (*entities.Asset, error) {
	genFS, err := libfs.GetLibFS("generated")
	if err != nil {
		return nil, fmt.Errorf("render error getting fs: %w", err)
	}
	imgName := fmt.Sprintf("%s.r.png", asset.ID)

	if _, err := fs.Stat(genFS, imgName); err == nil {
		return entities.NewAsset(genFS.(entities.LibFS), imgName, false, asset), nil
	}

	slog.Info("Rendering", "asset", *asset.Path, "img", imgName, "asset", asset)

	objFs, err := libfs.GetAssetFS(ctx, *asset)
	if err != nil {
		return nil, fmt.Errorf("error getting fs: %w", err)
	}

	f, err := objFs.Open(utils.VoZ(asset.Path))
	if err != nil {
		return nil, err
	}

	resolver, err := fauxgl.LoadSTLF(f)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	f.Close()

	f, err = objFs.Open(utils.VoZ(asset.Path))
	if err != nil {
		return nil, err
	}
	mesh, err := resolver(f)

	if err != nil {
		return nil, err
	}
	// fit mesh in a bi-unit cube centered at the origin
	mesh.BiUnitCube()

	// smooth the normals
	mesh.SmoothNormalsThreshold(fauxgl.Radians(30))

	// create a rendering context
	context := fauxgl.NewContext(s.width*s.scale, s.height*s.scale)
	context.ClearColorBufferWith(fauxgl.HexColor(config.Cfg.Render.BackgroundColor))

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

	writer, err := genFS.Create(imgName)
	if err != nil {
		return nil, err
	}
	defer writer.Close()
	if err := png.Encode(writer, image); err != nil {
		return nil, err
	}

	return entities.NewAsset(genFS.(entities.LibFS), imgName, false, asset), nil
}
