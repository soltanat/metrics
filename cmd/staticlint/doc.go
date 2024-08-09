// staticlint реализует анализ статических ошибок.
//
// Список анализатрров:
//
//  1. golang.org/x/tools/go/analysis/passes/printf
//     golang.org/x/tools/go/analysis/passes/shadow
//     golang.org/x/tools/go/analysis/passes/structtag
//
// 2. Все анализаторы класса SA https://staticcheck.io/docs/checks/
//
// Usage:
//
//	staticlint <path to project>
//
// Example:
//
//	staticlint -SA1000 <project path>
package main
