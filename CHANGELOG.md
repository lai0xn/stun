# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Comprehensive logging system with configurable levels and formats
- Structured logging with contextual fields for STUN operations
- Specialized logging methods for client and server operations
- JSON and text logging formats
- Caller information support for debugging
- Production-ready logging configuration options
- Detailed documentation for all public APIs
- Comprehensive README with usage examples
- Logging system documentation (LOGGING.md)
- API reference documentation
- Protocol details and error handling documentation
- String method for MessageType for better logging
- Enhanced error handling with contextual information

### Changed
- Improved server logging with detailed request/response tracking
- Enhanced client logging with connection and transaction details
- Updated examples to use improved logging system
- Better error messages with contextual information
- More descriptive log messages with structured fields

### Fixed
- Logger type issues in server configuration
- Missing documentation for public APIs
- Inconsistent error handling patterns
- Log message formatting and structure

## [0.1.0] - 2023-12-01

### Added
- Initial STUN protocol implementation
- Basic client and server functionality
- XOR-MAPPED-ADDRESS attribute support
- Basic error handling
- Simple logging with logrus
- Example client and server applications

### Features
- STUN binding request/response support
- IPv4 address handling
- Transaction ID generation and validation
- Magic cookie validation
- Basic UDP communication
- Simple message parsing and encoding

## Types of Changes

- **Added**: New features or functionality
- **Changed**: Changes to existing functionality
- **Deprecated**: Features that will be removed in future releases
- **Removed**: Features that have been removed
- **Fixed**: Bug fixes
- **Security**: Security-related fixes 