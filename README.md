## Brainfuck Compiler e Executer

Este repositório traz dois utilitários em Go para avaliar expressões e imprimir resultados usando Brainfuck:

1. **`bfc`**: recebe uma entrada `VAR=EXPR`, gera um programa Brainfuck que calcula `EXPR` em tempo de execução.
2. **`bfe`**: executa o código Brainfuck gerado e imprime o resultado como `VAR=valor`.

---

### Funcionamento

* **Parser**: `bfc` faz parsing LL(1) de expressões com `+`, `-`, `*`, parênteses e números.

* **Execução**: `bfe` interpreta Brainfuck, gerencia a fita e imprime o resultado.

---

### Compilação

```bash
go build -o bfc bfc.go
go build -o bfc bfc.go
```
---

### Exemplo

```bash
echo 'VAR=2+5*10' | ./bfc
# Saída: código Brainfuck correspondente

echo 'VAR=2+5*10' | ./bfc | ./bfe
# Saída: VAR=20
```
