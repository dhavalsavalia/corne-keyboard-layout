package firmware

import (
	"os"
	"path/filepath"
	"sort"
	"time"
)

// File represents a firmware file
type File struct {
	Name string
	Path string
	Size int64
}

// Build represents a firmware build (dated directory)
type Build struct {
	Date     string // YYYYMMDD format
	Path     string
	Files    []File
	HasLeft  bool
	HasRight bool
	HasReset bool
}

// Scanner scans firmware directories
type Scanner struct {
	baseDir string
}

// NewScanner creates a new firmware scanner
func NewScanner(baseDir string) *Scanner {
	return &Scanner{baseDir: baseDir}
}

// Scan scans for firmware builds and returns them sorted by date (newest first)
func (s *Scanner) Scan() ([]Build, error) {
	entries, err := os.ReadDir(s.baseDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []Build{}, nil
		}
		return nil, err
	}

	var builds []Build
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		// Check if directory name looks like a date (YYYYMMDD)
		if len(name) != 8 || !isNumeric(name) {
			continue
		}

		buildPath := filepath.Join(s.baseDir, name)
		build := Build{
			Date: name,
			Path: buildPath,
		}

		// Scan for UF2 files
		files, err := os.ReadDir(buildPath)
		if err != nil {
			continue
		}

		for _, f := range files {
			if f.IsDir() {
				continue
			}
			if filepath.Ext(f.Name()) != ".uf2" {
				continue
			}

			info, err := f.Info()
			if err != nil {
				continue
			}

			file := File{
				Name: f.Name(),
				Path: filepath.Join(buildPath, f.Name()),
				Size: info.Size(),
			}
			build.Files = append(build.Files, file)

			switch f.Name() {
			case "corne_left.uf2":
				build.HasLeft = true
			case "corne_right.uf2":
				build.HasRight = true
			case "settings_reset.uf2":
				build.HasReset = true
			}
		}

		if len(build.Files) > 0 {
			builds = append(builds, build)
		}
	}

	// Sort by date descending (newest first)
	sort.Slice(builds, func(i, j int) bool {
		return builds[i].Date > builds[j].Date
	})

	return builds, nil
}

// FormatDate formats YYYYMMDD to human-readable format
func FormatDate(date string) string {
	if len(date) != 8 {
		return date
	}
	t, err := time.Parse("20060102", date)
	if err != nil {
		return date
	}
	return t.Format("2006-01-02")
}

// FormatSize formats bytes to human-readable size
func FormatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return formatInt(bytes) + " B"
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return formatFloat(float64(bytes)/float64(div)) + " " + string("KMGTPE"[exp]) + "B"
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func formatInt(n int64) string {
	if n == 0 {
		return "0"
	}
	var result []byte
	for n > 0 {
		result = append([]byte{byte('0' + n%10)}, result...)
		n /= 10
	}
	return string(result)
}

func formatFloat(f float64) string {
	// Simple formatting: one decimal place
	intPart := int64(f)
	decPart := int64((f - float64(intPart)) * 10)
	return formatInt(intPart) + "." + formatInt(decPart)
}
