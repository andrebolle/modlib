package main

import (
	"math"
	"math/rand"

	"github.com/ByteArena/box2d"
)

func setupPhysics() (*box2d.B2World, *box2d.B2Body) {
	// Box2D is tuned for meters, kilograms, and seconds.
	// Define the gravity vector.
	gravity := box2d.MakeB2Vec2(0.0, 0.0)

	// Construct a world object, which will hold and simulate the rigid bodies.
	world := box2d.MakeB2World(gravity)

	// Short names
	//kinetic := box2d.B2BodyType.B2_kinematicBody
	dynamic := box2d.B2BodyType.B2_dynamicBody

	// A place to store bodies by name
	//characters := make(map[string]*box2d.B2Body)

	// 1. Create a bodydef. Initial Position, Type (Dynamic, Static,  Kinematic)
	// 2. Create the body from the def.
	// 3. Create a shape - Polygon, Chain, Circle
	// 4. Create a fixture - glues body and shape. denisty, friction, restitution

	// ----------------- Left
	leftBodyDef := box2d.MakeB2BodyDef()
	leftBodyDef.Position.Set(-20, 0)

	leftBody := world.CreateBody(&leftBodyDef)

	leftBox := box2d.MakeB2PolygonShape()

	leftBox.SetAsBox(19.0, 5000.0)

	leftBody.CreateFixture(&leftBox, 0.0)

	// ----------------- Right
	rightBodyDef := box2d.MakeB2BodyDef()
	rightBodyDef.Position.Set(20, 0)

	rightBody := world.CreateBody(&rightBodyDef)

	rightBox := box2d.MakeB2PolygonShape()

	rightBox.SetAsBox(10.0, 5000.0)

	rightBody.CreateFixture(&rightBox, 0.0)

	// ----------------- Ground
	groundBodyDef := box2d.MakeB2BodyDef()
	groundBodyDef.Position.Set(0.0, -20.0)

	groundBody := world.CreateBody(&groundBodyDef)

	groundBox := box2d.MakeB2PolygonShape()

	groundBox.SetAsBox(5000.0, 10.0)

	groundBody.CreateFixture(&groundBox, 0.0)

	// ----------------- Ceiling
	ceilingBodyDef := box2d.MakeB2BodyDef()
	ceilingBodyDef.Position.Set(0.0, 20.0)

	ceilingBody := world.CreateBody(&ceilingBodyDef)

	ceilingBox := box2d.MakeB2PolygonShape()

	ceilingBox.SetAsBox(5000.0, 10.0)

	ceilingBody.CreateFixture(&ceilingBox, 0.0)

	// ----------------- Box
	boxBodyDef := box2d.MakeB2BodyDef()
	boxBodyDef.Position.Set(0, 10)

	boxBodyDef.Type = dynamic
	boxBodyDef.AllowSleep = false
	boxBodyDef.LinearVelocity.Set(4, 2)
	// Body instance
	boxBody := world.CreateBody(&boxBodyDef)
	boxBody.SetTransform(box2d.B2Vec2{X: rand.Float64(), Y: rand.Float64()}, rand.Float64()*2*math.Pi)

	// Create a box shape
	boxShape := box2d.MakeB2PolygonShape()
	boxShape.SetAsBox(1.0, 1.0)

	// Fixture
	// A fixture binds a shape to a body and adds material properties such as density, friction, and restitution.
	// A fixture puts a shape into the collision system (broad-phase) so that it can collide with other shapes.
	fixtureDef := box2d.MakeB2FixtureDef()
	fixtureDef.Shape = &boxShape
	fixtureDef.Density = 1.0
	fixtureDef.Friction = 0.0
	fixtureDef.Restitution = 1
	boxBody.CreateFixtureFromDef(&fixtureDef)

	// physObjList := world.GetBodyList()
	// for i := 0; physObjList != nil; i++ {
	// 	fmt.Println(i)
	// 	physObjList = physObjList.GetNext()
	// }

	return &world, boxBody

}
