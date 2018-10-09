package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"encoding/binary"
	"log"
	"flag"
)

/*var files = []struct {
	Name, Body string
}{
	{"readme.txt", "This archive contains some text files."},
	{"gopher.txt", "Gopher names:\nGeorge\nGeoffrey\nGonzo"},
	{"todo.txt", "Get animal handling licence."},
}*/

func main() {
	/*mpath := "ets_datas_all.tar.gz"
	//tarIt(mpath)
	untarIt(mpath)
	numberOfFiles("ets_datas_all.tar.gz")
    */
	arg := os.Args[1]
	untarIt(arg)
	numberOfFiles(arg)
}



/*func tarIt(mpath string) {
	f, err := overwrite(mpath)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	gw := gzip.NewWriter(f)
	defer gw.Close()
	if err != nil {
		panic(err)
	}
	tw := tar.NewWriter(gw)
	for _, file := range files {
		hdr := &tar.Header{
			Name:     file.Name,
			Mode:     0600,
			Size:     int64(len(file.Body)),
			Typeflag: byte('0'),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			panic(err)
		}
		if _, err := tw.Write([]byte(file.Body)); err != nil {
			panic(err)
		}
	}
	// Make sure to check the error on Close.
	if err := tw.Close(); err != nil {
		panic(err)
	}
}*/

func untarIt(mpath string) {
	fr, err := read(mpath)
	defer fr.Close()
	if err != nil {
		panic(err)
	}
	gr, err := gzip.NewReader(fr)
	defer gr.Close()
	if err != nil {
		panic(err)
	}
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
			panic(err)
		}
		path := hdr.Name
		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, os.FileMode(hdr.Mode)); err != nil {
				panic(err)
			}
		case tar.TypeReg:
			ow, err := overwrite(path)
			defer ow.Close()
			if err != nil {
				panic(err)
			}
			if _, err := io.Copy(ow, tr); err != nil {
				panic(err)
			}
		default:
			fmt.Printf("Can't: %c, %s\n", hdr.Typeflag, path)
		}
	}
}

func overwrite(mpath string) (*os.File, error) {
	f, err := os.OpenFile(mpath, os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		f, err = os.Create(mpath)
		if err != nil {
			return f, err
		}
	}
	return f, nil
}

func read(mpath string) (*os.File, error) {
	f, err := os.OpenFile(mpath, os.O_RDONLY, 0444)
	if err != nil {
		return f, err
	}
	return f, nil
}

func readAndPrintFile(mpath string) {
	f, err := os.Open(mpath)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	for {
		var val float64
		err = binary.Read(f, binary.LittleEndian, &val)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalln(err)
		}
		fmt.Println(val)
	}
}

func numberOfFiles (sourceFile string) {

	numPtr := flag.Int("n", 4, "an integer")
	flag.Parse()

	if sourceFile == "" {
		fmt.Println("You didn't pass in a tar file!")
		os.Exit(1)
	}


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
			
			
			readAndPrintFile(name)
			
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