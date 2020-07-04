/**
 * An entry point.
 */

import mapboxgl from 'mapbox-gl'

import 'mapbox-gl/dist/mapbox-gl.css'

// gets directory name from a given path.
function getDirname (path) {
  if (path.endsWith('/')) {
    return path
  } else {
    const slashIndex = path.lastIndexOf('/')
    if (slashIndex !== -1) {
      return path.substring(0, slashIndex + 1)
    } else {
      return path
    }
  }
}

const dirPath = getDirname(window.location.pathname)
const baseUrl = `${window.location.origin}${dirPath}`
if (process.env.NODE_ENV !== 'production') {
  console.log('baseUrl', baseUrl)
}
const urlParams = new URLSearchParams(window.location.search)
const accessToken = urlParams.get('access_token')
mapboxgl.accessToken = accessToken
const map = new mapboxgl.Map({
  container: 'map',
  style: {
    version: 8,
    name: "imaginary",
    sources: {
      imaginary: {
        type: 'vector',
        tiles: [
          `${baseUrl}tiles/{z}/{x}/{y}.pbf`
        ],
        minzoom: 0,
        maxzoom: 1
      }
    },
    layers: [
      {
        id: 'background',
        type: 'background',
        paint: {
          'background-color': '#505050'
        }
      },
      {
        id: 'imaginary',
        type: 'fill',
        source: 'imaginary',
        'source-layer': 'test',
        paint: {
          'fill-color': '#AFDB1C',
          'fill-outline-color': '#181E04'
        }
      }
    ]
  },
  center: [-74.5, 40],
  zoom: 9,
  // renderWorldCopies:false makes only one map rendered
  // https://github.com/mapbox/mapbox-gl-js/pull/3885
  renderWorldCopies: false
})
