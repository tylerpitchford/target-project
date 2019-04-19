# Overview

As specified, this project implements a basic text search ( `O(MN) time` ), a regex search ( `O(MN) time` ), and an index search ( `O(1) time` ). The project also implements a positional-index search feature ( `O(MN) time` ). It wasn't requested, but it opens up the possibility of better matching capabilities, distance-based queries, and other exotic queries. It relies upon nested hash sets to store the positional index data. 

Finally, all search options have been implemented in concurrent and non-concurrent forms.

# Performance considerations

As shown by the benchmarks below, the single-token index search is the fastest option by far. Due to the precompilation of the indexes, the lookup time is a single hash retrieval which is `O(1)` complexity. Space wise, the storage required is O(N) in the worst case assuming a completely unique corpus.

# Notes

The text search relies upon golang's `strings.Count` function which produces partial matches (i.e. The matches both The and There). Similarly, the regex search also produces partial matches as compared to whole word matches.

As for the indexers, the do not support partial matches. The single-token indexer tokenizes based upon whitepsace, punctuation, and some special conditions for quoted text and numbers. The positional-indexer tokenizes on only punctuation and whitespace.

# Real-world Optimizations and TODOs

  1. *Don't regenerate indexes on every search* - It's pretty obvious that the indexes should only be rebuilt when something changes. I regenerate them every time in this code to make debugging and changing the index formats easier, but in production, the indexes should remain as static as possible between runs.
  
  2. *Caching* - Assuming a larger corpus and non-random searching, caching results could greatly enhance performance times at the cost of extra memory utilization.
  
  3. *Map-reduce* - Again, assuming a larger corpus, instead of splitting functions into merely concurrent processing on a local machine, searches in each file could be split into map-reduce functions. You could split up parts of files into map-reduce searches, but you would need to be careful to handle possible matches between the overlap of the data buffers.

  4. *Parallel position look-ups* - The positional-indexer relies upon recursive positional look-ups. It, however, would be possible to load the positions for each token in parallel and then reduce them together to potentially speed up performance. Depending on the hit and miss rates of the queries you could be generating a lot of needless execution for look-ups that would have terminated early, but mixed with caching and some usage data on miss rates, it could be an option for converting the serial logic into parallel performance at the cost.
  
  5. *Documentation* - There are some inline comments in some areas but if this were to be consumed as a package, it would require better documentation. Further, I would want to refactor a few areas that I streamlined for screen printing vs. consumption as a package.
  
  6. *Unit-testing* - You can always use more unit tests.
  
# Command-line options

```
> ./target-project -h
  -benchmark
    	Run the benchmarks.
  -concurrent
    	Run the search concurrently.
  -directory string
    	Provide a directory where files should be searched or indexed. Only files with the extension .txt are considered. (default "data")
  -positional
    	Use a positional search indices.
  -token string
    	Provide the search token non-interactively.
  -type int
    	Provide the search type non-interactively. (default -1)
```

# Example usage

## Interactive Mode
```
> ./target-project
Enter the search term: of the

Search Method: 1) String Match 2) Regular Expression 3) Indexed: 1

	 hitchhikers.txt - 7 matches

	 french_armed_forces.txt - 6 matches

	 warp_drive.txt - 1 matches

Elapsed time: 147.657µs
```

## Interactive Mode ( concurrent / positional )
```
> ./target-project -concurrent -positional
Enter the search term: of the

Search Method: 1) String Match 2) Regular Expression 3) Indexed: 3

	 hitchhikers.txt - 6 matches

	 french_armed_forces.txt - 6 matches

	 warp_drive.txt - 1 matches

Elapsed time: 1.542702ms
```

## Non-interactive Execution
```
./target-project -token="of the" -type=3 -positional -concurrent
	 hitchhikers.txt - 6 matches

	 french_armed_forces.txt - 6 matches

	 warp_drive.txt - 1 matches

Elapsed time: 306.279µs
```

## Benchmark Mode (2M searches)
```
> ./target-project -benchmark
2019/04/19 15:58:20 Starting benchmarks.
2019/04/19 15:58:20 
2019/04/19 15:58:20 Text Search (non-concurrent)
2019/04/19 15:58:28 Total time:  8.157646504s
2019/04/19 15:58:28 Text Search (concurrent)
2019/04/19 15:58:45 Total time:  17.110769586s
2019/04/19 15:58:45 RegEx Search (non-concurrent)
2019/04/19 15:59:07 Total time:  21.622161739s
2019/04/19 15:59:07 RegEx Search (concurrent)
2019/04/19 15:59:47 Total time:  39.763637052s
2019/04/19 15:59:47 Index Search (single-token, non-concurrent)
2019/04/19 15:59:48 Total time:  1.377887351s
2019/04/19 15:59:48 Index Search (single-token, concurrent)
2019/04/19 15:59:56 Total time:  8.252935719s
2019/04/19 15:59:56 Index Search (positional, non-concurrent)
2019/04/19 16:00:37 Total time:  40.875302939s
2019/04/19 16:00:37 Index Search (positional, concurrent)
2019/04/19 16:01:12 Total time:  34.862793409s
2019/04/19 16:01:12 
2019/04/19 16:01:12 Benchmarks complete.
```

## Golang Unit-testing 
```
go test -v ./...
=== RUN   TestNegativeSearchType
--- PASS: TestNegativeSearchType (0.00s)
=== RUN   TestTooLargeSearchType
--- PASS: TestTooLargeSearchType (0.00s)
=== RUN   TestStringSearchType
--- PASS: TestStringSearchType (0.00s)
=== RUN   TestValidSearchTypes
--- PASS: TestValidSearchTypes (0.00s)
PASS
ok  	target-project	0.006s
```

# Binary Releases 

Binary releases are available for 64-bit windows, osx, and linux operating systems:

 * [Windows](https://github.com/tylerpitchford/target-project/releases/download/v1/target-windows-amd64.exe)
 * [OSX](https://github.com/tylerpitchford/target-project/releases/download/v1/target-darwin-amd64)
 * [Linux](https://github.com/tylerpitchford/target-project/releases/download/v1/target-linux-amd64)

Instructions on building or running your own instances in Docker are below.

# Building the project

## Golang Native
```
cd GOPATH

git clone https://github.com/tylerpitchford/target-project

cd target-project

go build
```

## Golang Native 
( *multiple architectures - windows/64-bit; linux/64-bit, osx/64-bit* )
```
cd GOPATH

git clone https://github.com/tylerpitchford/target-project

cd target-project

./build.sh
```

## Docker Compilation [no golang required]
( *multiple architectures - windows/64-bit; linux/64-bit, osx/64-bit* )
```
git clone https://github.com/tylerpitchford/target-project

cd target-project

docker run --rm -it -v "$PWD":/go/src/target-project -w /go/src/target-project golang:1.12.4 ./build.sh 
```

## Docker Interactive Execution [no golang required]

```
git clone https://github.com/tylerpitchford/target-project

cd target-project

docker build -t target .

docker run -i -t target bash
```
