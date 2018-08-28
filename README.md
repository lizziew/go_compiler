### About

Interpreter and compiler for toy language

Written in Go

See http://datasieve.blogspot.com/2018/08/building-toy-language-interpreter-in-go.html and http://datasieve.blogspot.com/2018/08/building-toy-language-compiler-in-go.html for more details

Based on https://interpreterbook.com/ 

### Features 

Supports:
- integers, booleans, strings, arrays, hashmaps 
- prefix, infix operators
- index operators
- conditionals
- global and local bindings 
- first class functions
- return statements
- closures 

### How to Run

Build: 
```shell 
➜ go build -o toy
```

Run in interpreter mode: 
```shell
➜ ./toy -engine=eval
`Welcome to the Monkey programming language, elizabethwei!`
`Feel free to type in commands. Engine = eval`
>> 
```

Run in compiler mode:
```shell
➜ ./toy -engine=vm
Welcome to the Monkey programming language, elizabethwei!
Feel free to type in commands. Engine = vm
>> 
```

### Logging 

Run with or without intermediate print statements: 

![alt text](https://github.com/lizziew/go_interpreter/blob/master/img/without_print.png)

![alt text](https://github.com/lizziew/go_interpreter/blob/master/img/withprint.png)

