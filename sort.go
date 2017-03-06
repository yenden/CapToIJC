package main
import (
	"fmt"
	"os"
	"bufio"
	"sort"
	"strings"
)

/*remove duplicates from a string slice*/
func removeDuplicates(elements []string) []string {
	encountered := map[string]bool{}	//mapping each string with
						//a bool in target duplicates
	finalResult := []string{}
	for v:= range elements {
		if encountered[elements[v]] == true {
			//nothing to be done
		} else {
			encountered[elements[v]] = true
			finalResult = append(finalResult, elements[v])
		}
	}
	return finalResult
}

/* read lines from a slice of files, and put all the files lines in "lines"*/
func readLines(files []string) (lines []string, err error) {
	var line string
	for _,file := range files {
		f, err:= os.Open(file)
		if err!= nil {
			fmt.Println(err)
			return nil, err
		}
		defer f.Close()
		r:= bufio.NewScanner(f)
		for r.Scan() {
			line= r.Text()
			for _, fl := range line {
				//if a line strat with a number
				//between 0 and 9
				if fl<=57 && fl>=48 {
					result:= strings.Split(line, " ")
					result[1] += string("\r\n")
					lines = append(lines, result[1])
				}
				break
			}
		}
	}
	lines = removeDuplicates(lines)
	return lines, err
}

/* write lines in a file */
func writeLines(file string, lines[]string) (err error ){
	f, err := os.Create(file)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()
	w:= bufio.NewWriter(f)
	defer w.Flush()
	for _, line := range lines {
		_, err := w.WriteString(line)
		if err != nil {
			return err
		}

	}
	return nil
}

//main function
func main() {
	tabFiles := []string{
		`EMVConstants.txt`,
		`EMVCrypto.txt`,
		`EMVProtocolState.txt`,
		`EMVStaticData.txt`,
		`SimpleEMVApplet.txt`,
		`Passeport.txt`,

	}
	/*tabFiles := []string{
		`Passeport.txt`,
	}*/
	lines, err := readLines(tabFiles)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sort.Strings(lines)
	file2 :=`result2.txt`
	err = writeLines(file2, lines)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}


