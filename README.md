# cf
Codeforces command line client.

The goal of this tool is to make your life as a problem solver a little bit easier, you focus on solving the main problems, cf takes care of setting your environment up.

With cf you can easily download test cases, generate boilerplate code for various programming languages, automate tests, etc.

# Usage examples

Download test cases for a single problem:

    $ cf parse http://codeforces.com/contest/459/problem/C
    $ tree -a
    .
    ├── .in_0.txt
    ├── .in_1.txt
    ├── .out_0.txt
    └── .out_1.txt

    0 directories, 4 files


Download test cases for all problems in a contest:

    $ cf setup 459
    Problem A is ready!
    Problem B is ready!
    Problem C is ready!
    Problem D is ready!
    Problem E is ready!

    $ tree -a
    .
    ├── A
    │   ├── .in_0.txt
    │   ├── .in_1.txt
    │   ├── ...
    ├── B
    │   ├── .in_0.txt
    │   ├── .in_1.txt
    │   ├── ...
    ├── C
    │   ├── .in_0.txt
    │   ├── .in_1.txt
    │   ├── ...
    ├── D
    │   ├── .in_0.txt
    │   ├── .in_1.txt
    │   ├── ...
    └── E
        ├── .in_0.txt
        ├── .in_1.txt
        ├── ...

    5 directories, 28 files
