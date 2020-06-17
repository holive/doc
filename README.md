# doc

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

APIs specifications in one place with ReDoc.

## how to run
```bash
go get github.com/holive/doc
cd $GOPATH/src/github.com/holive/doc/app/cmd
go run main.go
```

There is also a Docker image:

```bash
docker run -p 3000:3000 hbliveira/doc
```

## how it works

*You need a squad key to be able to create or delete a doc. So get your squad key (replace <squad-name> by whatever):
```bash
curl -X POST --data '{"name": "<squad-name>"}' http://localhost:3000/squad
```

### create
Endpoint: `POST http://localhost:3000/{your-squad-name}/{project}/{version}`
```bash
curl -F "fileupload=@redoc.yaml" http://localhost:3000/{your-squad-name}/{project}/{version} -H 'x-squad-key: <your-squad-key>'
```

### delete
Endpoint: `DELETE http://localhost:3000/{your-squad-name}/{project}/{version}`
```bash
curl -X DELETE http://localhost:3000/{your-squad-name}/{project}/{version} -H 'x-squad-key: <your-squad-key>'
```
  
### get one
Endpoint: `GET http://localhost:3000/{your-squad-name}/{project}/{version}`

### get all
Endpoint: `GET http://localhost:3000/`
