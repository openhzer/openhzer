package TAAG

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	//go:embed "ANSI Shadow.flf"
	TemplateAnsiShadowFont []byte
)

// FLFFont 解析后的 Figlet 字体结构
type FLFFont struct {
	Signature      string
	HardBlank      rune
	Height         int
	Baseline       int
	MaxLength      int
	OldLayout      int
	CommentLines   int
	PrintDirection int
	FullLayout     int
	CodetagCount   int
	Characters     map[rune][]string
}

func ParseFLFFromFile(filename string) (*FLFFont, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return nil, fmt.Errorf("empty file")
	}
	return ParseFLF(scanner)
}

// ParseFLFFromBytes 从字节集解析 FLF
func ParseFLFFromBytes(data []byte) (*FLFFont, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	if !scanner.Scan() {
		return nil, fmt.Errorf("empty file")
	}
	return ParseFLF(scanner)
}

// ParseFLFFromString 从字符串解析 FLF
func ParseFLFFromString(data string) (*FLFFont, error) {
	scanner := bufio.NewScanner(strings.NewReader(data))
	if !scanner.Scan() {
		return nil, fmt.Errorf("empty file")
	}
	return ParseFLF(scanner)
}

// ParseFLF 从流解析 FLF
func ParseFLF(scanner *bufio.Scanner) (*FLFFont, error) {
	header := scanner.Text()
	parts := strings.Split(header, " ")
	if len(parts) < 6 {
		return nil, fmt.Errorf("invalid header format")
	}

	// 解析文件头
	font := &FLFFont{
		Signature:  parts[0],
		Characters: make(map[rune][]string),
	}

	// 提取硬空格字符（签名最后一个字符）
	if len(font.Signature) > 0 {
		font.HardBlank = rune(font.Signature[len(font.Signature)-1])
		font.Signature = font.Signature[:len(font.Signature)-1]
	}

	// 解析数值参数
	var nums []int
	for _, p := range parts[1:] {
		if p == "" {
			continue
		}
		num, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("invalid number in header: %v", err)
		}
		nums = append(nums, num)
	}

	if len(nums) < 8 {
		return nil, fmt.Errorf("not enough parameters in header")
	}

	font.Height = nums[0]
	font.Baseline = nums[1]
	font.MaxLength = nums[2]
	font.PrintDirection = nums[3]
	font.CommentLines = nums[4]
	font.FullLayout = nums[5]
	font.CodetagCount = nums[6]
	// nums[7] 是 horizontal smushing，忽略

	// 跳过注释行
	for i := 0; i < font.CommentLines; i++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("unexpected EOF while skipping comments")
		}
	}

	// 读取字符定义
	currentChar := ' '
	charLines := make([]string, font.Height)
	lineCount := 0

	for scanner.Scan() {
		line := scanner.Text()

		// 处理行尾的@和@@标记
		if strings.HasSuffix(line, "@@") {
			line = line[:len(line)-2]
		} else if strings.HasSuffix(line, "@") {
			line = line[:len(line)-1]
		}

		charLines[lineCount] = line
		lineCount++

		if lineCount == font.Height {
			// 处理硬空格替换
			for i := range charLines {
				charLines[i] = strings.ReplaceAll(charLines[i], string(font.HardBlank), " ")
			}

			font.Characters[currentChar] = charLines
			currentChar++
			charLines = make([]string, font.Height)
			lineCount = 0
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return font, nil
}

// Render 渲染文本使用解析的字体
func (f *FLFFont) Render(text string) string {
	var result strings.Builder

	for h := 0; h < f.Height; h++ {
		for _, c := range text {
			if char, ok := f.Characters[c]; ok && h < len(char) {
				result.WriteString(char[h])
			} else {
				// 如果字符未定义，使用空格代替
				result.WriteString(strings.Repeat(" ", f.MaxLength))
			}
		}
		result.WriteString("\n")
	}

	return result.String()
}
