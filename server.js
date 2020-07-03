// starts a local server for tests.

const express = require('express')
const path = require('path')

const port = 3000

const app = express()
app.get('/', app.use(express.static(path.resolve(__dirname, 'docs'))))
app.listen(port, () => console.log(`listening at http://localhost:${port}`))
