package main

import (
    "fmt"
    "log"
    "os"
    "errors"
    "sync"
    "math/rand"
    "time"
    "bufio"
)

var (
    err error
)

func main() {

    perform()

}

func perform() {
    seqFile, err := os.Create("Seq_exe.csv")
    if err != nil {
        log.Fatal("Could not read Seq_exe.csv.")
    }
    defer seqFile.Close()

    paralFile, err := os.Create("Paral_exe.csv")
    if err != nil {
        log.Fatal("Could not read Seq_exe.csv.")
    }
    defer paralFile.Close()

    // sequential
    w := bufio.NewWriter(seqFile)
    fmt.Fprintf(w, "n,avg_time_in_microseconds\n")

    for i := 2; i <= 50; i++ {
        start := time.Now()
        for j := 0; j < 1000; j++ {
            a := createRandomMatrix(i)
            b := createRandomMatrix(i)
            multiply(a, b, 1)
        }
        elapsed := time.Since(start)
        elapsed /= 1000
        fmt.Fprintf(w, "%v,%v\n", i, float32(elapsed) / 1000)
    }

    w.Flush()

    fmt.Println("Seq_exe.csv written.")

    // parallel
    pw := bufio.NewWriter(paralFile)
    fmt.Fprintf(pw, "n,avg_time_in_microseconds\n")

    for i := 2; i <= 50; i++ {
        start := time.Now()
        for j := 0; j < 1000; j++ {
            a := createRandomMatrix(i)
            b := createRandomMatrix(i)
            multiply(a, b, i)
        }
        elapsed := time.Since(start)
        elapsed /= 1000
        fmt.Fprintf(pw, "%v,%v\n", i, float32(elapsed) / 1000)
    }

    pw.Flush()

    fmt.Println("Paral_exe.csv written.")
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

                var dotProduct float32 = 0.0
                for k := range a[0] {
                    dotProduct += a[row][k] * b[k][col]
                }

                prod[row][col] = dotProduct

            }
            wg.Done()
        }(i + 1, pos, taskDistribution[i])
        pos += taskDistribution[i]
    }

    wg.Wait()

    return prod, nil
}

func createRandomMatrix(size int) [][]float32 {
    var randMatrix = make([][]float32, size)
    for i := 0; i < size; i++ {
        randMatrix[i] = make([]float32, size)
        for j := 0; j < size; j++ {
            randMatrix[i][j] = (rand.Float32() * 200) - 100
        }
    }
    return randMatrix
}
