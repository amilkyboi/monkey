// repl/repl.go

package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/lexer"
	"monkey/token"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	// Starts the REPL

	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, PROMPT)

		// Read from the input until encountering a newline
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		// Pass the read line into an instance of the lexer
		line := scanner.Text()
		l := lexer.New(line)

		// Print the tokens output by the lexer until encountering an EOF
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}
