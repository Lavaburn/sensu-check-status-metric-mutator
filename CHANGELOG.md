# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic
Versioning](http://semver.org/spec/v2.0.0.html).

## Unreleased

## [0.4.0] - 2023-02-09

Flags added to:
- Add occurrences as a separate metric instead of tag (avoids explosion of metrics in Prometheus)
- Make occurrences_watermark optional (in both scenarios)
- Change the tag names (useful when handler expects specific tags)

This release is backwards-compatible with nixwiz 0.3.0

## [0.3.0] - 2020-10-09

### Changed
- Fixed typo in README
- Updated dependencies to fix check runtime asset validation failure

## [0.2.1] - 2020-05-27

### Changed
- Add tags to README example output
- Formatting cleanup

## [0.2.0] - 2020-05-27

### Added
- Added tags to the check status metric point

## [0.1.1] - 2020-05-26

### Changed
- Changed default metric name template
- Changed check status timestamp

## [0.1.0] - 2020-05-25

### Added
- Initial release
