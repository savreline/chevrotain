# Links and References

## Golang
### Go Net and RPC Packages
https://golang.org/pkg/net/ \
https://godoc.org/google.golang.org/grpc \
https://github.com/grpc/grpc-go

### Go IO
https://golang.org/pkg/io/ioutil/#WriteFile

### Go HTML
https://godoc.org/golang.org/x/net/html

### Go Scope Rules
https://www.tutorialspoint.com/go/go_scope_rules.htm

## MongoDB
### Connecting to MondoDB:
https://www.mongodb.com/golang \
https://github.com/mongodb/mongo-go-driver \
https://godoc.org/go.mongodb.org/mongo-driver/mongo#Collection.Find \
https://stackoverflow.com/questions/56970867/how-to-use-mongo-driver-connection-into-other-packages \
https://www.mongodb.com/blog/post/mongodb-go-driver-tutorial-part-1-connecting-using-bson-and-crud-operations \
https://docs.mongodb.com/manual/reference/operator/query/elemMatch/

### MongoD commands:
https://docs.mongodb.com/manual/tutorial/manage-mongodb-processes/

### Mongo Shell commands:
https://docs.mongodb.com/manual/reference/mongo-shell/
```
mongo --port 27017
mongod --config "C:\Program Files\MongoDB\Server\4.0\bin\replica1.cfg"
cls
show dbs
show collections
use chevrotain
chevrotain.createCollection("kvs")
db.kvs.find()
db.kvs.count()
db.dropDatabase()
```

Save vs. Insert vs. Update \
https://stackoverflow.com/questions/16209681/what-is-the-difference-between-save-and-insert-in-mongo-db

Context \
https://golang.org/pkg/context/

## Other
### Markdown Syntax
https://www.markdownguide.org/cheat-sheet/

### Git
https://stackoverflow.com/questions/1186535/how-to-modify-a-specified-commit
```
git ls-files | xargs wc -l
wc -l $(git ls-files | grep '.*\.cs')
```

### BFS
https://stackoverflow.com/questions/10258305/how-to-implement-a-breadth-first-search-to-a-certain-depth

### Windows
https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/del
