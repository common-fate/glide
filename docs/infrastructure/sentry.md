# Sentry

We use [Sentry](https://sentry.io) to track errors in our internal testing deployments of Common Fate.

If you wish to enable Sentry locally follow these steps:

- Set the `USE_SENTRY` environment variable to true when building the NextJS frontend.
- Set the `COMMONFATE_SENTRY_DSN` environment variable in .env for the backend.
