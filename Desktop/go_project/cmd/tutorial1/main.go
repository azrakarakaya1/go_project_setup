/*
	1. Statically Typed Language
	2. Strongly Typed Language
	3. GO is Compiled
	4. Fast Compilation Time
	5. Built In Concurrency
	6. Simplicity

	Package	= Collection of go files
	Module	= Collection of packages

	go mod init github.com/azrakarakaya1/go.git 
*/

package main		//look for an entry point function
import "fmt"

func main() {

	var num int8 = -10		//int8 16 32 64 to specify memory/bits
	var unum uint8 = 10
	fmt.Println(num, "\t", unum)

	var fnum float64 = 12345678.90000	//float32 64
	fmt.Println(fnum)

	var str string = `abc
	def`
	fmt.Println(len(str))
}