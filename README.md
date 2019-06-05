# kazaam

[![Travis Build Status](https://img.shields.io/travis/qntfy/kazaam.svg?branch=master)](https://travis-ci.org/qntfy/kazaam)
[![Coverage Status](https://coveralls.io/repos/github/qntfy/kazaam/badge.svg?branch=master)](https://coveralls.io/github/qntfy/kazaam?branch=master)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENSE)
[![GitHub release](https://img.shields.io/github/release/qntfy/kazaam.svg?maxAge=3600)](https://github.com/qntfy/kazaam/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/qntfy/kazaam)](https://goreportcard.com/report/github.com/qntfy/kazaam)
[![GoDoc](https://godoc.org/github.com/qntfy/kazaam?status.svg)](http://godoc.org/gopkg.in/qntfy/kazaam.v3)

## Description
Kazaam was created with the goal of supporting easy and fast transformations of JSON data with Golang.
This functionality provides us with an easy mechanism for taking intermediate JSON message representations
and transforming them to formats required by arbitrary third-party APIs.

Inspired by [Jolt](https://github.com/bazaarvoice/jolt), Kazaam supports JSON to JSON transformation via a
transform "specification" also defined in JSON. A specification is comprised of one or more "operations". See
Specification Support, below, for more details.

## Documentation
API Documentation is available at http://godoc.org/gopkg.in/qntfy/kazaam.v3.

## Features
Kazaam is primarily designed to be used as a library for transforming arbitrary JSON.
It ships with eleven built-in transform types, and twenty-one built-in converter types,
described below, which provide significant flexibility in reshaping JSON data.

Also included when you `go get` Kazaam, is a binary implementation, `kazaam` that can be used for
development and testing of new transform specifications.

Finally, Kazaam supports the implementation of custom transform and converter types. We encourage and appreciate
pull requests for new transform types so that they can be incorporated into the Kazaam distribution,
but understand sometimes time-constraints or licensing issues prevent this. See the API documentation
for details on how to write and register custom transforms.

Due to performance considerations, Kazaam does not fully validate that input data is valid JSON. The
`IsJson()` function is provided for convenience if this functionality is needed, it may significantly slow
down use of Kazaam.

## Transform Specification Support

Transforms are the main mechanism in Kazaam for shaping json documents. Transforms, unlike converters work
at the document level, whereas converters work at the value level.  There are many built-in transforms for
you to shape your document, but there is also a mechanism for developing your own custom transforms when the
need arises.

Kazaam currently supports the following built-in transforms:
- shift
- steps
- concat
- coalesce
- extract
- timestamp
- uuid
- default
- pass
- delete
- merge

### Shift
The shift transform is the current Kazaam workhorse used for remapping of fields. It supports a `"require"` field that when 
set to `true`, will throw an error if *any* of the paths in the source JSON are not present.

The shift transform by default is destructive. For in-place operation, an optional `"inplace"`
field may be set.

The specification supports jsonpath-esque JSON accesses and sets as well as a custom JSON Path Parameters.

Concretely
```json
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
```json
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
```json
{
  "object": {
    "id": 12345
  },
  "gid2": "guid2",
  "allGuids": ["guid0", "guid2", "guid4"]
}
```
##### JSON Path Syntax

The JSON Path implementation supports a few special cases:
- *Array accesses*: Retrieve `n`th element from array
- *Array wildcarding*: indexing an array with `[*]` will return every matching element in an array
- *Top-level object capture*: Mapping `$` into a field will nest the entire original object under the requested key
- *Array append/prepend and set*: Append and prepend an array with `[+]` and `[-]`. Attempting to write an array element that does not exist results in null padding as needed to add that element at the specified index (useful with `"inplace"`).
- *JSON Path Parameters*: Conditional Expressions and chained value conversions through Converter expressions

##### JSON Path Parameters

###### JSON Path Conditional Expressions

JSON Path Conditionals allow value skipping based on document existence or conditional expression evaluation and
take on the following forms:

| Path Structure              | Description    |
|:--------------------------- |:---------------|
| _path.existing.value_ **?** | Return the existing value or skip the value if it is not defined. (NOTE: skipping is allowed with conditionals even when the `"require"` option is used with Shift) |
| _path.missing.value_ **? "default value"**<br/>_path.missing.value_ **? 42** | Default values can be provided, and when a value is missing, the default value that was provided is returned instead. |
| _path.existing.value_ **? ston("other.value") > 3 && another.value == "test" :** | Existing value is skipped unless the *Conditional Expression* evaluates to `"true"`. **Note**: the colon is required here or the expression itself will attempt to be treated as a default value. |
| _path.value_ **? other.value == "something" : "default value"** | If the path exists and the expression evaluates to true, the existing value is returned. If the path is missing, the default is provided. Otherwise, if the expression evaluates to `"false"` the default is returned. |

**Notes**:
* Default values are simple JSON Values only (no composites). Strings must be quoted (and when embedded in JSON, the quote will need to be escaped.) Strings, Boolean, Nulls and Numbers are supported.

  e.g.
  `"gid2 ? \"default value\"": "guid2",`

**Notes**: 
* Function calls in Conditional Expressions call to named Converters and require 1 or 2 string parameters. The first parameter must be a JSON Path (without parameters) as a string to
a value that will be converted and the optional second parameter must be a string. If provided the arguments will be provided to the Converter as a single string for it to parse.

  e.g. 
  `"gid2 ? substr(\"guid2\",\"2 3\") == \"id\":": "guid2",`

###### JSON Path Converter Expressions

JSON Path Converter Expressions allow for values to be altered by chaining the existing (or default value) through Converter functions. The value
returned from the last Converter in the chain will become the returned value for the JSON Path query.

| Path Structure | Description |
|:---------------|:------------|
| _path.existing.value_ **&#124; converter1 arguments &#124; converter2 arguments**| Chained Converter syntax | 
| _path.value_ **? other.value == "something" : "default value" &#124; converter1** | Can be combined with Conditional Expressions |

**Notes**: 
* If **`|`** characters are required as part of the value, they can be escaped with a **`\\`** character, and **`\\`** characters themselves
can also be escaped.

**Notes**: 
* The whitespace between the converter name and arguments, as well as surrounding the argument is ignored.  Although whitespace
within the arguments are preserved, if the whitespace around the arguments is required, it must be escaped:

  e.g.
  `"path.value | converter1 \ arguments\ `" would cause **` arguments `** to be the arguments string.

Arguments are passed to the converter functions as a single string, and will require the converter function to parse out any meaningful parameters.


### Steps
The steps transform performs a series of shift transforms with each step working on the ouptput from the last step. This
transform is very similar to the shift transform, and takes the same optional parameters.

The following example produces the same results as the `Shift` transform example presented earlier. The only difference
is that the each of the steps are guaranteed to transform in the specified order.

```json 
{
  "operation": "steps",
  "spec": {
    "steps": [
      {
        "object.id": "doc.uid"
      },
      {
        "gid2": "doc.guid[1]"
      },
      {
        "allGuids": "doc.guidObjects[*].id"
      }
    ]
  }
}
```


### Concat
The concat transform allows the combination of fields and literal strings into a single string value.
```json
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
```json
{
    "a": {
        "timestamp": 1481305274
    }
}
```

would result in
```json
{
    "a": {
        "timestamp": "TEST,1481305274"
    }
}
```

**Notes**:
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
```json
{
  "operation": "coalesce",
  "spec": {
    "firstObjectId": ["doc.guidObjects[0].uid", "doc.guidObjects[0].id"]
  }
}
```

executed on a json message with format
```json
{
  "doc": {
    "uid": 12345,
    "guid": ["guid0", "guid2", "guid4"],
    "guidObjects": [{"id": "guid0"}, {"id": "guid2"}, {"id": "guid4"}]
  }
}
```

would result in
```json
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
```json
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
```json
{
  "operation": "extract",
  "spec": {
    "path": "doc.guidObjects[0].path.to.subobject"
  }
}
```

executed on a json message with format
```json
{
  "doc": {
    "uid": 12345,
    "guid": [
      "guid0",
      "guid2",
      "guid4"
    ],
    "guidObjects": [
      {
        "path": {
          "to": {
            "subobject": {
              "name": "the.subobject",
              "field": "field.in.subobject"
            }
          }
        }
      },
      {
        "id": "guid2"
      },
      {
        "id": "guid4"
      }
    ]
  }
}
```

would result in
```json
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
```json
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
```json
{
  "timestamp": [
    "Sat Jul 22 08:15:27 +0000 2017",
    "Sun Jul 23 08:15:27 +0000 2017",
    "Mon Jul 24 08:15:27 +0000 2017"
  ]
}
```

would result in
```json
{
  "timestamp": [
    "2017-07-22T08:15:27+0000",
    "Sun Jul 23 08:15:27 +0000 2017",
    "Mon Jul 24 08:15:27 +0000 2017"
  ],
  "nowTimestamp": "2017-09-08T19:15:27+0000"
}
```

### UUID
A `uuid` transform generates a UUID based on the spec. Currently supports UUIDv3, UUIDv4, UUIDv5.

For version 4 is a very simple spec

```json
{
    "operation": "uuid",
    "spec": {
        "doc.uuid": {
            "version": 4
        }
    }
}
```

executed on a json message with format
```json
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
```json
{
  "doc": {
    "author_id": 11122112,
    "document_id": 223323,
    "meta": {
      "id": 23
    },
    "uuid": "f03bacc1-f4e0-4371-a5c5-e8160d3d6c0c"
  }
}
```

For UUIDv3 & UUIDV5 are a bit more complex. These require a Name Space which is a valid UUID already, and a set of paths, which generate UUID's based on the value of that path. If that path doesn't exist in the incoming document, a default field will be used instead. **Note** both of these fields must be strings.
**Additionally** you can use the 4 predefined namespaces such as `DNS`, `URL`, `OID`, & `X500` in the name space field otherwise pass your own UUID.

```json
{
   "operation":"uuid",
   "spec":{
      "doc.uuid":{
         "version":5,
         "namespace":"DNS",
         "names":[
            {"path":"doc.author_name", "default":"some string"},
            {"path":"doc.type", "default":"another string"}
         ]
      }
   }
}
```

executed on a json message with format
```json
{
  "doc": {
    "author_name": "jason",
    "type": "secret-document",
    "document_id": 223323,
    "meta": {
      "id": 23
    }
  }
}
```

would result in
```json
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
```json
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
```json
{
  "operation": "delete",
  "spec": {
    "paths": ["doc.uid", "doc.guidObjects[1]"]
  }
}
```

executed on a json message with format
```json
{
  "doc": {
    "uid": 12345,
    "guid": ["guid0", "guid2", "guid4"],
    "guidObjects": [{"id": "guid0"}, {"id": "guid2"}, {"id": "guid4"}]
  }
}
```

would result in
```json
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

### Merge
A merge transform will take multiple arrays and join them in to an array of objects joining them by keys. The arrays should be equal length.

```json
{
  "operation": "merge",
  "spec": {
    "merge1": [
      {
        "name": "prop_1",
        "array": "array_a"
      },
      {
        "name": "prop_2",
        "array": "array_b"
      },
      {
        "name": "prop_3",
        "array": "array_c"
      }
    ]
  }
}
```

executed on a json message with format:
```json
{
  "array_a": [
    "a_1",
    "a_2",
    "a_3"
  ],
  "array_b": [
    "b_1",
    "b_2",
    "b_3"
  ],
  "array_c": [
    "c_1",
    "c_2",
    "c_3"
  ]
}
```

would result in:
```json
{
  "merge1": [
    {
      "prop_1": "a_1",
      "prop_2": "b_1",
      "prop_3": "c_1"
    },
    {
      "prop_1": "a_2",
      "prop_2": "b_2",
      "prop_3": "c_2"
    },
    {
      "prop_1": "a_3",
      "prop_2": "b_3",
      "prop_3": "c_3"
    }
  ]
}
```



## Converter Specification Support

Converters in Kazaam allow for value level transformations and work within and extend the current Transform
capabilities.

Kazaam currently supports the following built-in Conveters:

| Converter Name | Description |
|---------------:|:------------|
`add <num>` | adds a number value to a number value
`ceil` | converts the number value to the least integer greater than or equal to the number value
`div <num>` | divides a number value by a number value
`floor` | converts the number value to the greatest integer less than or equal to the number value
`format <string>` | converts the value to a string via a **`fmt`** string
`lower` | converts the string value to lowercase characters
`mapped <string>`| maps the string value to another string value using predefined named maps
`mul <num>` | multiples a number value by a number value
`ntos` | converts the number value to a string value
`regex` | alters the string value with named regex replacements
`round` | converts a number value to the closet integer value
`ston` | converts a string value to a number value
`substr <num> [<num>]` | converts a string value to a substring value
`trim` | converts a string value by removing the leading and trailing whitespace characters
`upper` | converts a string value to uppercase characters
`len` | converts a string to an integer value equal to the length of the string
`splitn <string> <num>` | splits a string by a delimiter string and returns the Nth token (1 based)
`eqs <any>` | returns `true` or `false` based on whether the value matches the parameter
`not` | returns `true` if value is `false` and `false` if the value is anything other than `false`
`split <delim>` | returns array of values split on delimiter
`join <delim>` | joins an array of strings by the delimiter

### Converter Examples ###

The following examples will use the same input JSON value:

```json 
{
  "tests": {
    "test_int": 500,
    "test_float": 500.01,
    "test_float2": 500.0,
    "test_number": "750",
    "test_fraction": 0.5,
    "test_trim": "    blah   ",
    "test_money": "$6,000,000",
    "test_chars": "abcdefghijklmnopqrstuvwxyz",
    "test_mapped": "Texas",
    "test_null": null,
    "pi": 3.141592653,
    "test_true": true,
    "test_false": false,
    "test_null": null,
    "test_string": "The quick brown fox",
    "test_naics_code": "531312",
    "test_split":"a|b|c",
    "test_join":["a","b","c"]
  },
  "test_bool": true
}
```

#### Add ####

Adds a number to a number value

Argument | Description
---------|------------
Number   | Number value to add

example:
```json
{
  "operation": "shift",
  "spec": {
    "output": "tests.test_int | add 1"
  }
}
```

produces:
```json
{
  "output": 501
}
```

#### Ceil ####

Converts a number value to the least closest integer greater than or equal to the number value

Argument | Description
---------|------------


example:
```json
{
  "operation": "shift",
  "spec": {
    "output": "tests.test_float | ceil"
  }
}
```

produces:
```json
{
  "output": 501
}
```

#### Div ####

Divides a number value by another number value

Argument | Description
---------|------------
Number   | dividend in a division operation

example:
```json
{
  "operation": "shift",
  "spec": {
    "output1": "tests.test_float | div 2",
    "output2": "tests.test_int | div .5"
  }
}
```

produces:
```json
{
  "output1": 250,
  "output2": 1000
}
```

#### Floor ####

Converts a number value to the greatest integer value less than or equal to the number value

Argument | Description
---------|------------

example:
```json
{
  "operation": "shift",
  "spec": {
    "output": "tests.test_float | floor"
  }
}
```

produces:
```json
{
  "output": 500
}
```

#### Format ####

Formats a value into a new string value using a **`fmt`** string

Argument | Description
---------|------------
string   | fmt string, if whitespace shouldn't be trimmed, it should be escaped with \\

example:
```json
{
  "operation": "shift",
  "spec": {
    "output1": "tests.pi | format %.4f",
    "output2": "tests.test_float | format %.0f",
    "output3": "tests.test_string | format %s jumps over the lazy dog",
    "output4": "tests.test_true | format %t is the value"
  }
}
```

produces:
```json
{
  "output1": "3.1416",
  "output2": "500",
  "output3": "The quick brown fox jumps over the lazy dog",
  "output4": "true is the value"
}
```

#### Lower ####

Converts a string value to lowsercase

Argument | Description
---------|------------

example:
```json
{
  "operation": "shift",
  "spec": {
    "output": "tests.test_string | lower"
  }
}
```

produces:
```json
{
  "output": "the quick brown fox"
}
```

#### Mapped ####

Maps a string value to another string value using a named JSON map defined in `$.converters.mapped`

Argument | Description
---------|------------
string   | name of the map to use

example:
```json
{
  "operation": "shift",
  "converters": {
    "mapped": {
      "states": {
        "Ohio": "OH",
        "Texas": "TX"
      }
    }
  },
  "spec": {
    "output": "tests.test_mapped | mapped states"
  }
}
```

produces:
```json
{
  "output": "TX"
}
```

#### Mul ####

Multiplies a number value by another number value

Argument | Description
---------|------------
Number   | multiplier value of a multiplication operation

example:
```json
{
  "operation": "shift",
  "spec": {
    "output1": "tests.test_int | mul 2",
    "output2": "tests.test_int | mul .5"
  }
}
```

produces:
```json
{
  "output1": 1000,
  "output2": 250
}
```

#### Ntos ####

Converts a number value to a string value

Argument | Description
---------|------------

example:
```json
{
  "operation": "shift",
  "spec": {
    "output": "tests.test_int | ntos"
  }
}
```

produces:
```json
{
  "output": "500"
}
```

#### Regex ####

Use Regexp ReplaceAll to match and replace values defined in the `$.converters.regex` configuration object. You can
also pass an array of configuration objects and they will all be applied in order, stopping after the first match is matched and replaced.

Argument | Description
---------|------------
string   | name of predefined regex match and replace

example:
```json
{
  "operation": "shift",
  "converters": {
    "regex": {
      "remove_dollar_sign": {
        "match": "\\$\\s*(.*)",
        "replace": "$1"
      },
      "remove_comma": {
        "match": ",",
        "replace": ""
      },
      "convert_naics": [
      	{
      		"match": "^8111.*",
       		"replace": "automotive_services"
      	},
      	{
      		"match": "^4413.*",
       		"replace": "automotive_services"
      	},
      	{
      		"match": "^531.*",
      		"replace": "real_estate"
      	}
      ]
    }
  },
  "spec": {
    "output": "tests.test_money | regex remove_dollar_sign | regex remove_comma"
  }
}
```

produces:
```json
{
  "output": "6000000"
}
```

#### Round ####

Rounds a number value to the closest integer value

Argument | Description
---------|------------

example:
```json
{
  "operation": "shift",
  "spec": {
    "output1": "tests.test_float | round",
    "output2": "tests.test_fraction | round"
  }
}
```

produces:
```json
{
  "output1": 500,
  "output2": 1
}
```

#### Ston ####

Converts a string value to a number value

Argument | Description
---------|------------

example:
```json
{
  "operation": "shift",
  "spec": {
    "output": "tests.test_number | ston"
  }
}
```

produces:
```json
{
  "output": 750
}
```

#### Substr ####

Returns a substring of a string value

Argument | Description
---------|------------
number   | 0 based index where to start the substring
number   | (optional) index of last character + 1 in the substring, if omitted uses the string's length

example:
```json
{
  "operation": "shift",
  "spec": {
    "output1": "tests.test_chars | substr 3 6",
    "output2": "tests.test_string | substr 10"
  }
}
```

produces:
```json
{
  "output1": "def",
  "output2": "brown fox"
}
```

#### Trim ####

Removes whitespace from the start and end of a string value

Argument | Description
---------|------------

example:
```json
{
  "operation": "shift",
  "spec": {
    "output": "tests.test_trim | trim"
  }
}
```

produces:
```json
{
  "output": "blah"
}
```

#### Upper ####

Converts a string value to uppercase

Argument | Description
---------|------------

example:
```json
{
  "operation": "shift",
  "spec": {
    "output": "tests.test_string | upper"
  }
}
```

produces:
```json
{
  "output": "THE QUICK BROWN FOX"
}
```

#### Len ####

Returns the length of a string value

Argument | Description
---------|------------

example:
```json
{
  "operation": "shift",
  "spec": {
    "output": "tests.test_string | len"
  }
}
```

produces:
```json
{
  "output": 19
}
```

#### Splitn ####

Returns the Nth token of a string split by a delimiter string

Argument | Description
---------|------------
string   | delimiter string
number   | one based position of token to return

example:
```json
{
  "operation": "shift",
  "spec": {
    "output": "tests.test_string | splitn o 2"
  }
}
```

produces:
```json
{
  "output": "wn f"
}
```

#### Eqs ####

Returns `true` or `false` based on whether the value equals the parameter

Argument | Description
---------|------------
any      | value to compare


example:
```json
{
  "operation": "shift",
  "spec": {
    "output": "tests.test_string | eqs \"The quick brown fox\""
  }
}
```

produces:
```json
{
  "output": true
}
```

#### Not ####

Negates a `false` value returning `true` and returns `false` for everything else

Argument | Description
---------|------------



example:
```json
{
  "operation": "shift",
  "spec": {
    "output1": "tests.test_true | not",
    "output2": "tests.test_false | not"
  }
}
```

produces:
```json
{
  "output1": false,
  "output2": true
}
```

#### Split ####



Argument | Description
---------|------------
delim    | string delimiter on which to split the string 


example:
```json
{
  "operation": "shift",
  "spec": {
    "output1": "tests.test_split | split \\|"
  }
}
```

produces:
```json
{
  "output1": ["a","b","c"]
}
```

#### Join ####



Argument | Description
---------|------------
delim    | string delimiter on which to join the array into a string 


example:
```json
{
  "operation": "shift",
  "spec": {
    "output1": "tests.test_join | join \\|"
  }
}
```

produces:
```json
{
  "output1": "a|b|c"
}
```


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
cd $GOPATH/src/gopkg.in/qntfy.kazaam.v3/kazaam
go install
```

This will create an executable in `$GOPATH/bin` like you
would expect from the normal `go` build behavior.

### Examples

See [godoc examples](https://godoc.org/pkg/gopkg.in/qntfy/kazaam.v3/#pkg-examples).
