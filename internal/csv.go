package internal

import (
	"encoding/csv"
	"io"
	"iter"

	"github.com/cockroachdb/errors"
)

type Result[T any] struct {
	Value   T
	LineNum int
	Error   error
}

func ParseCSV[T any](reader io.Reader, includesHeader bool, fromFunc func(data []string, headerMap map[string]int) (T, error)) iter.Seq[Result[T]] {

	return func(yield func(Result[T]) bool) {
		lineNum := 1
		csvReader := csv.NewReader(reader)
		csvReader.LazyQuotes = true
		headerMap := make(map[string]int)

		if includesHeader {
			headers, err := csvReader.Read()
			if err != nil {
				yield(Result[T]{
					LineNum: lineNum,
					Error:   errors.Wrap(err, "failed to read CSV headers"),
				})
				return
			}
			for i, h := range headers {
				if _, ok := headerMap[h]; ok {
					yield(Result[T]{
						LineNum: lineNum,
						Error:   errors.Errorf("duplicate header column found: %s", h),
					})
					return
				}
				headerMap[h] = i
			}
		}

		for {
			lineNum++
			record, err := csvReader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				yield(Result[T]{
					LineNum: lineNum,
					Error:   errors.Wrapf(err, "failed to read CSV line %d", lineNum),
				})
				return
			}

			data, err := fromFunc(record, headerMap)
			if err != nil {
				yield(Result[T]{
					LineNum: lineNum,
					Error:   errors.Wrapf(err, "failed to parse CSV line %d", lineNum),
				})
				return
			}

			if !yield(Result[T]{
				Value:   data,
				LineNum: lineNum,
				Error:   nil,
			}) {
				return
			}
		}
	}
}
