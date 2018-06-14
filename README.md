<p align="center">
    <img alt="Balrog logo, by Malenea" src="documentation/logo.png"></img>
</p>
<h4 align="center">An Alpine Linux package repository proxy with a simple WebUI</h3>
<p align="center">
    <a href="https://goreportcard.com/report/github.com/blippar/balrog">
        <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/blippar/balrog">
    </a>
    <a href="https://godoc.org/github.com/blippar/balrog">
        <img alt="GoDoc" src="https://godoc.org/github.com/blippar/balrog?status.svg">
    </a>
    <a href="https://github.com/blippar/balrog/blob/master/LICENSE">
       <img alt="License" src="  https://img.shields.io/github/license/blippar/balrog.svg">
    </a>
    <a href="https://github.com/blippar/balrog/releases/latest">
        <img alt="Latest Release" src="https://img.shields.io/github/release/blippar/balrog.svg">
    </a>
    <a href="https://hub.docker.com/r/blippar/balrog/">
        <img alt="Docker Image" src="https://img.shields.io/docker/automated/blippar/balrog.svg">
    </a>
</p>

---

Balrog is a lightweight **Alpine Linux package repository UI and Proxy** made to ease the distribution of privately built packages in a company network. Its goal is to be **simple to use and deploy** while providing enough information about what is currently served.

## Features

- **Host your own APK Repository** from multiple backends (filesystem or AWS S3)
- **Compatible with multiple repository branch** to support multiple Alpine Linux version or distribute testing packages
- **List and search** packages and versions in a repository branch with the **built-in UI**
- **Serve extra static files** if needed (useful to provide standard assets when building a Docker container for example)

## Screenshots

<p align="center">
    <a href="documentation/screenshots/repository-list.png">
        <img width="320" alt="Repository Listing" src="documentation/screenshots/repository-list.png"></img>
    </a>
    <a href="documentation/screenshots/package-search.png">
        <img width="320" alt="Package List & Search" src="documentation/screenshots/package-search.png"></img>
    </a>
</p>

## Install

Balrog is currently only packaged as a Docker container but can be built and deployed manually on most system using Go 1.8+ and NPM.

