# Params Linter
Detects if `.go` files contains multiple parameters with the same type

## Usage
```
go install github.com/lycb/go-params-linter/cmd/go-params-linter 
cd <your-go-project>
go-params-linter ./...
```

run with the `-fix` flag to automatically merge params with the same type

## Known Limitations with `-fix`
`-fix` flag current remove free-floating comments inside of a method, use with caution 😅

## Demo
https://user-images.githubusercontent.com/32417800/191370161-a73eccc7-ad42-4743-999f-172d95e2c448.mov