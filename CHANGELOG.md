# Change Log
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased]
### Removed
- Remove deprecated -url flag.

## [0.1.1] - 2016-11-17
### Added
- The `dir` command to expose a local directory via publichost.me.

### Changed
- Introduces the `http` command to start an http tunnel. This means the client does not accept an url as the first argument. Change `publichost http://localhost:3000` to `publichost http http://localhost:3000`

## [0.1] - 2016-11-11
### Added
- Initial release for client and server
