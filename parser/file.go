package parser

import (
	"sort"
	"strings"
	"sync"
)

type Index int

func (parser *Parser) IndexOf(offset int) Index {
	return Index(parser.baseOffset + offset)
}

type Position struct {
	FileName string
	Line     int
	Column   int
}

func (parser *Parser) Position(index Index) *Position {
	return parser.file.Position(int(index) - parser.baseOffset)
}

type File struct {
	Lock              sync.RWMutex
	BaseOffset        int
	Name              string
	Content           string
	LineOffsets       []int
	LastScannedOffset int
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
		offsetIndex := strings.Index(file.Content[file.LastScannedOffset:], "\r\n")
		if offsetIndex == -1 {
			file.LastScannedOffset = len(file.Content)
			return file.LineOffsets, len(file.LineOffsets) - 1
		}
		file.LastScannedOffset = file.LastScannedOffset + (offsetIndex + 2)
		file.LineOffsets = append(file.LineOffsets, file.LastScannedOffset)
	}

	if file.LastScannedOffset == offset {
		return file.LineOffsets, len(file.LineOffsets) - 1
	}
	return file.LineOffsets, len(file.LineOffsets) - 2
}

func CreateFile(baseOffset int, name string, content string) *File {
	return &File{BaseOffset: baseOffset, Name: name, Content: content}
}
