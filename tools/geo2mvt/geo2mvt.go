package main

import (
	"flag"
	"fmt"
	"github.com/kikuomax/imaginary-map/cdn/common"
	"github.com/paulmach/orb/encoding/mvt"
	"github.com/paulmach/orb/maptile"
	"github.com/paulmach/orb/simplify"
	"io/ioutil"
	"log"
	"os"
)

func LoadGeoJson (geoPath string) (*common.NamedFeatureCollections, error) {
	geoIn, err := os.Open(geoPath)
	if err != nil {
		return nil, err
	}
	defer geoIn.Close()
	jsonBytes, err := ioutil.ReadAll(geoIn)
	if err != nil {
		return nil, err
	}
	fcs, err := common.LoadLayersJson(jsonBytes)
	if err != nil {
		return nil, err
	}
	return fcs, err
}

func SaveMvt (mvtPath string, mvtBytes []byte) error {
	mvtOut, err := os.Create(mvtPath)
	if err != nil {
		return err
	}
	defer mvtOut.Close()
	_, err = mvtOut.Write(mvtBytes)
	if err != nil {
		return err
	}
	return nil
}

func main () {
	x := flag.Uint(
		"x",
		0,
		"x-coordinate value of a tile to be outputted")
	y := flag.Uint(
		"y",
		0,
		"y-coordinate value of a tile to be outputted")
	z := flag.Uint(
		"z",
		0,
		"z-coordinate value (zoom) of a tile to be outputted")
	flag.Parse()
	args := flag.Args()
	geoPath := args[0]
	mvtPath := args[1]
	fmt.Printf("loading: %v\n", geoPath)
	featureCollections, err := LoadGeoJson(geoPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("converting GeoJSON")
	layers := mvt.NewLayers(*featureCollections)
	layers.ProjectToTile(maptile.New(
		uint32(*x),
		uint32(*y),
		maptile.Zoom(*z)))
	layers.Clip(mvt.MapboxGLDefaultExtentBound)
	layers.Simplify(simplify.DouglasPeucker(1.0))
	layers.RemoveEmpty(1.0, 1.0)
	mvtBytes, err := mvt.MarshalGzipped(layers)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("mvt # of bytes: %v\n", len(mvtBytes))
	fmt.Printf("saving: %v\n", mvtPath)
	err = SaveMvt(mvtPath, mvtBytes)
	if err != nil {
		log.Fatal(err)
		return
	}
}
