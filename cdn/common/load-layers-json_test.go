package common

import (
	"reflect"
	"testing"
	"github.com/paulmach/orb"
)

const goodJson =
`{
  "papers": {
    "type": "FeatureCollection",
    "features": [
      {
        "type": "Feature",
        "geometry": {
          "type": "Point",
          "coordinates": [
            1,
            2
          ]
        }
      }
    ]
  }
}`

const goodJson2 =
`{
  "islands-0": {
    "type": "FeatureCollection",
    "features": [
      {
        "type": "Feature",
        "geometry": {
          "type": "Polygon",
          "coordinates": [
            [
              [1, 0],
              [1, 1],
              [0, 1],
              [1, 0]
            ]
          ]
        }
      }
    ]
  },
  "islands-1": {
    "type": "FeatureCollection",
    "features": [
      {
        "type": "Feature",
        "geometry": {
          "type": "Polygon",
          "coordinates": [
            [
              [-0.5, 1],
              [0.5, -0.5],
              [2, -0.5],
              [1.5, 2],
              [-0.5, 1]
            ]
          ]
        }
      }
    ]
  },
  "islands-2": {
    "type": "FeatureCollection",
    "features": [
      {
        "type": "Feature",
        "geometry": {
          "type": "Polygon",
          "coordinates": [
            [
              [1, -1],
              [1, 1],
              [-1, 1],
              [-1, -1],
              [1, -1]
            ]
          ]
        }
      }
    ]
  }
}`

const nonJson = `"hanage"`

const badJsonArray = `[1, 2, 3]`

const badJsonNonFeatureCollection =
`{
  "non-compliant": {
    "Name": "John Flimsy",
	"Message": "I'm not that compliant."
  }
}`

