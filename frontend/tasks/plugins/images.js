const CopyWebpackPlugin = require("copy-webpack-plugin");

module.exports = (env) => {
  const devConfig = new CopyWebpackPlugin({
    patterns: [
      {
        from: "./src/assets/images",
        to: "./images",
      },
    ],
  });
  const prodConfig = new CopyWebpackPlugin({
    patterns: [
      {
        from: "./src/assets/images",
        to: "./src/assets/images",
      },
    ],
  });

  const plugin = {
    production: prodConfig,
    development: devConfig,
  };

  return plugin[env];
};
