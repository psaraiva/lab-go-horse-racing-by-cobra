[![project](https://img.shields.io/badge/github-psaraiva%2Flab--go--horse--racing--by--cobra-blue)](https://img.shields.io/badge/github-psaraiva%2Flab--go--horse--racing--by--cobra-blue)
[![License](https://img.shields.io/badge/license-MIT-%233DA639.svg)](https://opensource.org/licenses/MIT)

[![Go Report Card](https://goreportcard.com/badge/github.com/psaraiva/lab-go-horse-racing-by-cobra)](https://goreportcard.com/report/github.com/psaraiva/lab-go-horse-racing-by-cobra)
![Codecov](https://img.shields.io/codecov/c/github/psaraiva/lab-go-horse-racing-by-cobra)

[![Idioma: Português](https://img.shields.io/badge/Idioma-Português-green?style=flat-square)](README.pt-br.md)

# 🐎 Lab Go: Horse Racing by Cobra 🐍

## 🎯 Objective
This lab aims to demonstrate the use of Goroutines in a simple, practical, and fun way using Cobra.

## ⚙️ How it works?
The horses run until the first one crosses the finish line.

## 💻 Commands
Example of use
```bash
git clone https://github.com/psaraiva/lab-go-horse-racing-by-cobra.git
cd lab-go-horse-racing-by-cobra
docker build -t lab-go-horse-racing-by-cobra .
```

Run with default settings
```bash
docker run --rm -it lab-go-horse-racing-by-cobra
```

Run with 5 horses and a score target of 50
```bash
docker run --rm -it lab-go-horse-racing-by-cobra --horses-quantity 5 --score-target 50
```

Run with label 'C' and a timeout of 15 seconds
```bash
docker run --rm -it lab-go-horse-racing-by-cobra --horse-label C --game-timeout 15s
```

Run with 20 horses, label 'P', score target 50, and a timeout of 90 seconds
```bash
docker run --rm -it lab-go-horse-racing-by-cobra --horses-quantity 20 --horse-label P --score-target 50 --game-timeout 90s
```

## 🔧 Parameters
- `--horse-label`
  - default value: `H`
  - valid value: `char(1)`
- `--horses-quantity`
  - default value: `2`
  - valid value: `int 2..99`
- `--score-target`
  - default value: `75`
  - valid value: `int 15..100`
- `--game-timeout`
  - default value: `10s`
  - valid value: `string 10s..90s`

## Timeout message
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
