package main

import (
	"fmt"
	"os"
	"bufio"
	"log"
	"regexp"
	"strconv"
)

// variables
var consumers int
var lines = []string{}
var totalCount int

// producer function
func produce(jobs chan<- string) {
	for _, line := range lines {
		jobs <- line
	}
	close(jobs)
}

// consumer function
func consume(worker int, jobs <-chan string, done chan<- bool) {
	for line := range jobs {
		// count num words and print
                var cnt int = Count(line)
		totalCount += cnt
		fmt.Printf("Line: `%v` | processed by worker: #%v | Word Count: %v\n", line, worker, cnt)
	}
	done <- true
}

// counts number of words
func Count(line string) int {
    // match non-space character sequences
    re := regexp.MustCompile(`[\S]+`)

    // find all matches and return count
    results := re.FindAllString(line, -1)
    return len(results)
}

// main function
func main() {

    // num consumers from command line arguments
    if len(os.Args) > 1 {
	    // convert arg to int
	    var err error
	    nums := make([]int, len(os.Args))
	    if nums[1], err = strconv.Atoi(os.Args[1]); err != nil {
		    panic(err)
	    }
	    consumers = nums[1]
	    fmt.Print("Number of consumers: ", consumers)
    } else {
	    consumers = 3
	    fmt.Print("No number of consumers specified!\nDefault Consumer Number: 3")
    }
    // create channels
    jobs := make(chan string)
    done := make(chan bool)
    
    // instructions for users
    fmt.Print("\nEnter file name or string:\n\n")

    // get file name from input
    var input string
    scanner := bufio.NewScanner(os.Stdin)
    scanner.Scan()

	    // save input
	    input = scanner.Text()

	    // open file
	    file, err := os.Open("./" + input)

	    // if file not found
	    if err != nil {
		    // add line to array
                    lines = append(lines, input)

	    } else {
		    // reset line array
		    lines = nil

		    // read file
		    scanner2 := bufio.NewScanner(file)
		    for scanner2.Scan()  {
			    // add lines to array
			    lines = append(lines, scanner2.Text())
		    }

		    // break if error
		    if err2 := scanner2.Err(); err2 != nil {
			    log.Fatal(err2)
		    }
	    }
	    // produce jobs
	    go produce(jobs)

	    // consume jobs
	    for i := 1; i <= consumers; i++ {
	    go consume(i, jobs, done)
	    }
	    <-done

    // break if error
    if errr := scanner.Err(); errr != nil { // if error
	    log.Println(errr)
    }
	
    // print total word count
    fmt.Println("\n---------------------------------------------------------------------\nThe total word count is:", totalCount)
    
    // completion confirmation
    fmt.Println("\n--------------------------------------------------------------------\n>> Program exectuion completed <<\n>> To run again, type: `go run project2.go {number of consumers}` <<\n--------------------------------------------------------------------\n\n")

}

