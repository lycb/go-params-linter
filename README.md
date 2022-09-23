# Params Linter
Detects if `.go` files contains multiple parameters with the same type

## Usage
```
go install github.com/lycb/go-params-linter/cmd/go-params-linter 
cd <your-go-project>
go-params-linter ./...
```

run with the `-fix` flag to automatically merge params with the same type

## Known Limitations with `-fix` use with caution 😅
1. `-fix` flag current [remove free-floating comments](https://github.com/golang/go/issues/20744) inside of a method
2. does not fix everything that was flagged from the linting
3. does not recursively fix params so need to run multiple times to fix

## Demo
https://user-images.githubusercontent.com/32417800/191370161-a73eccc7-ad42-4743-999f-172d95e2c448.mov
