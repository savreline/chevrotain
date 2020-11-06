# Chevrotain
`go run chevrotain.go`

### Connecting to MondoDB:
https://www.mongodb.com/golang \
https://github.com/mongodb/mongo-go-driver \
https://stackoverflow.com/questions/56970867/how-to-use-mongo-driver-connection-into-other-packages
https://www.mongodb.com/blog/post/mongodb-go-driver-tutorial-part-1-connecting-using-bson-and-crud-operations
https://godoc.org/go.mongodb.org/mongo-driver/mongo#Collection.Find

### Go Net Package
`time.Sleep(5 * time.Second)` \
https://golang.org/pkg/net/

### Go IO
https://golang.org/pkg/io/ioutil/#WriteFile

### Go HTML
https://godoc.org/golang.org/x/net/html

### Go Scope Rules:
https://www.tutorialspoint.com/go/go_scope_rules.htm

### Mongo Shell commands:
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
```
git ls-files | xargs wc -l
wc -l $(git ls-files | grep '.*\.cs')
```
https://docs.mongodb.com/manual/reference/mongo-shell/
https://docs.mongodb.com/manual/reference/method/db.createCollection/

Save versus Insert versus Update \
https://stackoverflow.com/questions/16209681/what-is-the-difference-between-save-and-insert-in-mongo-db

### MongoD commands:
`mongod`\
https://docs.mongodb.com/manual/tutorial/manage-mongodb-processes/

### Context:
https://golang.org/pkg/context/

### Vector Clocks
https://godoc.org/github.com/DistributedClocks/GoVector/govec

### Markdown Syntax
https://www.markdownguide.org/cheat-sheet/

### Git
`git reset --soft HEAD~1`
