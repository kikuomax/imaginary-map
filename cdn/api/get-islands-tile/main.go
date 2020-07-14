package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/paulmach/orb/encoding/mvt"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/maptile"
	"github.com/paulmach/orb/simplify"
)

// Environment variable names.
const GEO_JSON_BUCKET_NAME = "GEO_JSON_BUCKET_NAME"
const ISLANDS_GEO_JSON_VERSION = "ISLANDS_GEO_JSON_VERSION"

// GeoJSON file name.
const GEO_JSON_FILE_NAME = "islands.json"

// Common form of a GetTileEvent.
type GetTileEvent struct {
	Zoom int `json:"zoom"`
	X int `json:"x"`
	Y int `json:"y"`
}

// Converts a GetTileEvent into a maptile.Tile.
func (event GetTileEvent) ToTile () maptile.Tile {
	return maptile.New(
		uint32(event.X),
		uint32(event.Y),
		maptile.Zoom(event.Zoom),
	)
}

// Loads bytes of the GeoJSON for islands.
func LoadGeoJsonBytes () ([]byte, error) {
	bucketName, ok := os.LookupEnv(GEO_JSON_BUCKET_NAME)
	if !ok {
		return nil, errors.New(
			fmt.Sprintf(
				"environment variable %s is not set",
				GEO_JSON_BUCKET_NAME))
	}
	version, ok := os.LookupEnv(ISLANDS_GEO_JSON_VERSION)
	if !ok {
		return nil, errors.New(
			fmt.Sprintf(
				"environment variable %s is not set",
				ISLANDS_GEO_JSON_VERSION))
	}
	key := fmt.Sprintf("/%s/%s", version, GEO_JSON_FILE_NAME)
	objectIn := s3.GetObjectInput {
		Bucket: aws.String(bucketName),
		Key: aws.String(key),
	}
	log.Println("loading GeoJSON", bucketName, key)
	client := s3.New(session.New())
	objectOut, err := client.GetObject(&objectIn)
	if err != nil {
		return nil, err
	}
	bodyIn := objectOut.Body
	defer bodyIn.Close()
	body, err := ioutil.ReadAll(bodyIn)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// Loads the GeoJSON for islands.
func LoadGeoJson () (*geojson.FeatureCollection, error) {
	bytes, err := LoadGeoJsonBytes()
	if err != nil {
		return nil, err
	}
	return geojson.UnmarshalFeatureCollection(bytes)
}

// Generates a map tile vector at a given coordinate.
func GenerateMapTileVector (fc *geojson.FeatureCollection, event GetTileEvent) ([]byte, error) {
	layer := mvt.NewLayer("islands", fc)
	layer.ProjectToTile(event.ToTile())
	layer.Clip(mvt.MapboxGLDefaultExtentBound)
	layer.Simplify(simplify.DouglasPeucker(1.0))
	layer.RemoveEmpty(1.0, 1.0)
	return mvt.MarshalGzipped(mvt.Layers{ layer })
}

func HandleRequest (ctx context.Context, event GetTileEvent) ([]byte, error) {
	log.Printf(
		"getting tile at x=%v, y=%v, zoom=%v",
		event.X,
		event.Y,
		event.Zoom)
	fc, err := LoadGeoJson()
	if err != nil {
		return nil, err
	}
	return GenerateMapTileVector(fc, event)
}

func main () {
	lambda.Start(HandleRequest)
}
