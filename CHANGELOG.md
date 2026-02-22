# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## 2.0
### Major Changes
- Password can now be supplied via stdin or interactively. 
`bcrypt-tool` no longer relies on having the password being passed visibly on the commandline.
- CLI ported to [Kong](https://github.com/alecthomas/kong).

**Both points result in slight changes to the CLI structure and command names!**

### Added
- `match` gained a `--quiet`/`-q` option to suppress output
- `hash` and `match` gained a `--truncate` option to take as input the first 72 bytes of stdin.
This can be used, for example, to provide a binary file (e.g., a picture) as the password.
- Improved validation of input arguments

## 1.x.x
Please consult the releases in the original repository: https://github.com/shoenig/bcrypt-tool/releases