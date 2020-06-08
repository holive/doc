# doc

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

APIs specifications in one place with ReDoc.

## how to run
```bash
go get github.com/holive/doc
cd $GOPATH/src/github.com/holive/doc
make run
```

## how it works

### create
Endpoint: `POST http://localhost:8080/{squad}/{project}/{version}`
```bash
curl -F "fileupload=@api-document.yaml" http://localhost:8080/matrix/new/v1
```

### delete
Endpoint: `DELETE http://localhost:8080/{squad}/{project}/{version}`

### see one
Endpoint: `GET http://localhost:8080/{squad}/{project}/{version}`

### see all
Endpoint: `GET http://localhost:8080/`
