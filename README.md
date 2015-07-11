# Gohn -- Simple Hacker News Go client.
A simple Go client for the Hacker News API. Items are parsed into protocol buffers.

## Example
Download Hacker News items 100 to 1000 with 30 workers to a leveldb database:
```
go run download/main.go --start 100 --end 1000 --numworkers=30
```

Iterate through the leveldb database and print the text:
```
go run readdb/read.go 
```
