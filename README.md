# REST API Test

The `rest-test` command line tool is used to test and validate REST API.

## Features

* `YAML` based test definition and validation
  * Test Plan -> Test Suite -> Test Case
* Optimized for JSON, file and plain response body
* JSON response Body value filter
  * Support [JSON path](https://github.com/PaesslerAG/jsonpath) value filter
  * Support [gval](https://github.com/PaesslerAG/gval) expression
  * Customized [gval](https://github.com/PaesslerAG/gval) operators
    * `string` convert a number, boolean value to string
    * `int` convert a number string to int
    * `float` convert a number string to float
    * `bool` convert a number, string, list and map to bool
    * `len` get length of a string, a map or a slice/array
    * `contain` check whether an element in a list or map. Check whether a string is a substring of another string
* `JavaScript` to evaluate values
* `YAML` report

> For more details please refer to [REST API Test](https://exia.dev/project/rest-test/)