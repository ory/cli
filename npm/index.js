var binwrap = require('binwrap')
var path = require('path')

var packageInfo = require(path.join(__dirname, '..', 'package.json'))
var version = packageInfo.version
var root = 'https://github.com/ory/cli/releases/download/v' + version

module.exports = binwrap({
  dirname: __dirname,
  binaries: ['ory'],
  urls: {
    'linux-x64': root + '/ory_' + version + '-linux_64bit.tar.gz',
    'win32-x64': root + '/ory_' + version + '-windows_64bit.zip',
    'darwin-x64': root + '/ory_' + version + '-macOS_64bit.tar.gz',
    'darwin-arm64': root + '/ory_' + version + '-macOS_arm64.tar.gz'
  }
})
