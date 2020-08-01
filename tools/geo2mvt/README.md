# geo2mvt

A Go program that converts a [GeoJSON](https://geojson.org) object into a [Mapbox vector tile](https://docs.mapbox.com/vector-tiles/specification/) (mvt) object.

This tool is powered by [orb](https://github.com/paulmach/orb).

## Prerequisites

You need the following software installed,
- [Go](https://golang.org) v1.13 or later

## Building the tool

1. Run `go build`.

    ```
    go build
    ```

2. You will find `geo2mvt` in this directory.

## Running the tool

```
geo2mvt -x X -y Y -z ZOOM INPUT OUTPUT
```

Parameters
- `X`: x coordinate of a tile to be generated.
- `Y`: y coordinate of a tile to be generated.
- `ZOOM`: zoom of a tile to be generated.
- `INPUT`: path to an input JSON file.
- `OUTPUT`: path to an output map vector tile file, which is a zipped PBF file.

As a Mapbox's map vector tile can contain multiple layers in a single tile, an input JSON file, multilayer GeoJSON file, can contain multiple GeoJSON objects inside.
It associates a layer name with a corresponding GeoJSON object.
A multilayer GeoJSON file looks like the following,

```js
{
    "layers-1": {
        // GeoJSON object
    },
    "layers-2": {
        // GeoJSON object
    }, ...
}
```