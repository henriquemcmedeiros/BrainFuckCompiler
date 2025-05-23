package main

import (
    "fmt"
    "io"
    "os"
)

func main() {
    entrada, err := io.ReadAll(os.Stdin)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    codigo := string(entrada)

    saltos := map[int]int{}
    var pilha []int
    for indice, caractere := range codigo {
        switch caractere {
        case '[':
            pilha = append(pilha, indice)
        case ']':
            if len(pilha) == 0 {
                fmt.Fprintf(os.Stderr, "']' sem par em %d\n", indice)
                os.Exit(1)
            }
            j := pilha[len(pilha)-1]
            pilha = pilha[:len(pilha)-1]
            saltos[j] = indice
            saltos[indice] = j
        }
    }
    if len(pilha) != 0 {
        fmt.Fprintln(os.Stderr, "'[' sem par")
        os.Exit(1)
    }

    fita := make([]byte, 30000)
    ponteiro := 0

    for instrucao := 0; instrucao < len(codigo); instrucao++ {
        switch codigo[instrucao] {
        case '>':
            ponteiro++
            if ponteiro >= len(fita) {
                ponteiro = 0
            }
        case '<':
            ponteiro--
            if ponteiro < 0 {
                ponteiro = len(fita) - 1
            }
        case '+':
            fita[ponteiro]++
        case '-':
            fita[ponteiro]--
        case '.':
            os.Stdout.Write([]byte{fita[ponteiro]})
        case '[':
            if fita[ponteiro] == 0 {
                instrucao = saltos[instrucao]
            }
        case ']':
            if fita[ponteiro] != 0 {
                instrucao = saltos[instrucao]
            }
        }
    }

    os.Stdout.Write([]byte{'\n'})
}
