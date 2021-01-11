# Go test tool


## Basic
- file name must ends with `_test.go`, it won't get compiled when build the binary executable file.
- must import the `testing` pkg.
- Format: `func TestXxx(t *testing.T)`, Xxx must start with Upper case. Usually use the format `StructName_functionName(t *testing.T)`

## How to run specified test
- `go test <package_name>` it runs all test cases in the package.
- `go test foo_test.go` it runs all the test cases in the file
    + If `foo_test.go` and `foo.go` are the same package (a common case) then you must name all other files required to build foo_test. In this example it would be:
      ```bash
      $ go test foo_test.go foo.go
      ```
- `-run` flag
    ```
    -run regexp
        Run only those tests and examples matching the regular
        expression.
    ```



## Reference
- [pkg test](https://golang.org/pkg/cmd/go/internal/test/)

