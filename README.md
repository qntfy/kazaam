# kazaam

[![Travis Build Status](https://img.shields.io/travis/qntfy/kazaam.svg?branch=master)](https://travis-ci.org/qntfy/kazaam)
[![Coverage Status](https://coveralls.io/repos/github/qntfy/kazaam/badge.svg?branch=master)](https://coveralls.io/github/qntfy/kazaam?branch=master)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENSE)
[![GitHub release](https://img.shields.io/github/release/qntfy/kazaam.svg?maxAge=2592000)](https://github.com/qntfy/kazaam/releases/latest)
[![GoDoc](https://godoc.org/github.com/qntfy/kazaam?status.svg)](http://godoc.org/gopkg.in/qntfy/kazaam.v1)

## Description
Kazaam was created with the goal of supporting easy and fast transformations of JSON data with Golang.
This functionality provides us with an easy mechanism for taking intermediate JSON message representations
and transforming them to formats required by arbitrary third-party APIs.

Inspired by [Jolt](https://github.com/bazaarvoice/jolt), Kazaam supports JSON to JSON transformation via a
transform "specification" also defined in JSON. A specification is comprised of one or more "operations". See
Specification Support, below, for more details.

## Documentation
API Documentation is available at http://godoc.org/gopkg.in/qntfy/kazaam.v1.

## Specification Support
Kazaam currently supports the following transforms:
- shift
- default
- pass

### Shift
The shift transform is the current Kazaam workhorse used for remapping of fields.
The specification supports jsonpath-esque JSON accesses and sets. Concretely
```javascript
{
  "operation": "shift",
  "spec": {
    "object.id": "doc.uid",
    "gid2": "doc.guid[1]",
    "allGuids": "doc.guidObjects[*].id"
  }
}
```

executed on a JSON message with format
```javascript
{
  "doc": {
    "uid": 12345,
    "guid": ["guid0", "guid2", "guid4"]
    "guidObjects": [{"id": "guid0"}, {"id": "guid2"}, {"id": "guid4"}]
  },
  "top-level-key": null
}
```

would result in
```javascript
{
  "object": {
    "id": 12345
  },
  "gid2": "guid2",
  "allGuids": ["guid0", "guid2", "guid4"]
}
```

The jsonpath implementation supports a few special cases:
- *Array accesses*: Retrieve `n`th element from array
- *Array wildcarding*: indexing an array with `[*]` will return every matching element in an array
- *Top-level object capture*: Mapping `$` into a field will nest the entire original object under the requested key

### Default
A default transform provides the ability to set a key's value explicitly. For example
```javascript
{
  "operation": "default",
  "spec": {
    "type": "message"
  }
}
```
would ensure that the output JSON message includes `{"type": "message"}`.

### Pass
A pass transform, as the name implies, passes the input data unchanged to the output. This is used internally
when a null transform spec is specified, but may also be useful for testing.

## Usage

To start, go get the versioned repository::
```sh
go get gopkg.in/qntfy/kazaam.v1
```

### Examples

See [godoc examples](https://godoc.org/pkg/gopkg.in/qntfy/kazaam.v1/#pkg-examples).
