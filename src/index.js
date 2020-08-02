/**
 * An entry point.
 */

import mapboxgl from 'mapbox-gl'

import 'mapbox-gl/dist/mapbox-gl.css'

const urlParams = new URLSearchParams(window.location.search)
const tileApiUrl = urlParams.get('tile-api')
if (process.env.NODE_ENV !== 'production') {
  console.log('tileApiUrl', tileApiUrl)
}
const map = new mapboxgl.Map({
  container: 'map',
  style: {
    version: 8,
    name: "imaginary",
    sources: {
      islands: {
        type: 'vector',
        tiles: [
          `${tileApiUrl}/{z}/{x}/{y}/islands.pbf`
        ],
        minzoom: 0,
        maxzoom: 10
      },
      papers: {
        type: 'vector',
        tiles: [
          `${tileApiUrl}/{z}/{x}/{y}/papers.pbf`
        ],
        minzoom: 5,
        maxzoom: 10
      }
    },
    layers: [
      {
        id: 'background',
        type: 'background',
        paint: {
          'background-color': '#D1D1D1'
        }
      },
      {
        id: 'islands',
        type: 'fill',
        source: 'islands',
        'source-layer': 'islands',
        paint: {
          'fill-color': '#DCE89C',
          'fill-outline-color': '#000000'
        }
      },
      {
        id: 'islands-1',
        type: 'fill',
        source: 'islands',
        'source-layer': 'islands-1',
        paint: {
          'fill-color': '#E0D5A6',
          'fill-outline-color': '#A0A0A0'
        }
      },
      {
        id: 'islands-2',
        type: 'fill',
        source: 'islands',
        'source-layer': 'islands-2',
        paint: {
          'fill-color': '#B5A598',
          'fill-outline-color': '#A0A0A0'
        }
      },
      {
        id: 'papers',
        type: 'circle',
        source: 'papers',
        'source-layer': 'papers',
        paint: {
          'circle-opacity': 0.5
        }
      }
    ]
  },
  center: [0, 0],
  zoom: 5,
  // renderWorldCopies:false makes only one map rendered
  // https://github.com/mapbox/mapbox-gl-js/pull/3885
  renderWorldCopies: false
})
