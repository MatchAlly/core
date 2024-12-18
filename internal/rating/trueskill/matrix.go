package trueskill

import (
	"errors"
	"fmt"
	"math"
)

const (
	FractionalDigitsToRoundTo = 10    // The precision for equality checks
	ErrorTolerance            = 1e-10 // The maximum allowed difference for equality checks
)

type matrix struct {
	Rows, Cols int
	data       [][]float64
}

func newMatrix(rows, cols int, data ...float64) (*matrix, error) {
	if rows <= 0 || cols <= 0 {
		return nil, errors.New("rows and cols must be positive")
	}

	m := &matrix{
		Rows: rows,
		Cols: cols,
		data: make([][]float64, rows),
	}

	for i := range m.data {
		m.data[i] = make([]float64, cols)
	}

	if len(data) > 0 {
		if len(data) != rows*cols {
			return nil, errors.New("number of provided values does not match dimensions")
		}

		index := 0
		for i := 0; i < rows; i++ {
			for j := 0; j < cols; j++ {
				m.data[i][j] = data[index]
				index++
			}
		}
	}

	return m, nil
}

func from2DArray(data [][]float64) (*matrix, error) {
	if len(data) == 0 {
		return nil, errors.New("empty input slice")
	}

	cols := len(data[0])
	for _, row := range data {
		if len(row) != cols {
			return nil, errors.New("inconsistent row lengths")
		}
	}

	return newMatrix(len(data), cols, flatten(data)...)
}

func flatten(data [][]float64) []float64 {
	result := make([]float64, 0, len(data)*len(data[0]))
	for _, row := range data {
		result = append(result, row...)
	}
	return result
}

func (m *matrix) Get(row, col int) (float64, error) {
	if !m.isValidIndex(row, col) {
		return 0, fmt.Errorf("invalid index: (%d, %d)", row, col)
	}
	return m.data[row][col], nil
}

func (m *matrix) Set(row, col int, value float64) error {
	if !m.isValidIndex(row, col) {
		return fmt.Errorf("invalid index: (%d, %d)", row, col)
	}
	m.data[row][col] = value
	return nil
}

func (m *matrix) Transpose() *matrix {
	result, _ := newMatrix(m.Cols, m.Rows)
	for i := 0; i < m.Rows; i++ {
		for j := 0; j < m.Cols; j++ {
			result.data[j][i] = m.data[i][j]
		}
	}
	return result
}

func (m *matrix) IsSquare() bool {
	return m.Rows == m.Cols
}

func (m *matrix) Determinant() (float64, error) {
	if !m.IsSquare() {
		return 0, errors.New("determinant is only defined for square matrices")
	}

	if m.Rows == 1 {
		return m.data[0][0], nil
	}

	if m.Rows == 2 {
		a, b := m.data[0][0], m.data[0][1]
		c, d := m.data[1][0], m.data[1][1]
		return a*d - b*c, nil
	}

	var det float64
	for j := 0; j < m.Cols; j++ {
		det += m.data[0][j] * m.cofactor(0, j)
	}

	return det, nil
}

func (m *matrix) Adjugate() (*matrix, error) {
	if !m.IsSquare() {
		return nil, errors.New("adjugate is only defined for square matrices")
	}

	result, _ := newMatrix(m.Rows, m.Cols)
	for i := 0; i < m.Rows; i++ {
		for j := 0; j < m.Cols; j++ {
			result.data[j][i] = m.cofactor(i, j)
		}
	}

	return result, nil
}

func (m *matrix) Inverse() (*matrix, error) {
	if !m.IsSquare() {
		return nil, errors.New("inverse is only defined for square matrices")
	}

	det, err := m.Determinant()
	if err != nil {
		return nil, err
	}

	if math.Abs(det) < 1e-10 {
		return nil, errors.New("matrix is singular or nearly singular")
	}

	adj, err := m.Adjugate()
	if err != nil {
		return nil, err
	}

	return scalarMultiply(1.0/det, adj), nil
}

func scalarMultiply(scalar float64, m *matrix) *matrix {
	result, _ := newMatrix(m.Rows, m.Cols)
	for i := 0; i < m.Rows; i++ {
		for j := 0; j < m.Cols; j++ {
			result.data[i][j] = scalar * m.data[i][j]
		}
	}
	return result
}

func add(a, b *matrix) (*matrix, error) {
	if a.Rows != b.Rows || a.Cols != b.Cols {
		return nil, errors.New("matrices must have the same dimensions for addition")
	}

	result, _ := newMatrix(a.Rows, a.Cols)
	for i := 0; i < a.Rows; i++ {
		for j := 0; j < a.Cols; j++ {
			result.data[i][j] = a.data[i][j] + b.data[i][j]
		}
	}
	return result, nil
}

func multiply(a, b *matrix) (*matrix, error) {
	if a.Cols != b.Rows {
		return nil, errors.New("number of columns in a must match number of rows in b")
	}

	result, _ := newMatrix(a.Rows, b.Cols)
	for i := 0; i < a.Rows; i++ {
		for j := 0; j < b.Cols; j++ {
			for k := 0; k < a.Cols; k++ {
				result.data[i][j] += a.data[i][k] * b.data[k][j]
			}
		}
	}
	return result, nil
}

func (m *matrix) cofactor(row, col int) float64 {
	minor, _ := m.getMinor(row, col)
	sign := math.Pow(-1, float64(row+col))
	det, _ := minor.Determinant()
	return sign * det
}

func (m *matrix) getMinor(row, col int) (*matrix, error) {
	if !m.isValidIndex(row, col) {
		return nil, fmt.Errorf("invalid index: (%d, %d)", row, col)
	}

	minor, _ := newMatrix(m.Rows-1, m.Cols-1)
	i := 0
	for r := 0; r < m.Rows; r++ {
		if r == row {
			continue
		}
		j := 0
		for c := 0; c < m.Cols; c++ {
			if c == col {
				continue
			}
			minor.data[i][j] = m.data[r][c]
			j++
		}
		i++
	}

	return minor, nil
}

func (m *matrix) isValidIndex(row, col int) bool {
	return row >= 0 && row < m.Rows && col >= 0 && col < m.Cols
}
