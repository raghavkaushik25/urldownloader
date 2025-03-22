# URL Downloader  

## Overview  

This project is a concurrent URL downloader written in Go. It follows a three-stage pipeline to efficiently process a CSV file containing URLs, download their content, and persist the data to disk.  

## Flow

### Stage 1 - Read File  
- Runs in a separate goroutine.  
- Reads the CSV line by line and ingnores the header.
- As soon as a new URL is read from the csv. It starts a new go routine to download its contents. 

### Stage 2 - Download Contents  
- Uses semaphore channel buffered to maxNumOfGoroutines to limit the Goroutines.
- As soon as the number "Download" gorouritnes reach 50, it will stop spawning more and wait for the existing one to return.
- It will only start max 50 concurrent goroutines. 

### Stage 3 - Persist Contents  
- Runs in a single dedicated goroutine which is started from main().  
- Writes the downloaded data to disk with url.Host + time.Now() as the file name.  
- Ensures data consistency by avoiding multiple writes at the same time.  


## Architecture

- Uses a robust logging mechanism to ensure suffiecient DEBUG logs if in case debugging is required. 

- Ensure that program shuts down gracefully if in case it receives a os.Interupt from the system.

- Uses locks to maintain the stats.


## How to run  

## Clone

Make sure you have Go installed on your system, preferably go 1.23 or above. Then, clone this repository and navigate to the project folder:  

```
git clone https://github.com/raghavkaushik25/urldownloader.git
cd urldownloader
```

## Build

```go build .```

## Run
```./url-downloader -path="URLs.csv"```

- This will create an output folder in the current working directory and put all the files in that folder.
- You can change the path to whatever the path of your csv is.

```./url-downloader -path="URLs.csv" -debug```

- This will run the application in debug mode and you will get debug logs on STDOUT. This is recommended if in case you are seeing descrapencies.

## Unit Test

```go test ./... ```

```
To check for race conditions:
   go test -race 
```
