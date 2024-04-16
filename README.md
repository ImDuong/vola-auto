# vola-auto
The ultimate streamline for Volatility 3. Speed up process of memory artifacts extraction phase

# Why this
- Why this project is developed? We want to accelerate the artifacts extraction phase and combine automatic artifact analyzing strategy. Volatility 3 focuses on reading memory, vola-auto focuses on extracting & analyzing artifacts. 
- Why not a simple project using Volatility SDK? Frankly, Volatility 3 is written in Python, which makes the wrapper program should also be written in Python. And, vola-auto is not intended to replace Volatility 3, but rather become an extra flavor for researchers who had already installed and been familiar with one of the best memory forensic tools. On the other hand, with golang, it's easy to cross compile as standalone binaries for multiple OSes.

# Features
1. Auto run common plugins: info, filescan, process, etc. Auto dump common artifacts file: MFT, Logfile, prefetch, etc. Auto run customized artifacts analytics
2. Add regex for dumping files (currently Vol3 does not support it, which is different from Vol2)
3. Run multiple commands parallelly from a file

- Note: just support Vol3 and Windows yet

# Prerequisite
- Python3
- Volatility 3 (Vola-auto tested with Volatility 3 Framework 2.5.2 to 2.7.0)
    - Our tool relies on Volatility 3, a memory forensics framework, for analyzing memory dumps. Users need to obtain Volatility 3 separately and comply with its licensing terms.
    - **License**: Volatility 3 is licensed under the Volatility Software License Version 1.0. Please review the [Volatility Software License Version 1.0](https://www.volatilityfoundation.org/license/vsl-v1.0) for details on your obligations when obtaining and using Volatility 3.

# Getting started
## Environment
- Install requirements for Volatility 3 or activate the env that you already setup for Volatility 3
- If running golang code directly, install go modules first with `go install`

## Auto Streamline
- Run Volatility 3 auto streamline with `--vol` pointing to volatility 3 folder, and `--file or -f` pointing to memory dump file

    ```
    go run cmd\main.go --vol <path_to_volatility3> -f <path_to_memory_dump> -o <output_folder>
    ```

    - if `-o` is not specified, vola-auto will generate folder `artifacts` in folder containing `path_to_memory_dump`

## Dump files with regex
- To dump files with regex, use subcommand `dumpfiles` with `-reg` flag to pass regex.

    ```
    go run cmd\main.go --vol <path_to_volatility3> -f <path_to_memory_dump> -o <output_folder> dumpfiles -reg "SCHEDLGU\.TXT$"
    ```

    - if `output_folder` does not contain filescan.txt (output file of filescan plugin) yet, vola-auto will run filescan plugin first to write `filescan.txt`, then starting to dump files
    - To specify a specific `filescan.txt`, use `-fs` flag

        ```
        go run cmd\main.go --vol <path_to_volatility3> -f <path_to_memory_dump> -o <output_folder> dumpfiles -reg "SCHEDLGU\.TXT$" -fs <path_to_filescan.txt>
        ```

## Execute batch of commands parallely
- Write commands in each line in a file. For example:
    ```
    windows.pstree.PsTree
    windows.psscan.PsScan
    windows.pslist.PsList
    ```

- Use subcommand `batch`

    ```
    go run cmd\main.go --vol <path_to_volatility3> -f <path_to_memory_dump> -o <output_folder> batch -f <path_to_command_file>
    ```

    - Results of command will be logged to files in `temp` folder inside `output_folder`

# Tips
- To run vola-auto in verbose mode, run the program with environment variable `DEBUG` having value as `1`

    ```
    # windows
    SET DEBUG=1 && go run cmd\main.go --vol <path_to_volatility3> -f <path_to_memory_dump>
    ```

## Acknowledgments

We would like to acknowledge the Volatility Foundation for developing Volatility 3, which our tool utilizes for memory forensics analysis.
- Volatility Foundation. Volatility 3 [Computer software]. https://github.com/volatilityfoundation/volatility3
