var binwrap = require("binwrap");
var path = require("path");

var packageInfo = require(path.join(__dirname, '..', "package.json"));
var version = packageInfo.version;
var root = "https://github.com/ory/cli/releases/download/v" + version;

module.exports = binwrap({
  dirname: __dirname,
  binaries: [
    "ory"
  ],
  urls: {
    "linux-x64": root + "/ory_"+version+"_linux_64-bit.tar.gz",
    "win32-x64": root + "/ory_"+version+"_windows_64-bit.zip",
    "darwin-x64": root + "/ory_"+version+"_macOS_64-bit.tar.gz",
  }
});
