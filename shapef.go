package main

import (
	"encoding/json"
	"fmt"
	"github.com/jonas-p/go-shp"
	"io/ioutil"
	"os"
	"strings"
)

// jShape is the type that is stored in JSON files.
type jShape struct {
	Points [][]float64            `json:"points"`
	Data   map[string]interface{} `json:"data"`
}

func (shape jShape) Polygon() *shp.Polygon {
	pointCount := len(shape.Points)
	points := make([]shp.Point, pointCount)
	for i, jPoint := range shape.Points {
		points[i] = shp.Point{jPoint[0], jPoint[1]}
	}
	return &shp.Polygon{
		Box:       shp.BBoxFromPoints(points),
		NumParts:  1,
		NumPoints: int32(pointCount),
		Parts:     []int32{0},
		Points:    points}
}

func check(err error, message string) {
	if err != nil {
		fmt.Println(message)
		panic(err)
	}
}

func readJSON(filePath string) []jShape {
	var shapes []jShape
	bytes, err := ioutil.ReadFile(filePath)
	check(err, "failed to read "+filePath)
	err = json.Unmarshal(bytes, &shapes)
	check(err, "failed to parse"+filePath)
	return shapes
}

type fieldDef struct {
	name    string
	numeric bool
	length  int
}

func (def fieldDef) Field() shp.Field {
	if def.numeric {
		return shp.FloatField(def.name, 32, 4)
	}
	return shp.StringField(def.name, uint8(def.length))
}

func getFieldDefs(shapes []jShape) []fieldDef {
	m := make(map[string]fieldDef)
	for _, shape := range shapes {
		for key, val := range shape.Data {
			switch t := val.(type) {
			case float64:
				if _, in := m[key]; !in {
					fmt.Println("  found float field", key)
					m[key] = fieldDef{name: key, numeric: true}
				}
			case string:
				s := val.(string)
				if f, in := m[key]; in {
					if len(s) > f.length {
						f.length = len(s)
					}
				} else {
					fmt.Println("  found text field", key)
					m[key] = fieldDef{name: key, numeric: false, length: len(s)}
				}
			default:
				fmt.Printf("unexpected type %T for field %s\n", t, key)
			}
		}
	}
	defs := make([]fieldDef, 0, len(m))
	for _, def := range m {
		defs = append(defs, def)
	}
	return defs
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("No JSON file given")
		fmt.Println("Usage: shapef <json file>")
		return
	}

	jsonFile := os.Args[1]
	fmt.Println("\n  Convert", jsonFile)

	shapes := readJSON(jsonFile)
	fmt.Println("  Found", len(shapes), "polygon definitions")
	fieldDefs := getFieldDefs(shapes)
	fields := make([]shp.Field, 0, len(fieldDefs))
	for _, def := range fieldDefs {
		fields = append(fields, def.Field())
	}

	shpFile := strings.TrimSuffix(jsonFile, ".json") + ".shp"
	fmt.Println("  Create", shpFile)
	shapeFile, err := shp.Create(shpFile, shp.POLYGON)
	check(err, "failed to create shapefile")
	defer shapeFile.Close()

	shapeFile.SetFields(fields)
	for shapeIdx, shape := range shapes {
		shapeFile.Write(shape.Polygon())
		for fieldIdx, fieldDef := range fieldDefs {
			shapeFile.WriteAttribute(shapeIdx, fieldIdx, shape.Data[fieldDef.name])
		}
	}
	fmt.Println("  All done")
}
