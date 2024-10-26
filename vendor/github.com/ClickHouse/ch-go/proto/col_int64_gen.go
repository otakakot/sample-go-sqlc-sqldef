// Code generated by ./cmd/ch-gen-col, DO NOT EDIT.

package proto

// ColInt64 represents Int64 column.
type ColInt64 []int64

// Compile-time assertions for ColInt64.
var (
	_ ColInput  = ColInt64{}
	_ ColResult = (*ColInt64)(nil)
	_ Column    = (*ColInt64)(nil)
)

// Rows returns count of rows in column.
func (c ColInt64) Rows() int {
	return len(c)
}

// Reset resets data in row, preserving capacity for efficiency.
func (c *ColInt64) Reset() {
	*c = (*c)[:0]
}

// Type returns ColumnType of Int64.
func (ColInt64) Type() ColumnType {
	return ColumnTypeInt64
}

// Row returns i-th row of column.
func (c ColInt64) Row(i int) int64 {
	return c[i]
}

// Append int64 to column.
func (c *ColInt64) Append(v int64) {
	*c = append(*c, v)
}

// Append int64 slice to column.
func (c *ColInt64) AppendArr(vs []int64) {
	*c = append(*c, vs...)
}

// LowCardinality returns LowCardinality for Int64 .
func (c *ColInt64) LowCardinality() *ColLowCardinality[int64] {
	return &ColLowCardinality[int64]{
		index: c,
	}
}

// Array is helper that creates Array of int64.
func (c *ColInt64) Array() *ColArr[int64] {
	return &ColArr[int64]{
		Data: c,
	}
}

// Nullable is helper that creates Nullable(int64).
func (c *ColInt64) Nullable() *ColNullable[int64] {
	return &ColNullable[int64]{
		Values: c,
	}
}

// NewArrInt64 returns new Array(Int64).
func NewArrInt64() *ColArr[int64] {
	return &ColArr[int64]{
		Data: new(ColInt64),
	}
}