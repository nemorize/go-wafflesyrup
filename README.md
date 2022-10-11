<h1 style="text-align: center;">WaffleSyrup</h1>
<div style="text-align: center;">
    Simple backup solution written by Go.<br /><br />
    <img src="https://img.shields.io/badge/License-MIT-blue" alt="License MIT" />
    <img src="https://img.shields.io/badge/golang-1.17-01A7D0?logo=go" alt="Golang 1.17" />
    <img src="https://img.shields.io/github/workflow/status/qroffle/wafflesyrup/Build%20&%20Release" alt="Build status, sorry for screen readers." />
</div>

## Usage

WaffleSyrup runs in the current working directory.
It will create `./tmp` directory to save tarballs which includes archived files.
So you need to obtain proper permissions to cwd.

WaffleSyrup is a command line tool. Simply type `wafflesyrup` in the console for help.
The basic commands are:

* `wafflesyrup start`: Start backup.
* `wafflesyrup version`: Displays the current running version of WaffleSyrup. Aliased as v.

### Build from source
* Install go(1.17) from https://go.dev
* Download source or `git clone https://github.com/qroffle/wafflesyrup`
* Change GOOS or GOARCH if you want to change OS or ARCH.
* Execute `./build.sh` on UNIX based OS, `./build.bat` on windows.
* Check the `./bin` directory for generated binary.

## Configuration

First, you need to create `config.json` which contains configuration data.
You can simply start with renaming `config.json.example` to `config.json`.

### Backups
`config.json` has one `backups` array, which contains some backup sets.
Each backup item must contain `name`, `directories`, `savepoints` and can contain `postscripts`.

```json
{
  "backups": [
    {
      "name": "example",
      "directories": [],
      "savepoints": [],
      "postscripts": []
    }
  ]
}
```

### Directories
`directories` has array-typed value, that can contain some backup directory information.
Each directory must contain `path` and can contain `excluded` array.

```json
{
  "directories": [
    {
      "path": "/home/example/cute_waffle",
      "excluded": [
        "node_modules/",
        ".git"
      ]
    }
  ]
}
```

### Savepoints
`savepoints` has array-typed value, that can contain some backup destination information.
Each savepoint must contain `driver`, `identity`. Each driver should get another structures of identity.
Show details at [Supported Drivers](#supported-drivers)

```json
{
  "savepoints": [
    {
      "driver": "sftp",
      "identity": {}
    }
  ]
}
```

### Postscripts
Currently unavailable.

## Supported Drivers

WaffleSyrup supports 3 drivers.
* [sftp](#sftp)
* [gdrive (google-drive)](#gdrive)
* [dropbox](#dropbox)

### sftp
`sftp` requires `host`, `port`, `username`, `password` fields at `identity`.

```json
{
  "driver": "sftp",
  "identity": {
    "host": "127.0.0.1",
    "port": "22",
    "username": "bitchchecker",
    "password": "has_joined_#stopHipHop"
  }
}
```
~~LMAO https://ubuntuforums.org/archive/index.php/t-825400.html~~

### gdrive
`gdrive` requires `credentialFile`, `folderId` fields at `identity`.
Credential files can be obtained from https://console.cloud.google.com.

Gdrive requires oauth tokens.

WaffleSyrup will give you some authentication links, and ask you the token.
After receives, WaffleSyrup create ./tokens/ directory and will save tokens into.

```json
{
  "driver": "gdrive",
  "identity": {
    "credentialFile": "./credentials/client_secret_xxx.apps.googleusercontent.com.json",
    "folderId": ""
  }
}
```

### dropbox
`dropbox` requires `accessToken`, `folderPath` fields at `identity`.
You can create an app and obtain accessToken from https://www.dropbox.com/developers/apps/create.

```json
{
  "driver": "dropbox",
  "identity": {
    "accessToken": "sl.xxx",
    "folderPath": "/path/to/upload"
  }
}
```
