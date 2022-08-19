package task

import (
	"path/filepath"

	"github.com/compose-spec/compose-go/loader"
	compose "github.com/compose-spec/compose-go/types"
	"github.com/pkg/errors"
)

/// compose.yaml または　docker-compose.ymlを読み込む
/// ファイルシステム上にないyamlを読み込むことがありうるのでパスをごまかせるように
/// ファイルを読む処理は別でやることにしてる
func ParseComposeConfig(path string, buf []byte) (*compose.Project, error) {
	composeConfig, err := loader.ParseYAML(buf)
	if err != nil {
		return nil, errors.Wrap(err, "parse error")
	}

	dir := filepath.Dir(path)
	config, err := loader.Load(compose.ConfigDetails{
		WorkingDir: dir,
		ConfigFiles: []compose.ConfigFile{
			{
				Filename: path,
				Config:   composeConfig,
			},
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "load err")
	}
	return config, nil
}

func ValidateComposeConfig(conf *compose.Project) error {
	if len(conf.VolumeNames()) != 0 {
		return errors.New("volume is prohibited")
	}
	for _, svc := range conf.Services {
		if len(svc.CapAdd) != 0 {
			return errors.New("cap_add is prohibited")
		}
		if len(svc.Configs) != 0 {
			return errors.New("configs is prohibited")
		}
		if svc.CredentialSpec != nil {
			return errors.New("credential_spec is prohibited")
		}
		if svc.Deploy != nil {
			return errors.New("deploy is prohibited")
		}
		if len(svc.EnvFile) != 0 {
			return errors.New("envfile is prohibited")
		}
		if len(svc.Extends) != 0 {
			return errors.New("extends is prohibited")
		}
		if svc.Logging != nil {
			return errors.New("logging is prohibited")
		}
		if svc.OomKillDisable {
			return errors.New("oom_kill_disable is prohibited")
		}
		if svc.Privileged {
			return errors.New("privileged is prohibited")
		}
		if svc.Runtime != "" {
			return errors.New("runtime is prohibited")
		}
		if len(svc.Secrets) != 0 {
			return errors.New("secrets is prohibited")
		}
		if len(svc.Sysctls) != 0 {
			return errors.New("sysctls is prohibited")
		}
		if len(svc.Tmpfs) != 0 {
			return errors.New("tmpfs is prohibited")
		}
		if len(svc.Volumes) != 0 {
			return errors.New("volumes is prohibited")
		}
		if len(svc.VolumesFrom) != 0 {
			return errors.New("volumes_from is prohibited")
		}
		if len(svc.Extensions) != 0 {
			return errors.New("extensions is prohibited")
		}
	}
	return nil
}
