package main

import (
    "fmt"
)

func main() {
    fmt.Println("Starting the program")
    
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered from panic:", r)
        }
    }()
    
    fmt.Println("Calling function that may panic")
    panickingFunction()
    
    fmt.Println("This line won't be reached")
}

func panickingFunction() {
    fmt.Println("About to panic")
    panic("Something went wrong")
}