package mutator

import (
	"math/rand"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

type Caterpillar struct {
	CaterpillarID string `∆.mutator.Moth:"MothID"    ∆.mutator.Butterfly:"ButterflyID" `
	Name          string `∆.mutator.Moth:"MothName"  ∆.mutator.Butterfly:"OriginalCaterpillarName"`
	Ears          int    `∆.mutator.Moth:"EarCount"` //butterflys don't have ears
}

type Butterfly struct {
	ButterflyID             string
	OriginalCaterpillarName string
}

type Moth struct {
	MothID   string
	MothName string
	EarCount int
}

type Tadpole struct {
	Name    string `∆.mutator.Frog:"Name"`
	Species string `∆.mutator.Frog:"Species"`
	Leglets int    `∆.mutator.Frog:"Legs"`
}

type Frog struct {
	Name        string `∆.mutator.Prince:"Name"`
	Species     string `∆.mutator.Prince:"Nationality"`
	Legs        int    `∆.mutator.Prince:"Limbs"`
	PrinceTitle string `∆.mutator.Prince:"Title"`
}

type Prince struct {
	Title       string `∆.mutator.Frog:"PrinceTitle"`
	Name        string `∆.mutator.Frog:"Name"`
	Nationality string `∆.mutator.Frog:"Species"`
	Limbs       int    `∆.mutator.Frog:"Legs"`
}

//Define out mutators types and instances
type Crysalis func(Caterpillar) Butterfly
type Cocoon func(Caterpillar) Moth
type Apopotosis func(Tadpole) Frog
type Kiss func(Frog) Prince
type Magic func(Prince) Frog

var cocoon Cocoon
var crysalis Crysalis
var apopotosis Apopotosis
var kiss Kiss
var magic Magic

func TestMain(m *testing.M) {

	//seed the random number generator
	rand.Seed(time.Now().UTC().UnixNano())

	//create the mutators
	Create(&cocoon)
	Create(&crysalis)
	Create(&apopotosis)
	Create(&kiss)
	Create(&magic)

	//run the tests and exit
	os.Exit(m.Run())
}

// test that we correctly transform using a basic definition
func TestBasicTransform(t *testing.T) {
	d := Tadpole{Name: "Kermit", Leglets: 4, Species: "Muppet"}

	//define the type inline
	var lApopotosis func(Tadpole) Frog
	Create(&lApopotosis)

	f := lApopotosis(d) //Apopotosis is the process that transforms a tadpole to a frog

	dtype := strings.TrimSpace(reflect.TypeOf(d).String())
	ftype := strings.TrimSpace(reflect.TypeOf(f).String())

	if dtype == ftype {
		t.Error("Types are the same", ftype, dtype)
	}
	if dtype != "mutator.Tadpole" {
		t.Error("Wrong type tadpole")
	}
	if ftype != "mutator.Frog" {
		t.Error("Wrong type frog")
	}
	if f.Name != d.Name {
		t.Error("Names do not match", f.Name, d.Name)
	}
	if f.Species != d.Species {
		t.Error("Species do not match", f.Species, d.Species)
	}
	if f.Legs != d.Leglets {
		t.Error("Leg count does not match", f.Legs, d.Leglets)
	}
}

// test that we correctly transform into to two differnt types
func TestDualTransformation(t *testing.T) {
	id := randomString(20)
	c := Caterpillar{CaterpillarID: id, Name: "Hungry Little", Ears: 2}

	b := crysalis(c) //butterflys use a crysalis
	m := cocoon(c)   //moths use a cocoon

	ctype := strings.TrimSpace(reflect.TypeOf(c).String())
	btype := strings.TrimSpace(reflect.TypeOf(b).String())
	mtype := strings.TrimSpace(reflect.TypeOf(m).String())

	if ctype == btype {
		t.Error("Types are the same", btype, ctype)
	}
	if ctype != "mutator.Caterpillar" {
		t.Error("Wrong type caterpillar")
	}

	if btype != "mutator.Butterfly" {
		t.Error("Wrong type butterfly")
	}
	if b.ButterflyID != id {
		t.Error("IDs do not match", b.ButterflyID, c.CaterpillarID)
	}
	if b.ButterflyID != c.CaterpillarID {
		t.Error("IDs do not match", b.ButterflyID, c.CaterpillarID)
	}
	if b.OriginalCaterpillarName != c.Name {
		t.Error("Names do not match", b.OriginalCaterpillarName, c.Name)
	}

	if mtype != "mutator.Moth" {
		t.Error("Wrong type moth")
	}
	if m.MothID != id {
		t.Error("IDs do not match", m.MothID, c.CaterpillarID)
	}
	if m.MothID != c.CaterpillarID {
		t.Error("IDs do not match", m.MothID, c.CaterpillarID)
	}
	if m.MothName != c.Name {
		t.Error("Names do not match", m.MothName, c.Name)
	}
	if m.EarCount != c.Ears {
		t.Error("Ear count does not match", m.EarCount, c.Ears)
	}

}

// test that we correctly transform into and back from a type
func TestCycleTranformation(t *testing.T) {
	p := Prince{Name: "Charming", Title: "Prince of the Blood", Limbs: 4, Nationality: ""}

	f := magic(p) //we use magic to turn a prince into a frog
	q := kiss(f)  //and a kiss to turn a frong back into a prince

	ptype := strings.TrimSpace(reflect.TypeOf(p).String())
	ftype := strings.TrimSpace(reflect.TypeOf(f).String())
	qtype := strings.TrimSpace(reflect.TypeOf(q).String())

	if ptype == ftype {
		t.Error("Types are the same", ptype, ftype)
	}
	if qtype == ftype {
		t.Error("Types are the same", qtype, ftype)
	}

	if ptype != "mutator.Prince" {
		t.Error("Wrong type prince")
	}
	if qtype != "mutator.Prince" {
		t.Error("Wrong type prince")
	}
	if ftype != "mutator.Frog" {
		t.Error("Wrong type frog")
	}

	if p.Name != q.Name {
		t.Error("Names do not match", p.Name, q.Name)
	}
	if p.Title != q.Title {
		t.Error("Titles do not match", p.Title, q.Title)
	}
	if p.Limbs != q.Limbs {
		t.Error("Limbs do not match", p.Limbs, q.Limbs)
	}
	if p.Nationality != q.Nationality {
		t.Error("Nationalities do not match", p.Nationality, q.Nationality)
	}

}

//utility functions

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
