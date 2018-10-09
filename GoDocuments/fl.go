package main

import (
    "archive/tar"
    "fmt"
    "io"
    "log"
    "os"
    "bytes"
    "compress/gzip"
	"flag"
)

func main() {

	numPtr := flag.Int("n", 4, "an integer")
	flag.Parse()

	sourceFile := flag.Arg(0)

	if sourceFile == "" {
		fmt.Println("You didn't pass in a tar file!")
		os.Exit(1)
	}

	fmt.Println("arg 1: ", flag.Arg(0))

	processFile(sourceFile, *numPtr) //


    archive, err := gzip.NewReader(sourceFile)

    if err != nil {
        fmt.Println("There is a problem with os.Open")
    }
    tr := tar.NewReader(archive)

    for {
        hdr, err := tr.Next()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Fatal(err)
        }

        fmt.Printf("Contents of %s:\n", hdr.Name)

        //Using a bytes buffer is an important part to print the values as a string

        bud := new(bytes.Buffer)
        bud.ReadFrom(tr)
        s := bud.String()
        fmt.Println(s)
        fmt.Println()
    }

}