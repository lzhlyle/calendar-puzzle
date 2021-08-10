# calendar-puzzle

calendar puzzle, 50 cells, 10 blocks

## Board Style

```
Jan Feb Mar Apr May Jun
Jul Aug Sep Oct Nov Dec
1   2   3   4   5   6   7
8   9   10  11  12  13  14
15  16  17  18  19  20  21
22  23  24  25  26  27  28
29  30  31  Sun Mon Tue Wed
                Thr Fri Sat
```

## How to generate puzzle solution

```go
go run main.go // generate puzzle solution for today

go run main.go 2021 08 10 // generate puzzle solution for 2021-08-10
``` 
