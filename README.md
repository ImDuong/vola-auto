# vola-auto
The ultimate streamline for Volatitlity 3. Speed up process of memory artifacts extraction phase

# Why this
- Why this project is developed? I want to accelerate the artifacts extraction phase and combine automatic artifact analytics strategy. Volatility 3 focuses on reading memory, vola-auto focuses on extracting & analyzing artifacts. 
- Why not a simple project using Volatility SDK? Frankly, Volatility 3 is written in Python, which makes the wrapper program should also be written in Python. And this is not what I want, because I just want to use Golang, and it's easy to cross compile as standalone binaries for multiple OSes. This is not intended to replace Volatility 3, but rather become an extra flavor for researchers who had already installed and been familiar with one of the best memory forensic tools. 

# Features
1. Auto run common plugins: info, filescan, process, etc. Auto dump common artifacts file: MFT, Logfile, prefetch, etc. Auto run customized artifacts analytics
2. [TODO] Add regex for dumping files (currently Vol3 does not support it, which is different from Vol2)
3. [TODO] Run multiple commands parallelly

- Note: just support Vol3 and Windows yet

# Prerequisite
- python3
- Volatility 3

# Getting started
1. Install requirements for Volatility 3 or activate the env that you already setup for Volatility 3
2. Run Volatility 3 auto streamline with 

    ```
    go run cmd\main.go -v <path_to_volatility3> -f <path_to_memory_dump> -o <output_folder>
    ```