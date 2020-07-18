package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kikuomax/imaginary-map/cdn/common"
	"github.com/paulmach/orb/encoding/mvt"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/simplify"
)

// Environment variable names.
const GEO_JSON_BUCKET_NAME = "GEO_JSON_BUCKET_NAME"
const PAPERS_GEO_JSON_VERSION = "PAPERS_GEO_JSON_VERSION"

// GeoJSON file name.
const GEO_JSON_FILE_NAME = "papers.json"

// MVT layer name.
const LAYER_NAME = "papers"

// Loads a GeoJson for papers.
func LoadGeoJson () (*geojson.FeatureCollection, error) {
	bucket, err := common.GetEnv(GEO_JSON_BUCKET_NAME)
	if err != nil {
		return nil, err
	}
	version, err := common.GetEnv(PAPERS_GEO_JSON_VERSION)
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("/%s/%s", version, GEO_JSON_FILE_NAME)
	log.Println("loading GeoJSON", bucket, key)
	return common.LoadGeoJsonFromS3(bucket, key)
}

// Generates a map vector tile at a given coordinate.
func GenerateMapVectorTile (fc *geojson.FeatureCollection, event common.GetTileEvent) ([]byte, error) {
	layer := mvt.NewLayer(LAYER_NAME, fc)
	tile, ok := event.ToTile()
	if !ok {
		return nil, errors.New(
			fmt.Sprintf("invalid tile coordinate: %v", event))
	}
	layer.ProjectToTile(tile)
	layer.Clip(mvt.MapboxGLDefaultExtentBound)
	layer.Simplify(simplify.DouglasPeucker(1.0))
	layer.RemoveEmpty(1.0, 1.0)
	return mvt.MarshalGzipped(mvt.Layers{ layer })
}

func HandleRequest (ctx context.Context, event common.GetTileEvent) ([]byte, error) {
	log.Printf(
		"getting tile at x=%v, y=%v, zoom=%v",
		event.X,
		event.Y,
		event.Zoom)
	fc, err := LoadGeoJson()
	if err != nil {
		return nil, err
	}
	return GenerateMapVectorTile(fc, event)
}

func main () {
	lambda.Start(HandleRequest)
}