You can get the packaged container on the [DockerHub](https://hub.docker.com/r/blippar/balrog/) and fetch it locally using:

```sh
docker pull blippar/balrog
```

## Getting started with Docker

The easiest way to run Balrog is using Docker, this is done by just creating a configuration file and running the container with the right options.

In this example we will be using [`config.docker.json`](config.docker.json) as our config file, this files contains everything needed to host an APK Repository from your filesystem but assumes a number of things:
- that your repository folder (the generated `packages` folder from `abuild`) is available locally at `/opt/apk/packages` with a `main` repository branch created
- that your public signing key (used to sign your packages and `APKINDEX`) is available at `/opt/apk/.abuild/builder.rsa.pub`

In order to start host your Repository UI at `localhost:8000` run the following commands:
```sh
docker pull blippar/balrog
docker run --restart=always \
           --volume "./config.docker.json:/app/config.json:ro" \
           --volume "/opt/apk/packages:/data/packages:ro" \
           --volume "/opt/apk/.abuild/builder.rsa.pub:/data/packages/builder.rsa.pub" \
           -p 8000:8000 blippar/balrog
```

## Configuration
#### Usage
```
Usage: balrog [--config CONFIG]

Options:
  --config CONFIG, -c CONFIG
                         path to your config file [env: BALROG_CONFIG] [default: config.json]
  --help, -h             display this help and exit
  --version              display version and exit
```

#### Configuration file

Balrog can only be configured using a JSON configuration file (`config.json` in the current working directory by default), this configuration file looks like the following document:

```json
{
  "title": "My APK Repository",
  "http": {
    "addr": "0.0.0.0:8000",
    "templates": "templates/*.html.tmpl",
    "dist": "dist/"
  },
  "storage": {
    "type": "filesystem",
    "filesystem": {
      "prefix": "/data/packages",
      "cache":  60
    },
    "minio": {
      "endpoint":    "s3.amazonaws.com",
      "access_key":  "AKIAIOSFODNN7EXAMPLE",
      "private_key": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
      "region":      "us-west-2",
      "ssl":         true,
      "bucket":      "my-aports",
      "prefix":      "path/to/my/root"
    }
  },
  "repositories": [
    {
      "name": "main",
      "available_arch": ["x86","x86_64","armhf", "aarch64","ppc64le","s390x"],
      "storage_prefix": "main",
      "public_keys":[
        "me@myorg.com-deadbeef.rsa.pub"
      ]
    }
  ],
  "static": [
    {
      "web_prefix": "assets",
      "storage_prefix": "my_assets"
    }
  ]
}
```

Each section of this document configures a different part of balrog, following is a quick run-down of what each section does and what is expected for each of those:

##### Root

| Name         | Type       | Default                 | Description                       |
|:-------------|:-----------|:------------------------|:----------------------------------|
| `title`      | `string`   | `APK Browser`           | Title of the page in the WebUI    |
| `http`       | `object`   | See `http` config       | HTTP server configuration         |
| `storage`    | `object`   | See `storage` config    | Backend storage configuration     |
| `repository` | `[]object` | See `repository` config | List of repository to serve       |
| `static`     | `[]object` | See `static` config     | List of static folder to serve    |

##### `http`
| Name             | Type     | Default                 | Description                                       |
|:-----------------|:---------|:------------------------|:--------------------------------------------------|
| `http.addr`      | `string` | `0.0.0.0:8000`          | HTTP router listen address                        |
| `http.templates` | `string` | `templates/`            | Path to Balrog's UI template folder               |
| `http.dist`      | `string` | `dist/`                 | Path to Balrog's UI assets folder                 |

##### `storage`
| Name                 | Type     | Default  | Description                                                                    |
|:---------------------|:---------|:---------|:-------------------------------------------------------------------------------|
| `storage.type`       | `string` | n/a      | Required, defines what storage backend to use (can be `minio` or `filesystem`) |
| `storage.filesystem` | `object` | n/a      | Filesystem storage backend configuration                                       |
| `storage.minio`      | `object` | n/a      | Minio storage backend configuration (AWS S3 compatible)                        |

##### `storage.filesystem`
| Name                        | Type     | Default | Description                                                            |
|:----------------------------|:---------|:--------|:-----------------------------------------------------------------------|
| `storage.filesystem.prefix` | `string` | n/a     | Required, a path prefix that will be used as your storage backend root |
| `storage.filesystem.cache`  | `int`    | n/a     | Optional, if > 0 activate in memory cache of files for X seconds       |

##### `storage.minio`
| Name                        | Type     | Default            | Description                                                           |
|:----------------------------|:---------|:-------------------|:----------------------------------------------------------------------|
| `storage.minio.endpoint`    | `string` | `s3.amazonaws.com` | Minio backend endpoint (AWS S3 by default)                            |
| `storage.minio.region`      | `string` | n/a                | Bucket region (AWS Region for S3)                                     |
| `storage.minio.access_key`  | `string` | n/a                | Required, bucket's access key (AWS AccessKeyID for S3)                |
| `storage.minio.private_key` | `string` | n/a                | Required, bucket's private key (AWS SecretAccessKey for S3)           |
| `storage.minio.ssl`         | `bool`   | `false`            | Activate SSL transport                                                |
| `storage.minio.bucket`      | `string` | n/a                | Required, name of the bucket where your AlpineLinux package folder is |
| `storage.minio.prefix`      | `string` | `/`                | A path prefix that will be used as your storage backend root          |

##### `repository`
| Name                        | Type       | Default | Description                                                                                                                                          |
|:----------------------------|:-----------|:--------|:-----------------------------------------------------------------------------------------------------------------------------------------------------|
| `repository.name`           | `string`   | n/a     | Name of the repository (will be used as an URL path prefix, use an URL safe string here)                                                             |
| `repository.available_arch` | `[]string` | n/a     | List of available arch for your architecture (see the content of `packages/reponame`, this should list all architecture supported by the repository) |
| `repository.storage_prefix` | `string`   | n/a     | Path to your repository package folder (the `packages/reponame` folder created by `abuild`) relative to your storage backend root                    |
| `repository.public_keys`    | `[]string` | n/a     | List of path to your signing public keys relative to your storage backend root (this allows you to serve your public keys easily)                    |

##### `static`
| Name                    | Type     | Default | Description                                                                                                                                                                 |
|:------------------------|:---------|:--------|:----------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `static.web_prefix`     | `string` | n/a     | URL path under which your static folder will be served by Balrog web server (do not use `/dist/` or the name of any of your repository as this can cause unwanted behavior) |
| `static.storage_prefix` | `string` | n/a     | Path to your static folder relative to your storage root                                                                                                                    |

## Development

In order to build `balrog` you will first need to make sure that you have Go 1.8+ installed, we do commit all assets and if you are not working with it you should not need anything else.

From there all elements of Balrog can be built via a Makefile at the root of this repository, available rules are the followings:

Using this `Makefile` you will have access to the following rules:
- `make balrog`: build a binary for your current system under `bin/balrog`
- `make static`: build a static binary targeted at Linux based system under `bin/balrog`
- `make assets`: minify all assets needed for Balrog's UI (require NodeJS 6+)
- `make all`: build the binary for your current system under `bin/balrog` then minify all assets (require NodeJS 6+)
- `make docker`: build the `blippar/balrog` Docker container (tagged with the current version)
- `make test`: test balrog source code using `go test -cover`
- `make lint`: lint balrog source code using `gometalinter`
- `make clean`: remove all files that would be created by `make balrog`

If you want to make any change to the Javascript or CSS assets used in Balrog (other than the ones included in the `templates/` folder) you will need to have NodeJS 6+ installed and run `npm install` before being able to rebuild the asset collection.

If any change is done, or extra assets need to be included you might need to update `Gulpfile.js` and run either `make assets` (which will use `gulp` from `node_modules`) or use `gulp` directly, here is the list of exposed rules in our `gulpfile.js`:
- `gulp build`: cleanup previous build, then concat and minify JS and CSS assets as well as move the necessary fonts to the `dist` folder
- `gulp scripts`: concat and minify all JS assets to the `dist` folder
- `gulp styles`: concat and minify all CSS assets to the `dist` folder
- `gulp fonts`: copy the necessary fonts to the `dist` folder
- `guild clean`: cleanup the `dist` folder`

Using those 2 tools you should be able to build and extend Balrog functionality and UI as you wish, just remember to create your own configuration file pointing to a local AlpineLinux repository to be able to test it properly.

## Credits

Thanks to [Malenea](https://github.com/malenea) for his awesome work on the logo.

Balrog's logo is licensed under the [Creative Commons Attribution 4.0 International License](http://creativecommons.org/licenses/by/4.0/) and was widely inspired on the original Go gopher designed by [Renee French](http://reneefrench.blogspot.com/).
