package platform // import "github.com/xtls/xray-core/common/platform"

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	PluginLocation  = "caddy.location.plugin"
	ConfigLocation  = "caddy.location.config"
	ConfdirLocation = "caddy.location.confdir"
	ToolLocation    = "caddy.location.tool"
	AssetLocation   = "caddy.location.asset"
	CertLocation    = "caddy.location.cert"

	UseReadV         = "caddy.buf.readv"
	UseFreedomSplice = "caddy.buf.splice"
	UseVmessPadding  = "caddy.vmess.padding"
	UseCone          = "caddy.cone.disabled"

	BufferSize           = "caddy.ray.buffer.size"
	BrowserDialerAddress = "caddy.browser.dialer"
	XUDPLog              = "caddy.xudp.show"
	XUDPBaseKey          = "caddy.xudp.basekey"
)

type EnvFlag struct {
	Name    string
	AltName string
}

func NewEnvFlag(name string) EnvFlag {
	return EnvFlag{
		Name:    name,
		AltName: NormalizeEnvName(name),
	}
}

func (f EnvFlag) GetValue(defaultValue func() string) string {
	if v, found := os.LookupEnv(f.Name); found {
		return v
	}
	if len(f.AltName) > 0 {
		if v, found := os.LookupEnv(f.AltName); found {
			return v
		}
	}

	return defaultValue()
}

func (f EnvFlag) GetValueAsInt(defaultValue int) int {
	useDefaultValue := false
	s := f.GetValue(func() string {
		useDefaultValue = true
		return ""
	})
	if useDefaultValue {
		return defaultValue
	}
	v, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return defaultValue
	}
	return int(v)
}

func NormalizeEnvName(name string) string {
	return strings.ReplaceAll(strings.ToUpper(strings.TrimSpace(name)), ".", "_")
}

func getExecutableDir() string {
	exec, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(exec)
}

func getExecutableSubDir(dir string) func() string {
	return func() string {
		return filepath.Join(getExecutableDir(), dir)
	}
}

func GetPluginDirectory() string {
	pluginDir := NewEnvFlag(PluginLocation).GetValue(getExecutableSubDir("plugins"))
	return pluginDir
}

func GetConfigurationPath() string {
	configPath := NewEnvFlag(ConfigLocation).GetValue(getExecutableDir)
	return filepath.Join(configPath, "config.json")
}

// GetConfDirPath reads "xray.location.confdir"
func GetConfDirPath() string {
	configPath := NewEnvFlag(ConfdirLocation).GetValue(func() string { return "" })
	return configPath
}
