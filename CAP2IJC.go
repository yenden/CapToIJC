package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func readCAP(sCapFileName string, option int) {

	var baHeader []byte
	var baDirectory []byte
	var baApplet []byte
	var baImport []byte
	var baConstantPool []byte
	var baClass []byte
	var baMethod []byte
	var baStaticField []byte
	var baRefLocation []byte
	var baExport []byte
	var baDescriptor []byte
	var baDebug []byte

	r, err := zip.OpenReader(sCapFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			//nothing to be done
		} else {
			sFileName := f.Name
			sDirs := strings.Split(sFileName, "/")
			switch sDirs[len(sDirs)-1] {
			case "Header.cap":
				fmt.Println("Processing Header.cap")
				baHeader = setComponent(f)
			case "Directory.cap":
				fmt.Println("Processing Directory.cap")
				baDirectory = setComponent(f)
			case "Applet.cap":
				fmt.Println("Processing Applet.cap")
				baApplet = setComponent(f)
			case "Import.cap":
				fmt.Println("Processing Import.cap")
				baImport = setComponent(f)
			case "ConstantPool.cap":
				fmt.Println("Processing ConstantPool.cap")
				baConstantPool = setComponent(f)
			case "Class.cap":
				fmt.Println("Processing Class.cap")
				baClass = setComponent(f)
			case "Method.cap":
				fmt.Println("Processing Method.cap")
				baMethod = setComponent(f)
			case "StaticField.cap":
				fmt.Println("Processing StaticField.cap")
				baStaticField = setComponent(f)
			case "RefLocation.cap":
				fmt.Println("Processing RefLocation.cap")
				baRefLocation = setComponent(f)
			case "Export.cap":
				fmt.Println("Processing Export.cap")
				baExport = setComponent(f)
			case "Descriptor.cap":
				fmt.Println("Processing Descriptor.cap")
				baDescriptor = setComponent(f)
			case "Debug.cap":
				fmt.Println("Processing Debug.cap")
				baDebug = setComponent(f)
			default:
				//nothing

			}
		}
	}
	//recuparate the name of the cap file
	sDirsOut := strings.Split(sCapFileName, "/")
	sCapFile := sDirsOut[len(sDirsOut)-1]
	n := len(sCapFile) - 4
	sIJCFileName := sCapFile[:n]
	sIJCFileName = sIJCFileName + ".ijc"

	fout, err := os.Create(sIJCFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer fout.Close()
	w := bufio.NewWriter(fout)
	defer w.Flush()

	/*option is related to the installation order
	 * describe in JCVM specification
	 * option =O => descriptor component will be not included.
	 * option =1 => all components will be included
	 */

	if option == 0 {
		outComponent(w, baHeader)
		outComponent(w, baDirectory)
		outComponent(w, baImport)
		outComponent(w, baApplet)
		outComponent(w, baClass)
		outComponent(w, baMethod)
		outComponent(w, baStaticField)
		outComponent(w, baExport)
		outComponent(w, baConstantPool)
		outComponent(w, baRefLocation)

	} else if option == 1 {
		outComponent(w, baHeader)
		outComponent(w, baDirectory)
		outComponent(w, baImport)
		outComponent(w, baApplet)
		outComponent(w, baClass)
		outComponent(w, baMethod)
		outComponent(w, baStaticField)
		outComponent(w, baExport)
		outComponent(w, baConstantPool)
		outComponent(w, baRefLocation)
		outComponent(w, baDescriptor)
		outComponent(w, baDebug)
	}

}

//write in the .ijc file
func outComponent(wout *bufio.Writer, baComp []byte) {
	if baComp != nil {
		_, err := wout.Write(baComp)
		if err != nil {
			log.Fatal(err)
		}
	}
}

//copy component length bytes from cap file
func setComponent(file *zip.File) []byte {
	var baComp = make([]byte, file.FileHeader.FileInfo().Size())
	r, err := file.Open()
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadFull(r, baComp)

	if err != nil {
		log.Fatal(err)
		fmt.Println(b)
	}
	return baComp

}

func main() {
	fmt.Println("----CAP2IJC version 1.0 ----")
	args := os.Args[1:]
	if args == nil || len(args) == 0 {
		fmt.Println("Usage: \n\tCap2IJC filename.cap or\n\tCap2IJC -a filename.cap")
		fmt.Println("Options: \n\t -a to include descriptor and debug components")
	} else if len(args) == 1 {
		readCAP(args[0], 0)
	} else {
		readCAP(args[1], 1)
	}
	fmt.Println("\tConversion done")
}
