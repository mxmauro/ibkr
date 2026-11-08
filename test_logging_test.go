package ibkr_test

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// -----------------------------------------------------------------------------

var testLogFile = flag.String("log-file", "", "write test and client event logs to a per-process file derived from this path")

var (
	testLogSinkOnce sync.Once
	testLogSinkInst *testLogSink
	testLogSinkErr  error
)

// -----------------------------------------------------------------------------

type testLogSink struct {
	path string
	file *os.File
	mtx  sync.Mutex
}

// -----------------------------------------------------------------------------

func getTestLogSink(tb testing.TB) *testLogSink {
	tb.Helper()

	testLogSinkOnce.Do(func() {
		basePath := strings.TrimSpace(*testLogFile)
		if basePath == "" {
			return
		}

		path := deriveTestLogFilePath(basePath, os.Getpid())
		dir := filepath.Dir(path)
		if dir != "." {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				testLogSinkErr = fmt.Errorf("create log directory %q: %w", dir, err)
				return
			}
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
		if err != nil {
			testLogSinkErr = fmt.Errorf("open test log file %q: %w", path, err)
			return
		}

		testLogSinkInst = &testLogSink{
			path: path,
			file: file,
		}
	})

	if testLogSinkErr != nil {
		tb.Fatal(testLogSinkErr)
	}
	return testLogSinkInst
}

func deriveTestLogFilePath(basePath string, pid int) string {
	ext := filepath.Ext(basePath)
	if ext == "" {
		return fmt.Sprintf("%s.%d.log", basePath, pid)
	}
	return fmt.Sprintf("%s.%d%s", strings.TrimSuffix(basePath, ext), pid, ext)
}

func (s *testLogSink) Close() error {
	if s == nil || s.file == nil {
		return nil
	}
	return s.file.Close()
}

func (s *testLogSink) writeRecord(testName string, level string, msg string) error {
	if s == nil {
		return nil
	}

	msg = strings.ReplaceAll(msg, "\r\n", "\n")
	msg = strings.TrimRight(msg, "\n")
	if msg == "" {
		return nil
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()

	for _, line := range strings.Split(msg, "\n") {
		if _, err := fmt.Fprintf(s.file, "%s [%s] [%s] %s\n", time.Now().Format(time.RFC3339Nano), level, testName, line); err != nil {
			return err
		}
	}

	return nil
}

func sprintTestArgs(args ...any) string {
	return strings.TrimSuffix(fmt.Sprintln(args...), "\n")
}

func recordTestLog(tb testing.TB, level string, msg string) {
	tb.Helper()

	sink := getTestLogSink(tb)
	if sink == nil {
		return
	}
	if err := sink.writeRecord(tb.Name(), level, msg); err != nil {
		tb.Fatalf("write test log: %v", err)
	}
}

func testLog(tb testing.TB, args ...any) {
	tb.Helper()
	tb.Log(args...)
	recordTestLog(tb, "INFO", sprintTestArgs(args...))
}

func testLogf(tb testing.TB, format string, args ...any) {
	tb.Helper()
	msg := fmt.Sprintf(format, args...)
	tb.Logf(format, args...)
	recordTestLog(tb, "INFO", msg)
}

func testEventLog(tb testing.TB, msg string) {
	tb.Helper()
	tb.Log(msg)
	recordTestLog(tb, "EVENT", msg)
}

func testError(tb testing.TB, args ...any) {
	tb.Helper()
	tb.Error(args...)
	recordTestLog(tb, "ERROR", sprintTestArgs(args...))
}

func testFatal(tb testing.TB, args ...any) {
	tb.Helper()
	recordTestLog(tb, "FATAL", sprintTestArgs(args...))
	tb.Fatal(args...)
}

func testFatalf(tb testing.TB, format string, args ...any) {
	tb.Helper()
	msg := fmt.Sprintf(format, args...)
	recordTestLog(tb, "FATAL", msg)
	tb.Fatalf(format, args...)
}

func testSkip(tb testing.TB, args ...any) {
	tb.Helper()
	recordTestLog(tb, "SKIP", sprintTestArgs(args...))
	tb.Skip(args...)
}

// -----------------------------------------------------------------------------

func TestDeriveTestLogFilePath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		basePath string
		pid      int
		want     string
	}{
		{
			name:     "with extension",
			basePath: filepath.Join("logs", "ibkr-test.log"),
			pid:      1234,
			want:     filepath.Join("logs", "ibkr-test.1234.log"),
		},
		{
			name:     "without extension",
			basePath: filepath.Join("logs", "ibkr-test"),
			pid:      1234,
			want:     filepath.Join("logs", "ibkr-test.1234.log"),
		},
		{
			name:     "multiple dots",
			basePath: filepath.Join("logs", "ibkr.test.trace.txt"),
			pid:      77,
			want:     filepath.Join("logs", "ibkr.test.trace.77.txt"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := deriveTestLogFilePath(tt.basePath, tt.pid); got != tt.want {
				t.Fatalf("unexpected log path: got %q want %q", got, tt.want)
			}
		})
	}
}

func TestTestLogSinkWriteRecordSerializesConcurrentWriters(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "ibkr-test.log")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		t.Fatal(err)
	}

	sink := &testLogSink{
		path: path,
		file: file,
	}

	const records = 32

	var wg sync.WaitGroup
	wg.Add(records)
	for idx := 0; idx < records; idx++ {
		idx := idx
		go func() {
			defer wg.Done()
			if err := sink.writeRecord("TestConcurrent", "INFO", fmt.Sprintf("message-%02d", idx)); err != nil {
				t.Errorf("writeRecord failed: %v", err)
			}
		}()
	}
	wg.Wait()

	if err := sink.Close(); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != records {
		t.Fatalf("unexpected line count: got %d want %d", len(lines), records)
	}

	for idx := 0; idx < records; idx++ {
		if !strings.Contains(string(data), fmt.Sprintf("message-%02d", idx)) {
			t.Fatalf("missing message-%02d in log output", idx)
		}
	}
}
