# Imaginary Map

A PoC project that renders an imaginary map using a [Mapbox GL JS](https://docs.mapbox.com/mapbox-gl-js/api/) API.

This is a sister project of [COVID-19 research](https://github.com/metasphere-xyz/covid19-research).

## Prerequisites

You need the following software installed,
- [Node.js](https://nodejs.org/en/) (tested with 12.14.0)

## Building an application

Please take the following steps,

1. Install modules.

    ```
    npm ci
    ```

2. Build the application.

    ```
    npm run build
    ```

3. You will find the following files in the `docs` directory updated.
    - `index.html`
    - `main.js`

For production, specify a `--mode=production` option at the step 2.

```
npm run build -- --mode=production
```

## Making a Mapbox vector tiles

There is a Go program that converts a GeoJSON object into a Mapbox vector tile object.
Please refer to [`tools/geo2html`](tools/geo2html) for more information.

## Hosting a CDN for map vector tiles

Map vector tiles are delivered through a CDN powered by [AWS CloudFront](https://aws.amazon.com/cloudfront/).
Please refer to [`cdn`](cdn) for more information.

## Running a test service

The following command starts a test service at http://localhost:3000.

```
npm start
```

You have to specify a CDN URL to a `tile-api` query parameter.
Suppose you have a CDN running at `https://xyz.cloudfront.net`, open the following URL.

```
http://localhost:3000/?tile-api=https%3A%2F%2Fxyz.cloudfront.net
```

## References

- [Mapbox vector tile specification](https://docs.mapbox.com/vector-tiles/specification/)
- [Vector tile specification v2.1](https://github.com/mapbox/vector-tile-spec/tree/master/2.1)
- [Add a third party vector tile source](https://docs.mapbox.com/mapbox-gl-js/example/third-party/)