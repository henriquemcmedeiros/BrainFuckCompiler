package main

import (
    "fmt"
    "io"
    "os"
    "strings"
    "strconv"
    "unicode"
)

func main() {
    dados, err := io.ReadAll(os.Stdin)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    texto := strings.TrimSpace(string(dados))
    partes := strings.SplitN(texto, "=", 2)
    if len(partes) != 2 {
        fmt.Fprintln(os.Stderr, "uso: VAR=EXPR")
        os.Exit(1)
    }
    nomeVar, expr := partes[0], partes[1]

    analisador := &Analisador{s: expr}
    arvore := analisador.expr()

    gerador := &GeradorBF{}

    for _, c := range nomeVar + "=" {
        for _, b := range []byte(string(c)) {
            gerador.irPara(10)
            gerador.zerar()
            gerador.incrementar(int(b))
            gerador.sb.WriteByte('.')
        }
    }

    arvore.Gerar(gerador, 0)

    resultado := avaliarNo(arvore)
    for _, ch := range strconv.Itoa(resultado) {
        gerador.irPara(10)
        gerador.zerar()
        gerador.incrementar(int(ch))
        gerador.sb.WriteByte('.')
    }

    fmt.Print(gerador.String())
}

type Analisador struct {
    s   string
    pos int
}

func (a *Analisador) ver() rune {
    if a.pos >= len(a.s) {
        return 0
    }
    return rune(a.s[a.pos])
}

func (a *Analisador) consumir() rune {
    ch := a.ver()
    if ch != 0 {
        a.pos++
    }
    return ch
}

func (a *Analisador) expr() No {
    no := a.termo()
    for {
        switch a.ver() {
        case '+', '-':
            op := byte(a.consumir())
            direito := a.termo()
            no = &BinOp{op: op, esquerdo: no, direito: direito}
        default:
            return no
        }
    }
}

func (a *Analisador) termo() No {
    no := a.fator()
    for a.ver() == '*' {
        a.consumir()
        direito := a.fator()
        no = &BinOp{op: '*', esquerdo: no, direito: direito}
    }
    return no
}

func (a *Analisador) fator() No {
    if a.ver() == '(' {
        a.consumir()
        no := a.expr()
        if a.ver() == ')' {
            a.consumir()
        }
        return no
    }
    inicio := a.pos
    for unicode.IsDigit(a.ver()) {
        a.consumir()
    }
    numStr := a.s[inicio:a.pos]
    num, err := strconv.Atoi(numStr)
    if err != nil {
        fmt.Fprintf(os.Stderr, "número inválido: %s\n", numStr)
        os.Exit(1)
    }
    return &Numero{valor: num}
}

type No interface {
    Gerar(g *GeradorBF, celula int)
}

type Numero struct{ valor int }

func (n *Numero) Gerar(g *GeradorBF, celula int) {
    g.irPara(celula)
    g.zerar()
    g.incrementar(n.valor)
}

type BinOp struct {
    op                byte
    esquerdo, direito No
}

func (b *BinOp) Gerar(g *GeradorBF, celula int) {
    switch b.op {
    case '+':
        b.esquerdo.Gerar(g, celula)
        b.direito.Gerar(g, celula+1)
        g.somar(celula+1, celula)
    case '-':
        b.esquerdo.Gerar(g, celula)
        b.direito.Gerar(g, celula+1)
        g.subtrair(celula+1, celula)
    case '*':
        b.esquerdo.Gerar(g, celula)
        b.direito.Gerar(g, celula+1)
        g.multiplicar(celula, celula+1, celula+2, celula+3)
    }
}

func avaliarNo(n No) int {
    switch v := n.(type) {
    case *Numero:
        return v.valor
    case *BinOp:
        E := avaliarNo(v.esquerdo)
        D := avaliarNo(v.direito)
        switch v.op {
        case '+':
            return E + D
        case '-':
            return E - D
        case '*':
            return E * D
        }
    }
    return 0
}

type GeradorBF struct {
    sb  strings.Builder
    pos int
}

func (g *GeradorBF) irPara(c int) {
    for g.pos < c {
        g.sb.WriteByte('>')
        g.pos++
    }
    for g.pos > c {
        g.sb.WriteByte('<')
        g.pos--
    }
}

func (g *GeradorBF) zerar() {
    g.sb.WriteString("[-]")
}

func (g *GeradorBF) incrementar(n int) {
    for i := 0; i < n; i++ {
        g.sb.WriteByte('+')
    }
}

func (g *GeradorBF) loop(c int, corpo func()) {
    g.irPara(c)
    g.sb.WriteByte('[')
    corpo()
    g.irPara(c)
    g.sb.WriteByte(']')
}

func (g *GeradorBF) somar(origem, destino int) {
    g.loop(origem, func() {
        g.sb.WriteByte('-')
        g.irPara(destino)
        g.sb.WriteByte('+')
        g.irPara(origem)
    })
    g.irPara(destino)
}

func (g *GeradorBF) subtrair(origem, destino int) {
    g.loop(origem, func() {
        g.sb.WriteByte('-')
        g.irPara(destino)
        g.sb.WriteByte('-')
        g.irPara(origem)
    })
    g.irPara(destino)
}

func (g *GeradorBF) multiplicar(a, b, res, tmp int) {
    g.irPara(res)
    g.zerar()
    g.irPara(tmp)
    g.zerar()
    g.loop(a, func() {
        g.irPara(a)
        g.sb.WriteByte('-')
        g.loop(b, func() {
            g.sb.WriteByte('-')
            g.irPara(res)
            g.sb.WriteByte('+')
            g.irPara(tmp)
            g.sb.WriteByte('+')
            g.irPara(b)
        })
        g.loop(tmp, func() {
            g.sb.WriteByte('-')
            g.irPara(b)
            g.sb.WriteByte('+')
            g.irPara(tmp)
        })
    })
    g.loop(res, func() {
        g.sb.WriteByte('-')
        g.irPara(a)
        g.sb.WriteByte('+')
        g.irPara(res)
    })
    g.irPara(a)
}

func (g *GeradorBF) String() string {
    return g.sb.String()
}
