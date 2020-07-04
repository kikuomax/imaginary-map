// starts a local server for tests.

const express = require('express')
const expressStaticGzip = require('express-static-gzip')
const path = require('path')

const port = 3000

const app = express()
// tiles are gzipped
app.use(
  '/tiles/',
  expressStaticGzip(path.resolve(__dirname, './docs/tiles')))
app.use(
  '/',
  express.static(path.resolve(__dirname, './docs')))
app.listen(port, () => console.log(`listening at http://localhost:${port}`))
