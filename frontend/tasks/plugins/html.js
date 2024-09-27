const HtmlWebpackPlugin = require("html-webpack-plugin");

module.exports = (env) => {
  const defaultConfig = new HtmlWebpackPlugin({
    hash: true,
    fileName: "index.html",
    template: "./src/index.html",
  });

  const plugin = {
    production: defaultConfig,
    development: defaultConfig,
  };

  return plugin[env];
};
