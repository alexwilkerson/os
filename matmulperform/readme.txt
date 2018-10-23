How to run matrixmul.exe on Windows:

- Open a Command Prompt.
- Navigate to the directory where matrixmul.exe is located.
- matrixmul.exe operates with the following command line arguments:
    > matrixmul.exe [threads] [inputfile]
    e.g.
    > matrixmul.exe 4 input0.txt

How to run matmulperform.exe on Windows:

- Open a Command Prompt.
- Navigate to the directory where matmulperform.exe is located.
- matmulperform.exe takes no command line arguments. Use as follows:
    > matmulperform.exe
- This operation will result in the following output files:
    - Seq_exe.csv (Sequential matrix multiplication)
    - Paral_exe.csv (Parallel matrix multiplication)
