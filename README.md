bcrypt-tool
===========

`bcrypt-tool` is a dandy CLI tool for generating and matching bcrypt hashes

It is forked from [shoening/bcrypt-tool](https://github.com/shoenig/bcrypt-tool).

### Installation

#### Install from Releases

The `bcrypt-tool` tool is available from the [Releases](https://github.com/Necoro/bcrypt-tool/releases) page.

It is pre-compiled for many operating systems and architectures including

- Linux
- Windows
- MacOS
- FreeBSD
- OpenBSD
- Plan9

#### Build from source
The `bcrypt-tool` command can be compiled by running
```bash
$ go install github.com/Necoro/bcrypt-tool@latest
```

### Usage
```bash
bcrypt-tool [action] parameter ...
```

#### Actions

- `hash  [password] <cost>` Use bcrypt to generate a hash given password and optional cost (4-31)

- `match [password] [hash]` Use bcrypt to check if a password matches a hash

- `cost  [hash]` Use bcrypt to determine the cost of a hash (4-31)


### Examples

#### Generate Hash from a Password
```bash
$ bcrypt-tool hash p4ssw0rd
```
    
#### Generate Hash from a Password with Cost
```bash
$ bcrypt-tool hash p4ssw0rd 31
```
    
#### Determine if Password matches Hash
```bash
$ bcrypt-tool match p4ssw0rd $2a$10$nWFwjoFo4zhyVosdYMb6XOxZqlVB9Bk0TzOvmuo16oIwMZJXkpanW
```

note: depending on your shell, you may need to escape the $ characters

#### Determine Cost of Hash
```bash
$ bcrypt-tool cost $2a$10$nWFwjoFo4zhyVosdYMb6XOxZqlVB9Bk0TzOvmuo16oIwMZJXkpanW
```
    
note: depending on your shell, you may need to escape the $ characters

# License

This module is open source under the [MIT](LICENSE) license.
