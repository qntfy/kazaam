# kazaam

[![Travis Build Status](https://img.shields.io/travis/qntfy/kazaam.svg?branch=master)](https://travis-ci.org/qntfy/kazaam)
[![Coverage Status](https://coveralls.io/repos/github/qntfy/kazaam/badge.svg?branch=master)](https://coveralls.io/github/qntfy/kazaam?branch=master)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENSE)
[![GitHub release](https://img.shields.io/github/release/qntfy/kazaam.svg?maxAge=3600)](https://github.com/qntfy/kazaam/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/qntfy/kazaam)](https://goreportcard.com/report/github.com/qntfy/kazaam) 
[![GoDoc](https://godoc.org/github.com/qntfy/kazaam?status.svg)](http://godoc.org/gopkg.in/qntfy/kazaam.v2)

## Description
Kazaam was created with the goal of supporting easy and fast transformations of JSON data with Golang.
This functionality provides us with an easy mechanism for taking intermediate JSON message representations
and transforming them to formats required by arbitrary third-party APIs.

Inspired by [Jolt](https://github.com/bazaarvoice/jolt), Kazaam supports JSON to JSON transformation via a
transform "specification" also defined in JSON. A specification is comprised of one or more "operations". See
Specification Support, below, for more details.

## Documentation
API Documentation is available at http://godoc.org/gopkg.in/qntfy/kazaam.v2.

## Features
Kazaam is primarily designed to be used as a library for transforming arbitrary JSON.
It ships with six built-in transform types, described below, which provide significant flexibility
in reshaping JSON data.

Also included when you `go get` Kazaam, is a binary implementation, `kazaam` that can be used for
development and testing of new transform specifications.

Finally, Kazaam supports the implementation of custom transform types. We encourage and appreciate
pull requests for new transform types so that they can be incorporated into the Kazaam distribution,
but understand sometimes time-constraints or licensing issues prevent this. See the API documentation
for details on how to write and register custom transforms.

Due to performance considerations, Kazaam does not fully validate that input data is valid JSON. The
`IsJson()` function is provided for convenience if this functionality is needed, it may significantly slow
down use of Kazaam.

## Specification Support
Kazaam currently supports the following transforms:
- shift
- concat
- coalesce
- extract
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
    "guid": ["guid0", "guid2", "guid4"],
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

The shift transform also supports a `"require"` field. When set to `true`,
Kazaam will throw an error if *any* of the paths in the source JSON are not
present.

### Concat
The concat transform allows the combination of fields and literal strings into a single string value.
```javascript
{
    "operation": "concat",
    "spec": {
        "sources": [{
            "value": "TEST"
        }, {
            "path": "a.timestamp"
        }],
        "targetPath": "a.timestamp",
        "delim": ","
    }
}
```

executed on a JSON message with format
```javascript
{
    "a": {
        "timestamp": 1481305274
    }
}
```

would result in
```javascript
{
    "a": {
        "timestamp": "TEST,1481305274"
    }
}
```

Notes:
- *sources*: list of items to combine (in the order listed)
  - literal values are specified via `value`
  - field values are specified via `path` (supports the same addressing as `shift`)
- *targetPath*: where to place the resulting string
  - if this an existing path, the result will replace current value.
- *delim*: Optional delimiter

The concat transform also supports a `"require"` field. When set to `true`,
Kazaam will throw an error if *any* of the paths in the source JSON are not
present.

### Coalesce
A coalesce transform provides the ability to check multiple possible keys to find a desired value. The first matching key found of those provided is returned.
```javascript
{
  "operation": "coalesce",
  "spec": {
    "firstObjectId": ["doc.guidObjects[0].uid", "doc.guidObjects[0].id"]
  }
}
```

executed on a json message with format
```javascript
{
  "doc": {
    "uid": 12345,
    "guid": ["guid0", "guid2", "guid4"],
    "guidObjects": [{"id": "guid0"}, {"id": "guid2"}, {"id": "guid4"}]
  }
}
```

would result in
```javascript
{
  "doc": {
    "uid": 12345,
    "guid": ["guid0", "guid2", "guid4"],
    "guidObjects": [{"id": "guid0"}, {"id": "guid2"}, {"id": "guid4"}]
  },
  "firstObjectId": "guid0"
}
```

### Extract
An `extract` transform provides the ability to select a sub-object and have kazaam return that sub-object as the top-level object. For example
```javascript
{
  "operation": "extract",
  "spec": {
    "path": "doc.guidObjects[0].path.to.subobject"
  }
}
```

executed on a json message with format
```javascript
{
  "doc": {
    "uid": 12345,
    "guid": ["guid0", "guid2", "guid4"],
    "guidObjects": [{"path": {"to": {"subobject": {"name": "the.subobject", "field", "field.in.subobject"}}}}, {"id": "guid2"}, {"id": "guid4"}]
  }
}
```

would result in
```javascript
{
  "name": "the.subobject",
  "field": "field.in.subobject"
}
```

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

To start, go get the versioned repository:
```sh
go get gopkg.in/qntfy/kazaam.v2
```

### Using as an executable program

If you want to create an executable binary from this project, follow
these steps (you'll need `go` installed and `$GOPATH` set):

``` shell
go get gopkg.in/qntfy/kazaam.v2
cd $GOPATH/src/gopkg.in/qntfy.kazaam.v2/kazaam
go install
```

This will create an executable in `$GOPATH/bin` like you
would expect from the normal `go` build behavior.

### Examples

See [godoc examples](https://godoc.org/pkg/gopkg.in/qntfy/kazaam.v2/#pkg-examples).
