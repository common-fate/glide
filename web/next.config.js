const withPlugins = require("next-compose-plugins");
const { PHASE_PRODUCTION_BUILD } = require("next/constants");
const withBundleAnalyzer = require("@next/bundle-analyzer")({
  enabled: process.env.ANALYZE === "true",
});
const { withSentryConfig } = require("@sentry/nextjs");

/** @type {import('@sentry/nextjs').SentryWebpackPluginOptions} */
const sentryWebpackPluginOptions = {
  // Additional config options for the Sentry Webpack plugin. Keep in mind that
  // the following options are set automatically, and overriding them is not
  // recommended:
  //   release, url, org, project, authToken, configFile, stripPrefix,
  //   urlPrefix, include, ignore

  silent: true, // Suppresses all logs
  // For all available options, see:
  // https://github.com/getsentry/sentry-webpack-plugin#options.
};

/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  productionBrowserSourceMaps: true,
};

const plugins = [withBundleAnalyzer];

// only add Sentry in CI builds. This avoids breaking local builds if you just
// want to ensure NextJS is building properly.
if (process.env.USE_SENTRY === "true") {
  console.log("NextJS Sentry plugin enabled");
  plugins.push([
    withSentryConfig,
    sentryWebpackPluginOptions,
    [PHASE_PRODUCTION_BUILD],
  ]);
}

module.exports = withPlugins([...plugins, nextConfig]);
