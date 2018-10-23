package main

import (
    "log"
    "fmt"
    "os"
    "strconv"
    "strings"
    "bufio"
    "errors"
    "sync"
)

var (
    err error
)

func main() {
    if len(os.Args) != 3 {
        log.Fatal("Invalid number of arguments.")
    }

    numberOfThreads, err := strconv.Atoi(os.Args[1])
    if err != nil {
        log.Fatal(err)
    }
    if numberOfThreads < 1 {
        numberOfThreads = 1
    }

    file, err := os.Open(os.Args[2])
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    matrixA, matrixB := createMatrices(file)

    fmt.Printf("Matrix A:\n")
    printMatrix(matrixA)

    fmt.Printf("\nMatrix B:\n")
    printMatrix(matrixB)


    fmt.Printf("\nNumber of threads: %d\n\n", numberOfThreads)

    m, err := multiply(matrixA, matrixB, numberOfThreads)
    if err != nil {
        fmt.Println(err)
    }

    fmt.Printf("\nResulting matrix:\n")

    printMatrix(m)

}

func createMatrices(file *os.File) ([][]float32, [][]float32) {
    var matrixA []string
    var matrixB []string

    scanner := bufio.NewScanner(file)
    numEmptyLines := 0
    for scanner.Scan() {
        line := scanner.Text()
        if len(line) == 0 {
            numEmptyLines++
            continue
        }
        if numEmptyLines == 2 {
            break
        } else if numEmptyLines == 0 {
            matrixA = append(matrixA, line)
        } else {
            matrixB = append(matrixB, line)
        }
    }

    var mA = make([][]float32, len(matrixA))
    for row := range matrixA {
        line := strings.Split(matrixA[row], ",")
        mA[row] = make([]float32, len(line))
        for col := range line {
            elem, err := strconv.ParseFloat(strings.TrimSpace(line[col]), 32)
            if err != nil {
                log.Fatal(err)
            }
            mA[row][col] = float32(elem)
        }
    }

    var mB = make([][]float32, len(matrixB))
    for row := range matrixB {
        line := strings.Split(matrixB[row], ",")
        mB[row] = make([]float32, len(line))
        for col := range line {
            elem, err := strconv.ParseFloat(strings.TrimSpace(line[col]), 32)
            if err != nil {
                log.Fatal(err)
            }
            mB[row][col] = float32(elem)
        }
    }

    return mA, mB
}

func multiply(a [][]float32, b [][]float32, threads int) ([][]float32, error) {
    // if a cols != b rows
    if len(a[0]) != len(b) {
        return nil, errors.New("First matrix's column size must match second matrix's row size.")
    }
    if len(a) > 100 || len(a[0]) > 100 || len(b) > 100 || len(b[0]) > 100 {
        return nil, errors.New("Matrices are too big to compute.")
    }

    tasks := len(a) * len(b[0])
    fmt.Println("Tasks:", tasks)

    var taskDistribution []int

    if float32(tasks) / float32(threads) <= 1.0 {
        taskDistribution = make([]int, tasks)
        for i := 0; i < tasks; i++ {
            taskDistribution[i] = 1
        }
    } else {
        r := tasks % threads
        n := tasks / threads
        taskDistribution = make([]int, threads)
        for i := 0; i < threads; i++ {
            taskDistribution[i] = n
        }
        for i := 0; i < r; i++ {
            taskDistribution[i] += 1
        }
    }

    fmt.Printf("\nTask Distribution: %v\n\n", taskDistribution)

    prod := make([][]float32, len(a))
    for i := range prod {
        prod[i] = make([]float32, len(b[0]))
    }

    var wg sync.WaitGroup
    wg.Add(len(taskDistribution))

    pos := 0
    for i := 0; i < len(taskDistribution); i++ {
        go func (id, start, numJobs int) {
            for j := 0; j < numJobs; j++ {
                var row, col int
                row = (start + j) / len(b[0])
                col = (start + j) % len(b[0])
                fmt.Printf("Worker %d starting task %d. Multiplying output[%d][%d]\n", id, start + j + 1, row, col)

                var dotProduct float32 = 0.0
                for k := range a[0] {
                    dotProduct += a[row][k] * b[k][col]
                }

                prod[row][col] = dotProduct

                fmt.Printf("Worker %d finished task %d.\n", id, start + j + 1)
            }
            wg.Done()
        }(i + 1, pos, taskDistribution[i])
        pos += taskDistribution[i]
    }

    wg.Wait()

    return prod, nil
}

func printMatrix(m [][]float32) {
    for row := range m {
        for col := range m[row] {
            fmt.Print(m[row][col])
            if col <= len(m) {
                fmt.Print(", ")
            }
        }
        fmt.Println()
    }
}
