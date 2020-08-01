package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kikuomax/imaginary-map/cdn/common"
	"github.com/paulmach/orb/encoding/mvt"
	"github.com/paulmach/orb/simplify"
)

// Environment variable names.
const GEO_JSON_BUCKET_NAME = "GEO_JSON_BUCKET_NAME"
const PAPERS_GEO_JSON_VERSION = "PAPERS_GEO_JSON_VERSION"

// GeoJSON file name.
const GEO_JSON_FILE_NAME = "papers.json"

// MVT layer name.
const LAYER_NAME = "papers"

// Loads a multilayer GeoJson file for papers.
func LoadLayersJson () (*common.NamedFeatureCollections, error) {
	bucket, err := common.GetEnv(GEO_JSON_BUCKET_NAME)
	if err != nil {
		return nil, err
	}
	version, err := common.GetEnv(PAPERS_GEO_JSON_VERSION)
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("/%s/%s", version, GEO_JSON_FILE_NAME)
	log.Println("loading a multilayer GeoJSON file", bucket, key)
	bytes, err := common.LoadS3Object(bucket, key)
	if err != nil {
		return nil, err
	}
	return common.LoadLayersJson(bytes)
}

// Generates a map vector tile at a given coordinate.
func GenerateMapVectorTile (fcs *common.NamedFeatureCollections, event common.GetTileEvent) ([]byte, error) {
	layers := mvt.NewLayers(*fcs)
	tile, ok := event.ToTile()
	if !ok {
		return nil, errors.New(
			fmt.Sprintf("invalid tile coordinate: %v", event))
	}
	layers.ProjectToTile(tile)
	layers.Clip(mvt.MapboxGLDefaultExtentBound)
	layers.Simplify(simplify.DouglasPeucker(1.0))
	layers.RemoveEmpty(1.0, 1.0)
	return mvt.MarshalGzipped(layers)
}

func HandleRequest (ctx context.Context, event common.GetTileEvent) ([]byte, error) {
	log.Printf(
		"getting tile at x=%v, y=%v, zoom=%v",
		event.X,
		event.Y,
		event.Zoom)
	fcs, err := LoadLayersJson()
	if err != nil {
		return nil, err
	}
	return GenerateMapVectorTile(fcs, event)
}

func main () {
	lambda.Start(HandleRequest)
}
