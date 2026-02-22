bcrypt-tool
===========

`bcrypt-tool` is a dandy CLI tool for generating and matching bcrypt hashes

It is forked from [shoening/bcrypt-tool](https://github.com/shoenig/bcrypt-tool), but deviates in CLI.

### Installation

#### Install from Releases

The `bcrypt-tool` tool is available from the [Releases](https://github.com/Necoro/bcrypt-tool/releases) page.

It is pre-compiled for many operating systems and architectures including

- Linux
- Windows
- MacOS

#### Build from source
The `bcrypt-tool` command can be compiled by running
```bash
$ go install github.com/Necoro/bcrypt-tool@latest
```

### Usage
```bash
bcrypt-tool [action] [flags] parameter ...
```

#### Password Provision

The password can be provided by different means:

* as a direct argument. Note that the password may be (shortly) visible in process browsers (like `ps aux` or `top`).
* on stdin, i.e., by piping the password (`gen-pwd | bcrypt-tool hash` or `bcrypt-tool hash < password_file`)
* interactively. `bcrypt-tool` queries for the password if it has not been provided by the previous means 
and is connected to a terminal

#### Actions

`hash`
: Use bcrypt to generate a hash from the provided password and optional cost (4–31).

`match`
: Use bcrypt to check if a password matches a hash.

`cost`
: Use bcrypt to determine the cost of a hash (4–31).


### Examples

#### Generate hash from a directly provided password
```bash
$ bcrypt-tool hash p4ssw0rd
```


#### Generate hash as output from `pass`
```bash
$ pass some/pwd | bcrypt-tool hash
```

#### Generate hash from an interactively provided password with cost
```bash
$ bcrypt-tool hash --cost 31
Password:
```
    
#### Determine if password from file matches hash
```bash
$ bcrypt-tool match '$2a$10$nWFwjoFo4zhyVosdYMb6XOxZqlVB9Bk0TzOvmuo16oIwMZJXkpanW' < password_file
```

#### Determine cost of hash
```bash
$ bcrypt-tool cost '$2a$10$nWFwjoFo4zhyVosdYMb6XOxZqlVB9Bk0TzOvmuo16oIwMZJXkpanW'
```

# License

This module is open source under the [MIT](LICENSE) license.
