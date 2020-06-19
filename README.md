# doc

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

APIs specifications in one place with ReDoc.

## How to run
```bash
go get github.com/holive/doc
cd $GOPATH/src/github.com/holive/doc/app/cmd
go run main.go
```

There is also a Docker image:

```bash
docker run -e MONGO_CONNECTION_STRING="your-connection-string" -p 3000:3000 hbliveira/doc
```

## How it works

**You need a squad key to be able to create or delete a doc**. 
So get your squad key:
```bash
curl -X POST --data '{"name": "<your-squad-name>"}' http://localhost:3000/squad
```

### Create
Endpoint: `POST http://localhost:3000/{project}/{version}`
```bash
curl -F 'fileupload=@redoc.yaml' -F 'squad=<your-squad-name>' -F "descricao=<optional-description>" http://localhost:3000/{project}/{version} -H 'X-Squad-Key: <your-squad-key>'
```

### Delete
Endpoint: `DELETE http://localhost:3000/{project}/{version}`
```bash
curl -X DELETE http://localhost:3000/{project}/{version} -F 'squad=<your-squad-name>' -H 'X-Squad-Key: <your-squad-key>'
```
 
### Get
Endpoint: `GET http://localhost:3000/{project}/{version}`

### List
Endpoint: `GET http://localhost:3000`
