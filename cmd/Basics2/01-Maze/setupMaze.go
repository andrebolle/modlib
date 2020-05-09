package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/ByteArena/box2d"
)

// Short names
var dynamic uint8 = box2d.B2BodyType.B2_dynamicBody
var static uint8 = box2d.B2BodyType.B2_staticBody

func initBox2D() box2d.B2World {
	// Box2D is tuned for meters, kilograms, and seconds.
	// Define the gravity vector.
	gravity := box2d.MakeB2Vec2(0.0, 0.0)

	// Construct a world object, which will hold and simulate the rigid bodies.
	world := box2d.MakeB2World(gravity)
	return world
}

func addFrame(world *box2d.B2World, frameSize float64) {
	// A place to store bodies by name
	//characters := make(map[string]*box2d.B2Body)

	// 1. Create a bodydef. Initial Position, Type (Dynamic, Static,  Kinematic)
	// 2. Create the body from the def.
	// 3. Create a shape - Polygon, Chain, Circle
	// 4. Create a fixture - glues body and shape. denisty, friction, restitution
	// shape - a polygon or circle
	// restitution - how bouncy the fixture is
	// friction - how slippery it is
	// density - how heavy it is in relation to its area
	// Fixtures are used to describe the size, shape, and material properties of an object in the physics scene.

	// ----------------- Left
	leftBodyDef := box2d.MakeB2BodyDef()
	leftBodyDef.Position.Set(-frameSize, 0)

	leftBody := world.CreateBody(&leftBodyDef)

	leftBox := box2d.MakeB2PolygonShape()
	leftBox.SetAsBox(10.0, 100.0)

	leftBody.CreateFixture(&leftBox, 0)

	// ----------------- Right
	rightBodyDef := box2d.MakeB2BodyDef()
	rightBodyDef.Position.Set(frameSize, 0)

	rightBody := world.CreateBody(&rightBodyDef)

	rightBox := box2d.MakeB2PolygonShape()

	rightBox.SetAsBox(10.0, 100.0)

	rightBody.CreateFixture(&rightBox, 00)

	// ----------------- Ground
	groundBodyDef := box2d.MakeB2BodyDef()
	groundBodyDef.Position.Set(0, -frameSize)

	groundBody := world.CreateBody(&groundBodyDef)

	groundBox := box2d.MakeB2PolygonShape()

	groundBox.SetAsBox(100.0, 10.0)

	groundBody.CreateFixture(&groundBox, 0)

	// ----------------- Ceiling
	ceilingBodyDef := box2d.MakeB2BodyDef()
	ceilingBodyDef.Position.Set(0.0, frameSize)

	ceilingBody := world.CreateBody(&ceilingBodyDef)

	ceilingBox := box2d.MakeB2PolygonShape()

	ceilingBox.SetAsBox(100.0, 10.0)

	ceilingBody.CreateFixture(&ceilingBox, 0.0)
}

func addStaticBox(world *box2d.B2World, pos, size box2d.B2Vec2) *box2d.B2Body {

	// BoxDef
	boxBodyDef := box2d.MakeB2BodyDef()
	boxBodyDef.Position.Set(pos.X, pos.Y)
	boxBodyDef.Type = static

	// Create body instance
	boxBody := world.CreateBody(&boxBodyDef)
	boxBody.SetTransform(box2d.B2Vec2{X: pos.X, Y: pos.Y}, 0)
	boxBody.SetUserData("box")

	// Create a box/circle shape
	boxShape := box2d.MakeB2PolygonShape()
	boxShape.SetAsBox(size.X, size.Y)

	return boxBody
}

func addDynamicBox(world *box2d.B2World, pos, vel, size box2d.B2Vec2) *box2d.B2Body {

	// ----------------- Box
	boxBodyDef := box2d.MakeB2BodyDef()
	boxBodyDef.Position.Set(pos.X, pos.Y)

	boxBodyDef.Type = dynamic
	boxBodyDef.AllowSleep = false
	boxBodyDef.LinearVelocity.Set(vel.X, vel.Y)
	// Body instance
	boxBody := world.CreateBody(&boxBodyDef)
	boxBody.SetTransform(box2d.B2Vec2{X: pos.X, Y: pos.Y}, rand.Float64()*2*math.Pi)
	boxBody.SetUserData("box")

	// Create a box/circle shape
	boxShape := box2d.MakeB2PolygonShape()
	boxShape.SetAsBox(size.X, size.Y)

	// Fixture
	// A fixture binds a shape to a body and adds material properties such as density, friction, and restitution.
	// A fixture puts a shape into the collision system (broad-phase) so that it can collide with other shapes.
	fixtureDef := box2d.MakeB2FixtureDef()
	fixtureDef.Shape = &boxShape
	fixtureDef.Density = 1.0
	fixtureDef.Friction = 0.0
	fixtureDef.Restitution = 1
	boxBody.CreateFixtureFromDef(&fixtureDef)

	return boxBody
}

func buildMaze(m *Maze) (*box2d.B2World, int) {

	world := initBox2D()
	// frameSize := float64(48)
	// addFrame(&world, frameSize)
	wallCount := 0
	for i := 0; i < m.height; i++ {
		for j := 0; j < m.width; j++ {
			if m.data[i][j] == wall {
				addStaticBox(&world, box2d.B2Vec2{X: float64(i) * 2.1, Y: float64(j) * 2.1}, box2d.B2Vec2{X: 1, Y: 1})
				wallCount++
			}
		}
	}

	fmt.Println("wallCount Box2D", wallCount)

	return &world, wallCount

}
