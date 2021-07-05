# kazaam

[![Travis Build Status](https://img.shields.io/travis/qntfy/kazaam.svg?branch=master)](https://travis-ci.com/qntfy/kazaam)
[![Coverage Status](https://coveralls.io/repos/github/qntfy/kazaam/badge.svg?branch=master)](https://coveralls.io/github/qntfy/kazaam?branch=master)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENSE)
[![GitHub release](https://img.shields.io/github/release/qntfy/kazaam.svg?maxAge=3600)](https://github.com/qntfy/kazaam/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/qntfy/kazaam)](https://goreportcard.com/report/github.com/qntfy/kazaam)
[![GoDoc](https://godoc.org/github.com/qntfy/kazaam?status.svg)](https://pkg.go.dev/github.com/qntfy/kazaam/v4)

## Description

Kazaam was created with the goal of supporting easy and fast transformations of JSON data with Golang.
This functionality provides us with an easy mechanism for taking intermediate JSON message representations
and transforming them to formats required by arbitrary third-party APIs.

Inspired by [Jolt](https://github.com/bazaarvoice/jolt), Kazaam supports JSON to JSON transformation via a
transform "specification" also defined in JSON. A specification is comprised of one or more "operations". See
Specification Support, below, for more details.

## Documentation

API Documentation is available on [pkg.go.dev](https://pkg.go.dev/github.com/qntfy/kazaam/v4).

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
- timestamp
- uuid
- default
- pass
- delete

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
- *Array append/prepend and set*: Append and prepend an array with `[+]` and `[-]`. Attempting to write an array element that does not exist results in null padding as needed to add that element at the specified index (useful with `"inplace"`).

The shift transform also supports a `"require"` field. When set to `true`,
Kazaam will throw an error if *any* of the paths in the source JSON are not
present.

Finally, shift by default is destructive. For in-place operation, an optional `"inplace"`
field may be set.

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

Coalesce also supports an `ignore` array in the spec. If an otherwise matching key has a value in `ignore`, it is not considered a match.
This is useful e.g. for empty strings

```javascript
{
  "operation": "coalesce",
  "spec": {
    "ignore": [""],
    "firstObjectId": ["doc.guidObjects[0].uid", "doc.guidObjects[0].id"]
  }
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

### Timestamp

A `timestamp` transform parses and formats time strings using the golang
syntax. **Note**: this operation is done in-place. If you want to preserve the
original string(s), pair the transform with `shift`. This transform also
supports the `$now` operator for `inputFormat`, which will set the current
timestamp at the specified path, formatted according to the `outputFormat`.
`$unix` is supported for both input and output formats as a Unix time, the
number of seconds elapsed since January 1, 1970 UTC as an integer string.

```javascript
{
  "operation": "timestamp",
  "timestamp[0]": {
    "inputFormat": "Mon Jan _2 15:04:05 -0700 2006",
    "outputFormat": "2006-01-02T15:04:05-0700"
  },
  "nowTimestamp": {
    "inputFormat": "$now",
    "outputFormat": "2006-01-02T15:04:05-0700"
  },
  "epochTimestamp": {
    "inputFormat": "2006-01-02T15:04:05-0700",
    "outputFormat": "$unix"
  }
}

```

executed on a json message with format

```javascript
{
  "timestamp": [
    "Sat Jul 22 08:15:27 +0000 2017",
    "Sun Jul 23 08:15:27 +0000 2017",
    "Mon Jul 24 08:15:27 +0000 2017"
  ]
}
```

would result in

```javascript
{
  "timestamp": [
    "2017-07-22T08:15:27+0000",
    "Sun Jul 23 08:15:27 +0000 2017",
    "Mon Jul 24 08:15:27 +0000 2017"
  ]
  "nowTimestamp": "2017-09-08T19:15:27+0000"
}
```

### UUID

A `uuid` transform generates a UUID based on the spec. Currently supports UUIDv3, UUIDv4, UUIDv5.

For version 4 is a very simple spec

```javascript
{
    "operation": "uuid",
    "spec": {
        "doc.uuid": {
            "version": 4, //required
        }
    }
}
```

executed on a json message with format

```javascript
{
  "doc": {
    "author_id": 11122112,
    "document_id": 223323,
    "meta": {
      "id": 23
    }
  }
}
```

would result in

```javascript
{
  "doc": {
    "author_id": 11122112,
    "document_id": 223323,
    "meta": {
      "id": 23
    }
    "uuid": "f03bacc1-f4e0-4371-a5c5-e8160d3d6c0c"
  }
}
```

For UUIDv3 & UUIDV5 are a bit more complex. These require a Name Space which is a valid UUID already, and a set of paths, which generate UUID's based on the value of that path. If that path doesn't exist in the incoming document, a default field will be used instead. **Note** both of these fields must be strings.
**Additionally** you can use the 4 predefined namespaces such as `DNS`, `URL`, `OID`, & `X500` in the name space field otherwise pass your own UUID.

```javascript
{
   "operation":"uuid",
   "spec":{
      "doc.uuid":{
         "version":5,
         "namespace":"DNS",
         "names":[
            {"path":"doc.author_name", "default":"some string"},
            {"path":"doc.type", "default":"another string"},
         ]
      }
   }
}
```

executed on a json message with format

```javascript
{
  "doc": {
    "author_name": "jason",
    "type": "secret-document"
    "document_id": 223323,
    "meta": {
      "id": 23
    }
  }
}
```

would result in

```javascript
{
  "doc": {
    "author_name": "jason",
    "type": "secret-document",
    "document_id": 223323,
    "meta": {
      "id": 23
    },
    "uuid": "f03bacc1-f4e0-4371-a7c5-e8160d3d6c0c"
  }
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

### Delete

A delete transform provides the ability to delete keys in place.

```javascript
{
  "operation": "delete",
  "spec": {
    "paths": ["doc.uid", "doc.guidObjects[1]"]
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
    "guid": ["guid0", "guid2", "guid4"],
    "guidObjects": [{"id": "guid0"}, {"id": "guid4"}]
  }
}
```

### Pass

A pass transform, as the name implies, passes the input data unchanged to the output. This is used internally
when a null transform spec is specified, but may also be useful for testing.

## Usage

To start, go get the versioned repository:

```sh
go get gopkg.in/qntfy/kazaam.v3
```

### Using as an executable program

If you want to create an executable binary from this project, follow
these steps (you'll need `go` installed and `$GOPATH` set):

``` shell
go get gopkg.in/qntfy/kazaam.v3
cd $GOPATH/src/gopkg.in/qntfy/kazaam.v3/kazaam
go install
```

This will create an executable in `$GOPATH/bin` like you
would expect from the normal `go` build behavior.

### Examples

See [godoc examples](https://godoc.org/pkg/gopkg.in/qntfy/kazaam.v3/#pkg-examples).
