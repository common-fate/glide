# Logging

Our Lambda functions are all configured to send logs to [CloudWatch](https://aws.amazon.com/cloudwatch/).

In development you can use the `gdeploy` CLI tool to easily view logs for your deployment.

This is also useful for debugging a staging or production deployment.

## Prerequisites

You'll need to build the `gdeploy` first:

```
make gdeploy
```

## Viewing logs

To view logs for your own deployment, run:

```bash
gdeploy logs get
```

To view logs for a particular deployment stage, run:

```bash
gdeploy logs get -s my-stage
```

## Streaming live logs

You can also stream live logs for a particular deployment.

To stream logs for your own deployment:

```bash
gdeploy logs watch
```

To stream logs for another deployment:

```bash
gdeploy logs watch -s my-stage
```
