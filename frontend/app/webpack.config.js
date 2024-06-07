const ForkTsCheckerWebpackPlugin = require('fork-ts-checker-webpack-plugin');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const TsconfigPathsPlugin = require('tsconfig-paths-webpack-plugin');
const path = require('path');
const Dotenv = require('dotenv-webpack');
const SOURCE_PATH = path.resolve(__dirname, 'src/app/');
const PUBLIC_PATH = path.resolve(__dirname, 'public');

module.exports = {
  entry: path.join(SOURCE_PATH, 'app.tsx'),
  mode: process.env.NODE_ENV || 'development',

  output: {
    filename: 'bundle.js',
    path: path.resolve(__dirname, 'dist'),
    clean: true,
  },

  resolve: {
    extensions: ['.ts', '.tsx', '.js', '.json'],
    plugins: [new TsconfigPathsPlugin()],
  },

  module: {
    rules: [
      {
        test: /\.(ts|tsx)$/,
        exclude: /node_modules/,
        use: ['ts-loader'],
      },
      {
        test: /\.s?css$/,
        exclude: /\.module\.scss$/,
        use: ['style-loader', 'css-loader', 'sass-loader'],
      },
      {
        test: /\.module\.scss$/,
        use: [
          'style-loader',
          {
            loader: 'css-loader',
            options: {
              modules: {
                exportLocalsConvention: 'dashes',
                localIdentName: '[local]___[hash:base64:5]',
              },
            },
          },
          'sass-loader',
        ],
      },
      {
        test: /\.(jpg|jpeg|png|gif|mp3|svg)$/,
        use: ['file-loader'],
      },
      {
        test: /\.woff(2)?$/,
        use: [
          {
            loader: 'url-loader',
            options: {
              limit: 10000,
              name: './font/[hash].[ext]',
              mimetype: 'application/font-woff',
            },
          },
        ],
      },
      {
        test: /\.(woff|woff2)$/,
        use: [{
          loader: 'file-loader',
          options: {
            name: '[name].[.ext]',
            outputPath: 'fonts/',
          },
        },
        ],
      },
    ],
  },

  plugins: [
    new HtmlWebpackPlugin({
      inject: true,
      template: path.join(PUBLIC_PATH, 'index.html'),
      manifest: path.join(PUBLIC_PATH, 'manifest.json'),
    }),
    new MiniCssExtractPlugin(),
    new ForkTsCheckerWebpackPlugin(),
    new Dotenv(),
  ],

  devtool: 'inline-source-map',

  performance: {
    hints: false,
  },

  devServer: {
    https: true,
    port: 8080,
    hot: true,
    compress: true,
    open: true,
    host: 'kubernetes.docker.internal',

    proxy: {
      '/browserkube/': {
        changeOrigin: true,
        target: 'https://kubernetes.docker.internal',
        secure: false,
        logLevel: 'debug',
      },
    }
  },
};
