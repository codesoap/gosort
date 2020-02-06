`gosort` sorts function definitions in a `*.go` file automatically.

**This is an experiment and the tool might never mature or its
behavior might change substantially!**

It uses these rules:
1. Sort by topology (place callees after their callers)
2. Sort by call order (first called first; only looks for calls within
   the given file)
3. Sort alphabetically


# Usage
```cosole
$ gosort list.go
ListFunctionCallsWithinFile
removeNonLocal
extractFunctionCalls
appendNewFunctionCalls
appendNewFunctionCall
```

# TODO
Right now the new order is just printed out and not applied. A `-w`
flag, to automatically apply the new order to the file, could be useful.

Maybe these features could prove useful:
- Put exported functions to the top of the file.
- Group methods for a certain struct together.
- Arrange getters and setters in the order their fields are arranged in
  the struct.
- Not only arrange functions, but constant definitions and other
  top-level declarations as well.
