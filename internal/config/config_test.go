package config_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/mokan-r/golook/internal/config"
)

// MockReader is a mock implementation of the Reader interface
type MockReader struct {
	ReadFileFunc func(path string) ([]byte, error)
}

func (r *MockReader) ReadFile(path string) ([]byte, error) {
	if r.ReadFileFunc == nil {
		return nil, errors.New("not implemented")
	}
	return r.ReadFileFunc(path)
}

func TestReadConfig(t *testing.T) {
	type args struct {
		path   string
		reader config.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    *config.Config
		wantErr bool
	}{
		{
			name: "Empty config file",
			args: args{
				path: "testdata/empty.yaml",
				reader: &MockReader{
					ReadFileFunc: func(path string) ([]byte, error) {
						return []byte{}, nil
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Successful read and parse",
			args: args{
				path: "testdata/config.yaml",
				reader: &MockReader{
					ReadFileFunc: func(path string) ([]byte, error) {
						return []byte(`
  - path: /var/logs
    log_file: /var/logs/app.log
    include_regexp:
      - ".+\\.log$"
    exclude_regexp:
      - ".+\\.gz$"
    commands:
      - "echo 'New logs available at /var/logs/app.log'"
`), nil
					},
				},
			},
			want: &config.Config{
				Directories: map[string]config.DirectoryConfig{
					"/var/logs": {
						LogFile:       "/var/logs/app.log",
						IncludeRegexp: getRegexpArr([]string{".+\\.log$"}),
						ExcludeRegexp: getRegexpArr([]string{".+\\.gz$"}),
						Commands:      []string{"echo 'New logs available at /var/logs/app.log'"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Read file error",
			args: args{
				path: "nonexistent.yaml",
				reader: &MockReader{
					ReadFileFunc: func(path string) ([]byte, error) {
						return nil, errors.New("file not found")
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Empty directory path",
			args: args{
				path: "testdata/config.yaml",
				reader: &MockReader{
					ReadFileFunc: func(path string) ([]byte, error) {
						return []byte(`
  - path: 
    log_file: /var/logs/app.log
    include_regexp:
      - ".+\\.log$"
    exclude_regexp:
      - ".+\\.gz$"
    commands:
      - "echo 'New logs available at /var/logs/app.log'"
`), nil
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Parse config error",
			args: args{
				path: "testdata/config.yaml",
				reader: &MockReader{
					ReadFileFunc: func(path string) ([]byte, error) {
						return []byte(`invalid yaml`), nil
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Config file with only commands key",
			args: args{
				path: "testdata/commands-only.yaml",
				reader: &MockReader{
					ReadFileFunc: func(path string) ([]byte, error) {
						return []byte(`
commands:
  - echo "Hello, world!"
`), nil
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Config file with only directory key",
			args: args{
				path: "testdata/directories-only.yaml",
				reader: &MockReader{
					ReadFileFunc: func(path string) ([]byte, error) {
						return []byte(`
  - path: /var/logs
`), nil
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Config file with multiple directories",
			args: args{
				path: "testdata/multiple-directories.yaml",
				reader: &MockReader{
					ReadFileFunc: func(path string) ([]byte, error) {
						return []byte(`
- path: /var/logs
  log_file: /var/logs/app.log
  include_regexp:
    - ".+\\.log$"
  exclude_regexp:
    - ".+\\.gz$"
  commands:
    - "echo 'New logs available at /var/logs/app.log'"
- path: /opt
  include_regexp:
    - ".+\\.html$"
    - ".+\\.css$"
    - ".+\\.js$"
  commands:
    - "echo 'New logs available at /var/logs/app.log'"
`), nil
					},
				},
			},
			want: &config.Config{
				Directories: map[string]config.DirectoryConfig{
					"/var/logs": {
						LogFile:       "/var/logs/app.log",
						IncludeRegexp: getRegexpArr([]string{".+\\.log$"}),
						ExcludeRegexp: getRegexpArr([]string{".+\\.gz$"}),
						Commands:      []string{"echo 'New logs available at /var/logs/app.log'"},
					},
					"/opt": {
						IncludeRegexp: getRegexpArr([]string{".+\\.html$", ".+\\.css$", ".+\\.js$"}),
						Commands:      []string{"echo 'New logs available at /var/logs/app.log'"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Config file with invalid directory configuration",
			args: args{
				path: "testdata/invalid-directory.yaml",
				reader: &MockReader{
					ReadFileFunc: func(path string) ([]byte, error) {
						return []byte(`
- path: /var/logs
  include_regexp:
    - ".+\\.log$"
  exclude_regexp:
    - ".+\\.gz$"
  commands:
    - "echo 'New logs available at /var/logs/app.log'"
- invalid_key: invalid_value
`), nil
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Config file with invalid key",
			args: args{
				path: "testdata/invalid-key.yaml",
				reader: &MockReader{
					ReadFileFunc: func(path string) ([]byte, error) {
						return []byte(`
invalid_key: invalid_value
`), nil
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Config file with invalid include regular expression",
			args: args{
				path: "testdata/invalid-include-regexp.yaml",
				reader: &MockReader{
					ReadFileFunc: func(path string) ([]byte, error) {
						return []byte(`
- path: /var/logs
  include_regexp:
    - "+log$"
`), nil
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Config file with invalid exclude regular expression",
			args: args{
				path: "testdata/invalid-exclude-regexp.yaml",
				reader: &MockReader{
					ReadFileFunc: func(path string) ([]byte, error) {
						return []byte(`
- path: /var/logs
  exclude_regexp:
    - "+gz$"
`), nil
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := config.ReadConfig(tt.args.path, tt.args.reader)

			if (err != nil) != tt.wantErr {
				t.Errorf("got error = %v\nexpected error = %v", err, tt.wantErr)
				return
			}

			if !configDeepEqual(got, tt.want) {
				t.Errorf("got config = %v\nexpected config %v", got, tt.want)
			}
		})
	}
}

func configDeepEqual(x *config.Config, y *config.Config) bool {
	if x == nil && y == nil {
		return true
	}
	if x == nil || y == nil {
		return false
	}
	for key, v := range x.Directories {
		vWant, ok := y.Directories[key]
		if !ok {
			return false
		}

		if v.LogFile != vWant.LogFile {
			return false
		}

		for i := range v.Commands {
			if v.Commands[i] != vWant.Commands[i] {
				return false
			}
		}

		for i := range v.IncludeRegexp {
			if v.IncludeRegexp[i].String() != vWant.IncludeRegexp[i].String() {
				return false
			}
		}

		for i := range v.ExcludeRegexp {
			if v.ExcludeRegexp[i].String() != vWant.ExcludeRegexp[i].String() {
				return false
			}
		}
	}
	return true
}

func getRegexpArr(expr []string) []*regexp.Regexp {
	ret := make([]*regexp.Regexp, len(expr))
	for i := range expr {
		reg, err := regexp.Compile(expr[i])
		if err != nil {
			return nil
		}
		ret[i] = reg
	}
	return ret
}
