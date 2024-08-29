# REST API Test

The `rest-test` command line tool is used to test REST API.

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
    * `contain` check whether an element in a list or map. Check whether a string is a substring of another string
* `JavaScript` to evaluate values
* `YAML` report
