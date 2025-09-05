[![project](https://img.shields.io/badge/github-psaraiva%2Flab--go--horse--racing--by--cobra-blue)](https://img.shields.io/badge/github-psaraiva%2Flab--go--horse--racing--by--cobra-blue)
[![License](https://img.shields.io/badge/license-MIT-%233DA639.svg)](https://opensource.org/licenses/MIT)

[![Go Report Card](https://goreportcard.com/badge/github.com/psaraiva/lab-go-horse-racing-by-cobra)](https://goreportcard.com/report/github.com/psaraiva/lab-go-horse-racing-by-cobra)
![Codecov](https://img.shields.io/codecov/c/github/psaraiva/lab-go-horse-racing-by-cobra)

[![Idioma: English](https://img.shields.io/badge/Idioma-English-blue?style=flat-square)](README.md)

# 🐎 Lab Go: Corrida de Cavalos por Cobra 🐍

## 🎯 Objetivo
Este laboratório demonstra o uso de Goroutines de uma forma simples, prática e divertida, utilizando a biblioteca Cobra.

## ⚙️ Como isso funciona?
Os cavalos correm até o primeiro cruzar a linha de chegada.

## 💻 Comandos
Exemplo de uso
```bash
git clone https://github.com/psaraiva/lab-go-horse-racing-by-cobra.git
cp lab-go-horse-racing-by-cobra
docker build -t lab-go-horse-racing-by-cobra .
```

Executa com configuração padrão
```bash
docker run --rm -it lab-go-horse-racing-by-cobra
```

Executa com 5 cavalos e alvo de 50 pontos
```bash
docker run --rm -it lab-go-horse-racing-by-cobra --horses-quantity 5 --score-target 50
```

Executa com o label 'C' e um timeout de 15 segundos
```bash
docker run --rm -it lab-go-horse-racing-by-cobra --horse-label C --game-timeout 15s
```

Executa com o 20 cavalos, label 'P', alvo 75 ponsto e um timeout de 90 segundos
```bash
docker run --rm -it lab-go-horse-racing-by-cobra --horses-quantity 20 --horse-label P --score-target 50 --game-timeout 90s
```

## 🔧 Parâmetros
- `--horse-label`
  - valor padrão `H`
  - valor válido `char(1)`
- `--horses-quantity`
  - valor padrão `2`
  - valor válido `int 2..99`
- `--score-target`
  - valor padrão `75`
  - valor válido `int 15..100`
- `--game-timeout`
  - valor padrão `10s`
  - valor válido `string 10s..90s`

## Exemplo
```bash
   +---------|---------|---------|---------|---------|---------|---------|---------|--+
H01|................................................................................H01|
H02|........................................................................H02       |
H03|..............................................................................H03 |
H04|............................................................................H04   |
H05|...............................................................................H05|
H06|..............................................................................H06 |
H07|.............................................................................H07  |
H08|..............................................................................H08 |
H09|.........................................................................H09      |
H10|.........................................................................H10      |
   +---------|---------|---------|---------|---------|---------|---------|---------|--+
```

Mensagem de tempo esgotado
```bash
   +---------|---------|---------|---------|---------|---------|---------|---------|--+
H01|..................H01                                                             |
H02|...............H02                                                                |
H03|.....................H03                                                          |
   +---------|---------|---------|---------|---------|---------|---------|---------|--+

Today is a very hot day, the horses are tired!
```
## Preview
![Preview](./asset/horse_race.gif)
