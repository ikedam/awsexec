# awsexec: execute command with aws configure export-credentials

## Background

* AWS CLI provides `aws configure export-credentials` command, which exports credentials in a format that can be set as environment variables. By executing the authentication process with AWS CLI and passing only the resulting credentials to applications, you can use authentication processes that are not directly supported by applications.
    * An example is the `aws login` command supported from AWS CLI v2.32.0. Applications using older AWS SDKs cannot directly use credentials from `aws login`, so credentials need to be passed to applications via the `aws configure export-credentials` command.
        * https://docs.aws.amazon.com/signin/latest/userguide/command-line-sign-in.html
        * https://docs.aws.amazon.com/cli/latest/reference/login/

* Passing the output of the `aws configure export-credentials` command to applications requires multiple steps. These steps are often tedious and require some basic knowledge of shell scripting.

    1. Execute `aws configure export-credentials`.
    2. Set the output credentials as environment variables.
    3. Execute the target command.

* `awsexec` is a very simple wrapper program that automates these steps.

## awsexec Implementation and Distribution Policy

* Two distribution formats are provided.
    1. Single binary built with golang ([GitHub Releases](https://github.com/ikedam/awsexec/releases))
    2. Shell script ([awsexec.sh](awsexec.sh))

* Since this is a security-related feature, no external libraries are used.
    * This is to avoid the need to check the safety of dependency libraries.
* Pre-built binaries are provided as GitHub Releases Assets, but you can also build them in your own environment if Docker is installed. Also, if using the shell script, you can create it in your own environment by copying and pasting.

## Prerequisites

* AWS CLI must be installed:

    ```sh
    aws --version
    ```

    * If not installed, please install it according to AWS installation instructions.
        * https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html

    * AWS CLI v2.32.0 or later is recommended. The `aws login` command is available.
        * https://docs.aws.amazon.com/signin/latest/userguide/command-line-sign-in.html

## Installation

### Binary Installation

* Download the awsexec binary appropriate for your environment from [GitHub Releases](https://github.com/ikedam/awsexec/releases) and place it in a location where the path is available.
* For example, you can install it with the following steps:

    1. Download the archive containing the binary:

        ```sh
        curl -L -o /tmp/awsexec.tar.gz https://github.com/ikedam/awsexec/releases/download/latest/awsexec-linux-amd64.tar.gz
        ```

    2. Extract the archive:

        ```sh
        tar -xzf /tmp/awsexec.tar.gz -C /tmp
        ```

    3. Place the binary:

        ```sh
        sudo cp /tmp/awsexec /usr/local/bin/awsexec
        ```

### Shell Script Installation

* If a shell script environment is available, you can also place the shell script in your own environment by copying and pasting.
    * Usually available on Linux and macOS.
    * Only available in WSL environment on Windows.
* For example, you can install it with the following steps:

    1. Download the shell script:

        ```sh
        curl -L -o /tmp/awsexec.sh https://raw.githubusercontent.com/ikedam/awsexec/refs/heads/main/awsexec.sh
        ```

    2. Place the shell script:

        ```sh
        sudo cp /tmp/awsexec.sh /usr/local/bin/awsexec
        ```

    3. Grant execute permission to the shell script:

        ```sh
        chmod +x /usr/local/bin/awsexec
        ```

## Usage

* To specify a profile:

    ```sh
    awsexec profile -- command
    ```

    * Executes the specified command with the credentials of the specified AWS profile set as environment variables.

* To specify a profile using the `AWS_PROFILE` environment variable:

    ```sh
    AWS_PROFILE=profile awsexec -- command
    ```

    * To avoid using the wrong profile, specifying the `AWS_PROFILE` environment variable is required.
    * If you don't need to use multiple profiles, you can set `AWS_PROFILE=default` in your user profile.

## Building the Binary

You can build the binary using Docker. The built binary will be generated at `build/awsexec`.

* Linux (amd64)

    ```sh
    GOOS=linux GOARCH=amd64 docker compose run --rm build
    ```

* Windows (amd64)

    ```sh
    GOOS=windows GOARCH=amd64 docker compose run --rm build
    ```

* macOS (arm64)

    ```sh
    GOOS=darwin GOARCH=arm64 docker compose run --rm build
    ```

## License

This application and source code are distributed under the [MIT License](LICENSE).
