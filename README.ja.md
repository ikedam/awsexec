# awsexec: execute command with aws configure export-credentials

## 背景

* awscli には `aws configure export-credentials` コマンドにより認証情報のパラメーターを環境変数に与える形式で出力する機能があります。認証プロセスを awscli で実行し、結果の認証情報だけをアプリケーションに渡すことで、アプリケーションでサポートしていない認証プロセスでもアプリケーションと併用することができます。
    * awscli v2.32.0 からサポートされる `aws login` コマンドが挙げられます。古い aws SDK を用いたアプリケーションでは `aws login` による認証情報を直接利用できないため、 `aws configure export-credentials` コマンドを経由して認証情報をアプリケーションに渡す必要があります。
        * https://docs.aws.amazon.com/ja_jp/signin/latest/userguide/command-line-sign-in.html
        * https://docs.aws.amazon.com/cli/latest/reference/login/

* `aws configure export-credentials` コマンドの出力をアプリケーションに渡すためには複数の手順が必要になります。この手順はしばしば手間がかかり、また、シェルにある程度の基礎知識が必要です。

    1. `aws configure export-credentials` を実行する。
    2. 出力された認証情報を環境変数に移す。
    3. 目的のコマンドを実行する。

* `awsexec` はこの手順を一括で実行できるようにするための非常に単純なラッパープログラムです。


## awsexec の実装・提供方針

* 2 つの配布形態を提供しています。
    1. golang によるシングルバイナリー ([GitHub Releases](https://github.com/ikedam/awsexec/releases))
    2. シェルスクリプト ([awsexec.sh](awsexec.sh))


* セキュリティに関わる機能であるため、外部ライブラリーは使用していません。
    * 依存ライブラリーの安全性のチェックなどの考慮を不要にするため。
* GitHub Releases の Assets としてビルド済みバイナリーを提供していますが、 Docker がインストールされている環境であればご自分の環境でビルドすることもできます。また、シェルスクリプトであればコピペでご自分の環境に作成することもできます。

## 前提条件

* awscli がインストールされていること:

    ```sh
    aws --version
    ```

    * インストールされていない場合は AWS のインストール手順にしたがってインストールしてください。
        * https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html

    * awscli v2.32.0 以降のインストールを推奨します。 `aws login` コマンドが利用できます。
        * https://docs.aws.amazon.com/ja_jp/signin/latest/userguide/command-line-sign-in.html


## インストール方法

### バイナリーでのインストール方法 (ダウンロード)

* ご自分の環境に合わせた awsexec のバイナリーを [GitHub Releases](https://github.com/ikedam/awsexec/releases) からダウンロードして、パスが通った場所に配置してください。
* 例えば以下の手順でインストールできます:

    1. バイナリーファイルをダウンロードする:

        ```sh
        curl -L -o /tmp/awsexec https://github.com/ikedam/awsexec/releases/download/latest/awsexec_darwin_arm64
        ```

    2. バイナリーを配置する:

        ```sh
        sudo cp /tmp/awsexec /usr/local/bin/awsexec
        ```

    3. バイナリーに実行権限を付与:

        ```sh
        chmod +x /usr/local/bin/awsexec
        ```

### バイナリーでのインストール方法 (ビルド)

* Docker がインストールされている環境であれば、バイナリーをビルドすることもできます。ビルドしたバイナリーは `build/awsexec` に生成されます。
* 例えば以下の手順でインストールできます:

    1. バイナリーをビルドする:

        ```sh
        GOOS=darwin GOARCH=arm64 docker compose run --rm build
        ```

    2. バイナリーを配置する:

        ```sh
        sudo cp build/awsexec /usr/local/bin/awsexec
        ```

### シェルスクリプトのインストール方法

* シェルスクリプトが利用可能な環境であれば、シェルスクリプトをコピペでご自分の環境に配置することもできます。
    * Linux や macOS では通常利用可能です。
    * Windows では WSL 環境でのみ利用可能です。
* 例えば以下の手順でインストールできます:

    1. シェルスクリプトをダウンロード:

        ```sh
        curl -L -o /tmp/awsexec.sh https://raw.githubusercontent.com/ikedam/awsexec/refs/heads/main/awsexec.sh
        ```

    2. シェルスクリプトを配置:

        ```sh
        sudo cp /tmp/awsexec.sh /usr/local/bin/awsexec
        ```

    3. シェルスクリプトを実行権限を付与:

        ```sh
        chmod +x /usr/local/bin/awsexec
        ```


## awsexec の使用方法

* プロファイルを指定する場合:

    ```sh
    awsexec profile -- command
    ```

    * 指定の AWS プロファイルの認証情報を環境変数に設定した状態で、指定のコマンドを実行します。

* プロファイルの指定を `AWS_PROFILE` 環境変数で行う場合:

    ```sh
    AWS_PROFILE=profile awsexec -- command
    ```

    * 誤ったプロファイルの利用を避けるため、 `AWS_PROFILE` 環境変数の指定を必須にしています。
    * 複数のプロファイルを使い分ける必要がない、という場合は `AWS_PROFILE=default` をユーザープロファイルで設定しておくとよいでしょう。

## バイナリーのビルド方法

Docker を使用してバイナリーのビルドができます。ビルドしたバイナリーは `build/awsexec` に生成されます。

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

## ライセンス

本アプリケーションおよびソースコードは [MIT ライセンス](LICENSE) で配布します。
