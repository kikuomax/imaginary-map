const path = require('path')

const HtmlWebpackPlugin = require('html-webpack-plugin')

const defaultMode = 'development'

module.exports = {
  mode: defaultMode,
  entry: './src/index.js',
  output: {
    filename: 'main.js',
    path: path.resolve(__dirname, 'dist')
  },
  module: {
    rules: [
      {
        test: /\.css$/,
        use: [
          'style-loader',
          'css-loader'
        ]
      }
    ]
  },
  plugins: [
    new HtmlWebpackPlugin({
      title: 'Imaginary Map',
      template: path.resolve(__dirname, './src/index.ejs')
    })
  ]
}
