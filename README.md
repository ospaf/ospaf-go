# ospaf-go
This is a light weight open source analysing project based on the [github API](https://developer.github.com/v3), written in Golang.

- github: github API implemented in Golang
- lib: the wrapper of the useful functions
- [programs](#programes): independent programs used to analyse org/developer/project based on github APIs.


##programes
Every program should have:
- `README.md`:  explain what the program do and how to use it.
- `Makefile`:  with `all` session to compile all the sub-programs

To make it easier to integrated with QA automation, it is recommended sub-programs should following these nameing rules:
- `collect.go`:  collect data from github
- `analyse.go`:  analyse data collected by collect.go
- `report.go`:  `report --html` and `report --plain`
  
