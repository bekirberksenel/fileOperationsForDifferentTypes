package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"os"
	"log"
	"encoding/binary"
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

	processFile(sourceFile, *numPtr)
}

func processFile(srcFile string, num int) {
	f, err := os.Open(srcFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	gzf, err := gzip.NewReader(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tarReader := tar.NewReader(gzf)

	i := 1
	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		name := header.Name
		
		/*type Integers struct { //
			I1 uint16		   //
			I2 int32		   //
			I3 int64		   //
		}*/					   //
		

		switch header.Typeflag {
		case tar.TypeDir:
			continue
		case tar.TypeReg:
			fmt.Println("(", i, ")", "Name: ", name)
			
			/*f, err := os.Open(name) //
			if err != nil {
				log.Fatalln(err)
			}
			defer f.Close()
	
			q := Integers{}
			err = binary.Read(f, binary.LittleEndian, &q)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println(q) //			
			*/
			fl, err := os.Open(name)
			if err != nil {
				log.Fatalln(err)
			}
			defer fl.Close()

			for {
				var val float64
				err = binary.Read(fl, binary.LittleEndian, &val)
				if err != nil {
					if err == io.EOF {
						break
					}
					log.Fatalln(err)
				}
				fmt.Println(val)
			}
			
			
			
			
			
			
			
			
			if i == num {
				fmt.Println(" --- ")
				io.Copy(os.Stdout, tarReader)
				fmt.Println(" --- ")
				os.Exit(0)
			}
		default:
			fmt.Printf("%s : %c %s %s\n",
				"Unable to figure out type",
				header.Typeflag,
				"in file",
				name,
			)
		}

		i++
	}
}