func TestLoadLayersJson (t *testing.T) {
	t.Run(
		"LoadLayersJson should correctly load the test JSON input goodJson",
		func (t *testing.T) {
			layers, err := LoadLayersJson([]byte(goodJson))
			if err != nil {
				t.Fatalf(
					"Expected LoadLayersJson(goodJson) not to return an error but got %v",
					err,
				)
			}
			papers, ok := (*layers)["papers"]
			if !ok {
				t.Fatal(`Expected LoadLayersJson(goodJson)["papers"] exist but not`)
			}
			if papers.Type != "FeatureCollection" {
				t.Errorf(
					`Expected LoadLayersJson(goodJson)["papers"].Type to be "FeatureCollection" but got %v`,
					papers.Type,
				)
			}
			if len(papers.Features) != 1 {
				t.Fatalf(
					`Expected LoadLayersJson(goodJson)["papers"].Features to be length of 1 but got %v`,
					len(papers.Features),
				)
			}
			if papers.Features[0].Type != "Feature" {
				t.Fatalf(
					`Expected LoadLayersJson(goodJson)["papers"].Features[0].Type to be "Feature" but got %v`,
					papers.Features[0].Type,
				)
			}
			switch geo := papers.Features[0].Geometry.(type) {
			case orb.Point:
				expected := orb.Point{1, 2}
				if !reflect.DeepEqual(geo, expected) {
					t.Errorf(
						`Expected LoadLayersJson(goodJson)["papers"].Features[0].Geometry to equal to %v but got %v`,
						expected,
						geo,
					)
				}
			default:
				t.Fatalf(
					`Expected LoadLayersJson(goodJson)["papers"].Features[0].Geometry to be orb.Point but got %v`,
					papers.Features[0].Geometry.GeoJSONType(),
				)
			}
		},
	)
	t.Run(
		"LoadLayerJson should correctly load the test JSON input goodJson2",
		func (t *testing.T) {
			layers, err := LoadLayersJson([]byte(goodJson2))
			if err != nil {
				t.Fatalf(
					`Expected LoadLayersJson(goodJson2) not to return an error but got %v`,
					err,
				)
			}
			// islands-0
			islands, ok := (*layers)["islands-0"]
			if !ok {
				t.Fatal(`Expected LoadLayersJson(goodJson2)["islands-0"] to exist but not`)
			}
			if islands.Type != "FeatureCollection" {
				t.Errorf(
					`Expected LoadLayersJson(goodJson2)["islands-0"].Type to be "FeatureCollection" but got %v`,
					islands.Type,
				)
			}
			if len(islands.Features) != 1 {
				t.Fatalf(
					`Expected LoadLayersJson(goodJson2)["islands-0"].Features to be length of 1 but got %v`,
					len(islands.Features),
				)
			}
			if islands.Features[0].Type != "Feature" {
				t.Fatalf(
					`Expected LoadLayersJson(goodJson2)["islands-0"].Features[0].Type to be "Feature" but got %v`,
					islands.Features[0].Type,
				)
			}
			switch geo := islands.Features[0].Geometry.(type) {
			case orb.Polygon:
				expected := orb.Polygon{ orb.Ring{
					{1, 0},
					{1, 1},
					{0, 1},
					{1, 0},
				} }
				if !reflect.DeepEqual(geo, expected) {
					t.Errorf(
						`Expected LoadLayersJson(goodJson2)["islands-0"].Features[0].Geometry to equal to %v but got %v`,
						expected,
						geo,
					)
				}
			default:
				t.Fatalf(
					`Expected LoadLayersJson(goodJson2)["islands-0"].Features[0].Geometry to be orb.Polygon but got %v`,
					islands.Features[0].Geometry.GeoJSONType(),
				)
			}
			// islands-1
			islands, ok = (*layers)["islands-1"]
			if !ok {
				t.Fatal(`Expected LoadLayersJson(goodJson2)["islands-1"] to exist but not`)
			}
			if islands.Type != "FeatureCollection" {
				t.Errorf(
					`Expected LoadLayersJson(goodJson2)["islands-1"].Type to be "FeatureCollection" but got %v`,
					islands.Type,
				)
			}
			if len(islands.Features) != 1 {
				t.Fatalf(
					`Expected LoadLayersJson(goodJson2)["islands-1"].Features to be length of 1 but got %v`,
					len(islands.Features),
				)
			}
			if islands.Features[0].Type != "Feature" {
				t.Fatalf(
					`Expected LoadLayersJson(goodJson2)["islands-1"].Features[0].Type to be "Feature" but got %v`,
					islands.Features[0].Type,
				)
			}
			switch geo := islands.Features[0].Geometry.(type) {
			case orb.Polygon:
				expected := orb.Polygon{ orb.Ring{
					{-0.5, 1},
					{0.5, -0.5},
					{2, -0.5},
					{1.5, 2},
					{-0.5, 1},
				} }
				if !reflect.DeepEqual(geo, expected) {
					t.Errorf(
						`Expected LoadLayersJson(goodJson2)["islands-1"].Features[0].Geometry to equal to %v but got %v`,
						expected,
						geo,
					)
				}
			default:
				t.Fatalf(
					`Expected LoadLayersJson(goodJson2)["islands-1"].Features[0].Geometry to be orb.Polygon but got %v`,
					islands.Features[0].Geometry.GeoJSONType(),
				)
			}
			// islands-2
			islands, ok = (*layers)["islands-2"]
			if !ok {
				t.Fatal(`Expected LoadLayersJson(goodJson2)["islands-2"] to exist but not`)
			}
			if islands.Type != "FeatureCollection" {
				t.Errorf(
					`Expected LoadLayersJson(goodJson2)["islands-2"].Type to be "FeatureCollection" but got %v`,
					islands.Type,
				)
			}
			if len(islands.Features) != 1 {
				t.Fatalf(
					`Expected LoadLayersJson(goodJson2)["islands-2"].Features to be length of 1 but got %v`,
					len(islands.Features),
				)
			}
			if islands.Features[0].Type != "Feature" {
				t.Fatalf(
					`Expected LoadLayersJson(goodJson2)["islands-2"].Features[0].Type to be "Feature" but got %v`,
					islands.Features[0].Type,
				)
			}
			switch geo := islands.Features[0].Geometry.(type) {
			case orb.Polygon:
				expected := orb.Polygon{ orb.Ring{
					{1, -1},
					{1, 1},
					{-1, 1},
					{-1, -1},
					{1, -1},
				} }
				if !reflect.DeepEqual(geo, expected) {
					t.Errorf(
						`Expected LoadLayersJson(goodJson2)["islands-2"].Features[0].Geometry to equal to %v but got %v`,
						expected,
						geo,
					)
				}
			default:
				t.Fatalf(
					`Expected LoadLayersJson(goodJson2)["islands-2"].Features[0].Geometry to be orb.Polygon but got %v`,
					islands.Features[0].Geometry.GeoJSONType(),
				)
			}
		},
	)
	t.Run(
		"LoadLayersJson should return an error for a non-JSON input (nonJson)",
		func (t *testing.T) {
			_, err := LoadLayersJson([]byte(nonJson))
			if err == nil {
				t.Error("Expected LoadLayersJson(nonJson) to return an error but not")
			}
		},
	)
	t.Run(
		"LoadLayersJson should return an error for a bad JSON input that is an array (badJsonArray)",
		func (t *testing.T) {
			_, err := LoadLayersJson([]byte(badJsonArray))
			if err == nil {
				t.Error("Expected LoadLayersJson(badJsonArray) to return an error but not")
			}
		},
	)
	t.Run(
		"LoadLayersJson should return an error for a bad JSON input that contains a non-FeatureCollection (badJsonNonFeatureCollection)",
		func (t *testing.T) {
			_, err := LoadLayersJson([]byte(badJsonNonFeatureCollection))
			if err == nil {
				t.Error("Expected LoadLayersJson(badJsonNonFeatureCollection) to return an error but not")
			}
		},
	)
}
