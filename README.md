# RESTful_API

Basic CRU (yes, without Delete) API implemented in Go. Used to accept different texts and keep track of the number of appearances of each word.
  
For high traffic it can be scaled up to 5 instances with a load balancer. The balancer has the round robin technique implemented for simplicity.

The words appearances are *stored in memory and persisted in a file*. Each instance has its own dedicated file.

After a restart, if the files are not deleted, the instances reload their memory with the data from the file.

For synchronization, when one instance receives the HTTP request, it will update its memory and its file *and then will send an RPC call to the other instances* in order to keep them updated and respect the eventual consistency.

Supported requests:

**POST**
```
POST http://localhost:7000
Body:
{
  "Text": "this is a sample text"
}
```

**GET**
```
GET http://localhost:7000
Body:
{
  "Words": ["this", "is"]
}

A successful response for this would be:
Status: 200 OK
Body:
{"this": 1, "is": 1}
```

### Build and test

Requierements: go 1.21, linux or git bash.

1. Clone the project
2. Go to server and then balancer directory and run ```go build .```
3. To run the unit test, in the same 2 directories run ```go test``` or ```go test -v``` to see each test

### Run

1. Direct way: run ```./run.sh <balancer_port> <number_of_instances>``` which will build the 2 modules and then start <number_of_instances> API instances and a load balancer on the specified <balancer_port>
2. Steps with load balancer: see **Build and test**, build the 2 modules then run  ```./server -port=<port> -no_of_instances=<no_of_instances>``` in server directory and  ```./balancer -port=<port> -number_of_instances=<number_of_instances>``` in balancer directory.
3. Just one instance: multiple instances are not mandatory, you can just run ```./server -port<port>```




