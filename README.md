# cf
Codeforces command line client.

The goal of this tool is to make your life as a problem solver a little bit easier, you focus on solving the main problems, cf takes care of setting your environment up.

With cf you can easily download test cases, generate boilerplate code for various programming languages, automate tests, etc.

## Usage examples

###Download test cases for a single problem:

    $ cf parse http://codeforces.com/contest/459/problem/C
    $ tree -a
    .
    ├── .in_0.txt
    ├── .in_1.txt
    ├── .out_0.txt
    ├── .out_1.txt
    └── .settings.yml

    0 directories, 5 files

###Download test cases for all problems in a contest:

    $ cf setup 459
    Problem A is ready!
    Problem B is ready!
    Problem C is ready!
    Problem D is ready!
    Problem E is ready!

    $ tree -a
    .
    └── CodeforcesRound261Div.2
        ├── A
        │   ├── .in_0.txt
        │   ├── ...
        │   └── .settings.yml
        ├── B
        │   ├── .in_0.txt
        │   ├── ...
        │   └── .settings.yml
        ├── C
        │   ├── .in_0.txt
        │   ├── ...
        │   └── .settings.yml
        ├── D
        │   ├── .in_0.txt
        │   ├── ...
        │   └── .settings.yml
        └── E
            ├── .in_0.txt
            ├── ...
            └── .settings.yml

    6 directories, 33 files

###Introducing the cf global configuration file
cf uses the `~/.cf.yml` file to store global settings, e.g.:

    ---
    template.cpp: /home/user/templates/template.cpp
    template.go: /home/user/templates/template.go

###Generate sample solution:

    $ cf gen main.go
    $ cat main.go
    package main

    import ()

    func main() {
    }

If you specify a valid template file in your `~/.cf.yml`, the sample solution will be a copy of that template.
