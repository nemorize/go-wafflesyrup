<h1 style="text-align: center;">WaffleSyrup</h1>
<div style="text-align: center;">
    Simple backup solution written by Go.<br /><br />
    <img src="https://img.shields.io/badge/License-MIT-blue" />
    <img src="https://img.shields.io/badge/golang-1.17-01A7D0?logo=go" />
    <img src="https://img.shields.io/github/workflow/status/qroffle/wafflesyrup/Build%20&%20Release" />
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