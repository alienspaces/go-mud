# ./server

This is the top level directory for all server side code.

We create this directory so the project then does not need to reorganise its structure when adding additional top level directories such as `./ui` or `./ci`.

## Subdirectories

| Directory     | Description                                                     |
| ------------- | --------------------------------------------------------------- |
| `./build`     | Project build and deployment configuration                      |
| `./client`    | Project HTTP client implementations for accessing service API's |
| `./constant`  | Project level constants                                         |
| `./core`      | Core package containing all common source code                  |
| `./schema`    | Project level JSON Schema definitions for API services          |
| `./script`    | Developer scripts for common tasks                              |
| `./service`   | Project service implementations                                 |
