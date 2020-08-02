package common

import (
	"encoding/json"
	"github.com/paulmach/orb/geojson"
)

type NamedFeatureCollections map[string]*geojson.FeatureCollection

func LoadLayersJson (bytes []byte) (*NamedFeatureCollections, error) {
	layers := make(NamedFeatureCollections)
	rawLayers := make(rawNamedFeatureCollections)
	err := json.Unmarshal(bytes, &rawLayers)
	if err != nil {
		return nil, err
	}
	for name, raw := range rawLayers {
		layers[name], err = geojson.UnmarshalFeatureCollection(raw)
		if err != nil {
			return nil, err
		}
	}
	return &layers, nil
}

type rawNamedFeatureCollections map[string]json.RawMessage
