# Example GitHub Action Using Golang

[![CI](https://github.com/pplmx/setup-custom-action-by-docker/workflows/CI/badge.svg)](https://github.com/pplmx/setup-custom-action-by-docker/actions)
[![Coverage Status](https://coveralls.io/repos/github/pplmx/setup-custom-action-by-docker/badge.svg?branch=main)](https://coveralls.io/github/pplmx/setup-custom-action-by-docker?branch=main)

## Overview

**Example GitHub Action Using Golang** is a custom GitHub Action designed to handle text processing, lists, files, and
API requests efficiently using Go. The action is configured using a TOML file and runs inside a Docker container.

## Features

- **Text Processing**: Perform find-and-replace operations on provided text.
- **List Operations**: Calculate word count, sum, and average for provided data.
- **API Requests**: Fetch data from APIs and extract specific fields.
- **Configuration**: Fully configurable via a TOML file.

## Inputs

| Name          | Description                         | Required | Default                                |
|---------------|-------------------------------------|----------|----------------------------------------|
| `config_path` | Path to the TOML configuration file | No       | `.github/configs/setup-custom-action-by-docker.toml` |

## Outputs

| Name             | Description                               |
|------------------|-------------------------------------------|
| `processed_text` | The processed text after find and replace |
| `word_count`     | The total number of words in the text     |
| `sum`            | The sum of the numbers                    |
| `average`        | The average of the numbers                |
| `response_field` | A specific field from the API response    |

## Usage

To use this GitHub Action in your workflow, you can define a step that utilizes it. Here's an example:

```yaml
name: Example Workflow
on: [ push ]

jobs:
    example-job:
        runs-on: ubuntu-latest
        steps:
            -   uses: actions/checkout@v4

            -   name: Run Example GitHub Action Using Golang
                uses: pplmx/setup-custom-action-by-docker@v1
```

You can also follow [this](.github/workflows/test_custom_action_itself.yml).

## Local Development

For instructions on setting up a local development environment and testing this GitHub Action, please refer to the [Development Documentation](docs/development.md).

## License

This project is dual-licensed under either of the following licenses, at your option:

- Apache License, Version 2.0 ([LICENSE-APACHE](LICENSE-APACHE) or [http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0))
- MIT license ([LICENSE-MIT](LICENSE-MIT) or [http://opensource.org/licenses/MIT](http://opensource.org/licenses/MIT))

## Contribution

Unless you explicitly state otherwise, any contribution intentionally submitted for inclusion in the work by you, as defined in the Apache-2.0 license, shall be dual licensed as above, without any additional terms or conditions.

For more details, see [CONTRIBUTING.md](CONTRIBUTING.md).
