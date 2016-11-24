# Datadog Agent

[![CircleCI](https://circleci.com/gh/DataDog/datadog-agent/tree/master.svg?style=svg&circle-token=dbcee3f02b9c3fe5f142bfc5ecb735fdec34b643)](https://circleci.com/gh/DataDog/datadog-agent/tree/master)

The Datadog Agent faithfully collects events and metrics and brings
them to [Datadog](https://app.datadoghq.com) on your behalf so that
you can do something useful with your monitoring and performance data.

## Requirements
To build the project you need:
 * `go` 1.6+
 * `rake`
 * an `agent` version 5.x installed under `/opt/`

 We use `pkg-config` to make compilers and linkers aware of CPython. If you need to adjust the build for your specific configuration, add or edit the files within the `pkg-config` folder.

## Getting started
Binary distributions are not provided yet, to try out the Agent you can build the `master` branch. Checkout the repo within your `GOPATH`, then install `glide`:
```
go get github.com/Masterminds/glide
```

Use `glide` to fetch project dependencies:
```
glide up
```

To run the test suite `golint` has to be available on your system, if it is not, just install it with:
```
go get -u github.com/golang/lint/golint
```

Build and tests are orchestrated by a `Rakefile`, write `rake -T` on a shell to see the available tasks.
If you're using the DogBox, ask `gimme` to provide a recent version of go, like:
```
eval "$(gimme 1.6.2)"
```

## Executing
Just execute the `agent` launch script within the `bin/agent` folder, it will take care of adjusting paths and run the binary.
You need to provide a valid API key, either through the config file or passing the environment variable like:
```
DD_API_KEY=12345678990 ./bin/agent/agent
```
