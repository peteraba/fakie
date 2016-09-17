[![Release](https://img.shields.io/github/release/peteraba/fakie.svg)](https://github.com/peteraba/fakie/releases/latest)
[![Build Status](https://travis-ci.org/peteraba/fakie.svg?branch=master)](https://travis-ci.org/peteraba/fakie)
[![Coverage Status](https://coveralls.io/repos/github/peteraba/fakie/badge.svg?branch=master)](https://coveralls.io/github/peteraba/fakie?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/peteraba/fakie)](https://goreportcard.com/report/github.com/peteraba/fakie)
[![GoDoc](https://godoc.org/github.com/peteraba/fakie?status.svg)](https://godoc.org/github.com/peteraba/fakie)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/peteraba/fakie/master/LICENSE.md)

# fakie

Tiny command line program that accepts a template and outputs fake data based on [phony](https://github.com/yields/phony).

## Examples

```bash
# publish email to nsq every 1ms.
echo '{"email":"{{email}}", "subject": "welcome!"}' \
  | fakie --tick 1ms \
  | json-to-nsq --topic users

# add users to FoundationDB.
echo "'set {{username}} {{avatar}}'" \
  | fakie \
  | xargs -L1 -n3 fdbcli --exec

# add users to MongoDB.
echo "'db.users.insert({ name: \"{{name}}\" })'" \
  | fakie \
  | xargs -L1 -n1 mongo --eval

# add users to Redis.
echo "set {{username}} {{avatar}}" \
  | fakie \
  | xargs -L1 -n3 redis-cli

# send a single request using curl.
echo 'country={{country}}' \
  | fakie --max 1 \
  | curl -d @- httpbin.org/post
```

## Installation

```bash
$ go get github.com/peteraba/fakie
```

## Usage

```text

Usage: fakie
  [--tick d]
  [--max n]
  [--batch n]
  [--list]
  [--concurrent]
  [--prefix s]
  [--postfix s]

  fakie -h | --help
  fakie -v | --version

Options:
  --tick d        generate data every d [default: 10ms]
  --max n         generate data up to n [default: -1]
  --batch n       batch size for concurrent runs [default: 100]
  --list d        list all available generators
  --concurrent    skip ticks and generate fake data concurrently
  --prefix s      text to display before the generated data [default: ]
  --postfix s     text to display after the generated data [default: ]
  -v, --version   show version information
  -h, --help      show help information

```

## Generators

```text
  avatar
  color
  country
  country.code
  domain
  domain.name
  domain.tld
  double
  email
  event.action
  http.method
  id
  ipv4
  ipv6
  latitude
  longitude
  mac.address
  name
  name.first
  name.last
  product.category
  product.name
  smartdouble:desiredStdDev,desiredMean,min,max
  smartunixtime:deviationDays,rangeDays
  smartdate:format,deviationDays,rangeDays
  state
  state.code
  timezone
  unixtime
  username
```

### Generator attributes

Basic generators coming from phony have no attributes, but the more advanced fakie ones do.

#### Smart double

All arguments are optional.

- Desired standard deviation
- Desired mean
- Desired minimum value
- Desired maximum value

More info: https://golang.org/pkg/math/rand/#Rand.NormFloat64

*Example:*

The following will output a random number with the standard deviation of 10.2 and mean of -8.7, but no smaller than -80.3 and no larger than 90.2.

```
echo '{{smartdouble:10.2,-8.7,-80.3,90.2}}' | fakie --max 1
```

#### Smart unixtime

Smart unixtime will - as the name suggests - give you a unix timestamp. It will start with `now`.

- Deviation in days (float): Value will be added to the unixtime of `now`.
- Range in days (float): Value will be used to randomize the unix timestamp.

*Example:*

If the current time is `2016-09-17 16:05:02`, than the following will get a unix timestamp representing a date between `2015-09-12 12:05:02` and `2015-09-22 22:05:02`

```
echo '{{smartunixtime:-366.0,5.25}}' | fakie --max 1
```

#### Smart date

- Date format. Allowed values:
  - ANSIC       ("Mon Jan _2 15:04:05 2006")
  - UnixDate    ("Mon Jan _2 15:04:05 MST 2006")
  - RubyDate    ("Mon Jan 02 15:04:05 -0700 2006")
  - RFC822      ("02 Jan 06 15:04 MST")
  - RFC822Z     ("02 Jan 06 15:04 -0700" // RFC822 with numeric zone)
  - RFC850      ("Monday, 02-Jan-06 15:04:05 MST")
  - RFC1123     ("Mon, 02 Jan 2006 15:04:05 MST")
  - RFC1123Z    ("Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone)
  - RFC3339     ("2006-01-02T15:04:05Z07:00")
  - RFC3339Nano ("2006-01-02T15:04:05.999999999Z07:00")
  - Kitchen     ("3:04PM")
  - Stamp       ("Jan _2 15:04:05")
  - StampMilli  ("Jan _2 15:04:05.000")
  - StampMicro  ("Jan _2 15:04:05.000000")
  - StampNano   ("Jan _2 15:04:05.000000000")
  - SqlDateTime ("2006-01-02 15:04:05")
  - SqlDate     ("2006-01-02")
  - SqlTime     ("15:04:05")
- Deviation in days
- Range in days

*Example:*

If the current time is `2016-09-17 16:05:02`, than the following will get a date between `2015-09-12 12:05:02` and 
`2015-09-22 22:05:02`, formatting will be compatible with MySQL.

```
echo '{{smartdate:SqlDateTime,-366.0,5.25}}' | fakie --max 1
```


## Numbering

You can define numbers in your template to reference previously defined templates:

```
echo '{"email":"{{email}}", "email_repeated": "{{0}}"}' \
  | fakie --max 1
```

This will output something similar to the following:
```
{"email":"tomaslau16@example.us","email_repeated":"tomaslau16@example.us"}
```

## Development

Pull requests are welcome, but if you intend to develop something for `fakie` than please make sure to install the
git hook shipped with the project:
```bash
./git/install_hooks.bash
```
