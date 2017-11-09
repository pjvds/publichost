# Change Log
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased]

### Added
- Support multiple hostnames. Only the subdomain is used for matching.
- TLS flag for the client.

### Changed
- Tunnel identifiers changed to xid's.
- Public addresses for tunnels are now names instead of numbers.
- Validate existance of local directory when serving it from the client.

### Removed
- Remove deprecated -url flag.

### Fixed
- Fixed public address not returned to client.

## [0.1.1]
### Added
- The `dir` command to expose a local directory via publichost.me.

### Changed
- Introduces the `http` command to start an http tunnel. This means the client does not accept an url as the first argument. Change `publichost http://localhost:3000` to `publichost http http://localhost:3000`

## [0.1]
### Added
- Initial release for client and server
