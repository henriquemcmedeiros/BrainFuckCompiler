package main

import (
    "fmt"
    "io"
    "os"
    "strings"
    "strconv"
)

func main() {
    data, err := io.ReadAll(os.Stdin)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    text := strings.TrimSpace(string(data))
    parts := strings.SplitN(text, "=", 2)
    if len(parts) != 2 {
        fmt.Fprintln(os.Stderr, "Entrada inválida, use VAR=EXPR")
        os.Exit(1)
    }
    varName := strings.TrimSpace(parts[0])
    expr := parts[1]

    parser := &Parser{s: expr}
    result := parser.parseExpr()

    output := varName + "=" + strconv.Itoa(result)

    // Gera código Brainfuck que, célula a célula,
    // ajusta o valor até o byte desejado, emite ('.')
    // e zera ([-]) para o próximo caractere.
    var bf strings.Builder
    for i := 0; i < len(output); i++ {
        b := output[i]
        for j := 0; j < int(b); j++ {
            bf.WriteByte('+')
        }
        bf.WriteByte('.')
        bf.WriteString("[-]")
    }
    fmt.Print(bf.String())
}

// ----- parser recursivo simples para + - * / e ( ) -----

type Parser struct {
    s   string
    pos int
}

func (p *Parser) peek() byte {
    if p.pos >= len(p.s) {
        return 0
    }
    return p.s[p.pos]
}

func (p *Parser) consume() byte {
    ch := p.peek()
    if ch != 0 {
        p.pos++
    }
    return ch
}

func (p *Parser) parseExpr() int {
    val := p.parseTerm()
    for {
        switch p.peek() {
        case '+':
            p.consume()
            val += p.parseTerm()
        case '-':
            p.consume()
            val -= p.parseTerm()
        default:
            return val
        }
    }
}

func (p *Parser) parseTerm() int {
    val := p.parseFactor()
    for {
        switch p.peek() {
        case '*':
            p.consume()
            val *= p.parseFactor()
        case '/':
            p.consume()
            val /= p.parseFactor()
        default:
            return val
        }
    }
}

func (p *Parser) parseFactor() int {
    if p.peek() == '(' {
        p.consume()
        val := p.parseExpr()
        if p.peek() == ')' {
            p.consume()
        }
        return val
    }
    start := p.pos
    for c := p.peek(); c >= '0' && c <= '9'; c = p.peek() {
        p.consume()
    }
    numStr := p.s[start:p.pos]
    num, err := strconv.Atoi(numStr)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Número inválido: %s\n", numStr)
        os.Exit(1)
    }
    return num
}
