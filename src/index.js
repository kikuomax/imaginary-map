/**
 * An entry point.
 */

import mapboxgl from 'mapbox-gl'

import 'mapbox-gl/dist/mapbox-gl.css'

console.log(window.location.origin)
const urlParams = new URLSearchParams(window.location.search)
const accessToken = urlParams.get('access_token')
mapboxgl.accessToken = accessToken
const map = new mapboxgl.Map({
  container: 'map',
  style: 'mapbox://styles/mapbox/streets-v11',
  center: [-74.5, 40],
  zoom: 9,
  // renderWorldCopies:false makes only one map rendered
  // https://github.com/mapbox/mapbox-gl-js/pull/3885
  renderWorldCopies: false
})
map.on('load', () => {
  map.addSource('imaginary', {
    type: 'raster',
    tiles: [
      `./tiles/{z}/{x}/{y}.png`
    ],
    minzoom: 0,
    maxzoom: 0
  })
  map.addLayer(
    {
      id: 'imaginary',
      type: 'raster',
      source: 'imaginary'
    }
  )
})
