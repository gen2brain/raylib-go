package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {

	screenWidth := int32(900)
	screenHeight := int32(500)
	rl.SetConfigFlags(rl.FlagMsaa4xHint) //ENABLE 4X MSAA IF AVAILABLE

	rl.InitWindow(screenWidth, screenHeight, "raylib [shaders] example - basic pbr")

	car := rl.LoadModel("./models/old_car_new.glb")
	plane := rl.LoadModel("./models/plane.glb")

	cam := rl.Camera3D{}
	cam.Fovy = 45
	cam.Position = rl.Vector3{2, 2, 8}
	cam.Projection = rl.CameraPerspective
	cam.Up = rl.Vector3{0, 1, 0}

	shader := rl.LoadShader("./glsl330/pbr.vs", "./glsl330/pbr.fs")

	l := Light{}
	l.SetCombineShader(&shader)
	l.Init(0.0, rl.Vector3{1, 1, 1})
	l1 := l.NewLight(LightTypePoint, rl.Vector3{-1, 1, -2}, rl.Vector3{}, rl.Yellow, 4, &l.Shader)
	l2 := l.NewLight(LightTypePoint, rl.Vector3{2, 1, 1}, rl.Vector3{}, rl.Green, 3.3, &l.Shader)
	l3 := l.NewLight(LightTypePoint, rl.Vector3{-2, 1, 1}, rl.Vector3{}, rl.Red, 8.3, &l.Shader)
	l4 := l.NewLight(LightTypePoint, rl.Vector3{1, 1, -2}, rl.Vector3{}, rl.Blue, 2, &l.Shader)

	p := PhysicRender{}
	p.SetCombineShader(&shader)
	p.Init()

	p.UseTexAlbedo()
	p.UseTexMRA()
	p.UseTexNormal()
	p.UseTexEmissive()

	car.GetMaterials()[0].Shader = shader

	p.AlbedoColorModel(&car, rl.White)
	p.MetallicValue(0.0)
	p.RoughnessValue(0)
	p.EmissivePower(0.0)
	p.AoValue(0.0)
	p.NormalValue(0.2)

	p.EmissiveColor(rl.NewColor(255, 162, 0, 255))
	p.TextureMapAlbedo(&car.GetMaterials()[0], rl.LoadTexture("models/old_car_d.png"))
	p.TextureMapMetalness(&car.GetMaterials()[0], rl.LoadTexture("models/old_car_mra.png"))
	p.TextureMapNormal(&car.GetMaterials()[0], rl.LoadTexture("models/old_car_n.png"))
	rl.SetMaterialTexture(&car.GetMaterials()[0], rl.MapEmission, rl.LoadTexture("models/old_car_e.png"))

	p.SetTiling(rl.NewVector2(1, 1))

	plane.GetMaterials()[0].Shader = shader
	p.TextureMapAlbedo(&plane.GetMaterials()[0], rl.LoadTexture("./models/road_a.png"))
	p.TextureMapNormal(&plane.GetMaterials()[0], rl.LoadTexture("./models/road_n.png"))
	p.TextureMapMetalness(&plane.GetMaterials()[0], rl.LoadTexture("./models/road_mra.png"))

	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {
		rl.UpdateCamera(&cam, rl.CameraOrbital)
		p.UpadteByCamera(cam.Position)
		rl.BeginDrawing()
		rl.DrawFPS(10, 20)
		rl.ClearBackground(rl.Gray)
		rl.BeginMode3D(cam)
		rl.DrawModel(car, rl.Vector3{0, 0.01, 0}, 0.25, rl.RayWhite)
		rl.DrawModel(plane, rl.Vector3{0, 0, 0}, 5, rl.RayWhite)
		l.DrawSpherelight(&l1)
		l.DrawSpherelight(&l2)
		l.DrawSpherelight(&l3)
		l.DrawSpherelight(&l4)
		rl.EndMode3D()
		rl.EndDrawing()
	}
	rl.CloseWindow()
}
