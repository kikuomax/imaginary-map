/**
 * An entry point.
 */

import mapboxgl from 'mapbox-gl'

import 'mapbox-gl/dist/mapbox-gl.css'

const urlParams = new URLSearchParams(window.location.search)
const accessToken = urlParams.get('access_token')
mapboxgl.accessToken = accessToken
const map = new mapboxgl.Map({
  container: 'map',
  style: 'mapbox://styles/mapbox/streets-v11',
  center: [-74.5, 40],
  zoom: 9
})
