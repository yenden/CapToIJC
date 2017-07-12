package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const (
	tagHeaderComp            = 0x01
	tagDirComp               = 0x02
	tagAppletComp            = 0x03
	tagImportComp            = 0x04
	tagConstantPoolComp      = 0x05
	tagClassComp             = 0x06
	tagMethodComp            = 0x07
	tagStaticFieldComp       = 0x08
	tagReferenceLocationComp = 0x09
	tagExportComp            = 0x0A
	tagDescriptorComp        = 0x0B
	tagDebugComp             = 0x0C
)

//capture the console output
func captureStdout(baComp []byte, comp int) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	writeToFile(baComp, comp)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}
func writeToFile(baComp []byte, comp int) {
	switch comp {
	case tagHeaderComp:
		fmt.Println("**************Header component*************")
	case tagDirComp:
		fmt.Println("**************Directory component*************")
	case tagImportComp:
		fmt.Println("**************Import component*************")
	case tagAppletComp:
		fmt.Println("**************Applet component*************")
	case tagClassComp:
		fmt.Println("**************Class component*************")
	case tagMethodComp:
		fmt.Println("**************Method component*************")
	case tagStaticFieldComp:
		fmt.Println("**************Static Field component*************")
	case tagExportComp:
		fmt.Println("**************Export component*************")
	case tagConstantPoolComp:
		fmt.Println("**************Constant Pool component*************")
	case tagReferenceLocationComp:
		fmt.Println("**************Reference Location component*************")
	case tagDescriptorComp:
		fmt.Println("**************Descriptor component*************")
	case tagDebugComp:
		fmt.Println("**************Debug component*************")
	}
	fmt.Println(baComp)
	fmt.Println()
}
func readCAP(sCapFileName string, option int) {

	var (
		baHeader       []byte
		baDirectory    []byte
		baApplet       []byte
		baImport       []byte
		baConstantPool []byte
		baClass        []byte
		baMethod       []byte
		baStaticField  []byte
		baRefLocation  []byte
		baExport       []byte
		baDescriptor   []byte
		baDebug        []byte
	)

	r, err := zip.OpenReader(sCapFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	for _, f := range r.File {
		if !f.FileInfo().IsDir() {
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
	 * option =O => debug component will be not included.
	 * option =1 => all components will be included
	 */

	outComponent(w, baHeader, tagHeaderComp)
	outComponent(w, baDirectory, tagDirComp)
	outComponent(w, baImport, tagImportComp)
	outComponent(w, baApplet, tagAppletComp)
	outComponent(w, baClass, tagClassComp)
	outComponent(w, baMethod, tagMethodComp)
	outComponent(w, baStaticField, tagStaticFieldComp)
	outComponent(w, baExport, tagExportComp)
	outComponent(w, baConstantPool, tagConstantPoolComp)
	outComponent(w, baRefLocation, tagReferenceLocationComp)
	outComponent(w, baDescriptor, tagDescriptorComp)

	//if option ==1 debug compwill be included
	if option == 1 {
		outComponent(w, baDebug, tagDebugComp)
	}
}

//write in the .ijc file
func outComponent(wout *bufio.Writer, baComp []byte, comp int) {
	if baComp != nil {
		toWrite := captureStdout(baComp, comp)
		_, err := wout.Write([]byte(toWrite))
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
		fmt.Println("Options: \n\t -a to include debug components")
	} else if len(args) == 1 {
		readCAP(args[0], 0)
	} else {
		readCAP(args[1], 1)
	}
	fmt.Println("\tConversion done")
}
