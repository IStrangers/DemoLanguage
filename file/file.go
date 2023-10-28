package file

import (
	"sort"
	"sync"
)

type Index int

type Position struct {
	FileName string
	Line     int
	Column   int
}

type File struct {
	Lock              sync.RWMutex
	BaseOffset        int
	Name              string
	Content           string
	LineOffsets       []int
	LastScannedOffset int
}

func (file *File) PositionByIndex(index Index) *Position {
	return file.Position(int(index) - file.BaseOffset)
}

func (file *File) Position(offset int) *Position {
	var line int
	var lineOffsets []int

	lock := file.Lock
	if offset > file.LastScannedOffset {
		lock.Lock()
		lineOffsets, line = file.scanToOffset(offset)
		lock.Unlock()
	} else {
		lock.RLock()
		lineOffsets = file.LineOffsets
		lock.RUnlock()
		line = sort.Search(len(lineOffsets), func(index int) bool {
			return lineOffsets[index] > offset
		}) - 1
	}

	col := offset + 1
	if len(lineOffsets) > 0 {
		col -= lineOffsets[line]
	}
	row := line + 2

	return &Position{
		FileName: file.Name,
		Line:     row,
		Column:   col,
	}
}

func (file *File) scanToOffset(offset int) ([]int, int) {
	for file.LastScannedOffset < offset {
		lineOffset := file.findLineOffset(file.Content[file.LastScannedOffset:])
		if lineOffset == -1 {
			file.LastScannedOffset = len(file.Content)
			return file.LineOffsets, len(file.LineOffsets) - 1
		}
		file.LastScannedOffset = file.LastScannedOffset + lineOffset
		file.LineOffsets = append(file.LineOffsets, file.LastScannedOffset)
	}

	if file.LastScannedOffset == offset {
		return file.LineOffsets, len(file.LineOffsets) - 1
	}
	return file.LineOffsets, len(file.LineOffsets) - 2
}

func (file *File) findLineOffset(content string) int {
	for index, ch := range content {
		switch ch {
		case '\r':
			if index < len(content)-1 && content[index+1] == '\n' {
				return index + 2
			}
			return index + 1
		case '\n':
			return index + 1
		}
	}
	return -1
}

func CreateFile(baseOffset int, name string, content string) *File {
	return &File{BaseOffset: baseOffset, Name: name, Content: content}
}
