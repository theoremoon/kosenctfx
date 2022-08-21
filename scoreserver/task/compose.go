package task

import (
	"context"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/compose-spec/compose-go/loader"
	"github.com/compose-spec/compose-go/types"
	composeTypes "github.com/compose-spec/compose-go/types"
	"github.com/docker/cli/cli/command"
	cliflags "github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/compose/v2/pkg/compose"
	"github.com/pkg/errors"
	"github.com/theoremoon/kosenctfx/scoreserver/task/registry"
)

/// compose.yaml または　docker-compose.ymlを読み込む
/// ファイルシステム上にないyamlを読み込むことがありうるのでパスをごまかせるように
/// ファイルを読む処理は別でやることにしてる
func ParseComposeConfig(path string, buf []byte) (*composeTypes.Project, error) {
	composeConfig, err := loader.ParseYAML(buf)
	if err != nil {
		return nil, errors.Wrap(err, "parse error")
	}

	dir := filepath.Dir(path)
	config, err := loader.Load(composeTypes.ConfigDetails{
		WorkingDir: dir,
		ConfigFiles: []composeTypes.ConfigFile{
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

func ValidateComposeConfig(conf *composeTypes.Project) error {
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

type Compose struct {
	compose api.Service
}

func NewCompose(registryConf *registry.RegistryConfig) (*Compose, error) {
	dockerCli, err := command.NewDockerCli()
	if err != nil {
		return nil, errors.Wrap(err, "docker cli")
	}
	dockerCli.Initialize(cliflags.NewClientOptions())

	// login
	creds := dockerCli.ConfigFile().GetCredentialsStore(registryConf.GetRegistry())
	creds.Store(registryConf.GetAuthConfig())

	return &Compose{
		compose: compose.NewComposeService(dockerCli),
	}, nil
}

func (svc *Compose) Up(ctx context.Context, project *types.Project, port int) error {
	for i := range project.Services {
		// ここで設定されたlabelをもとにdocker apiで検索とかする
		if project.Services[i].CustomLabels == nil {
			project.Services[i].CustomLabels = map[string]string{
				api.ProjectLabel:     project.Name,
				api.ServiceLabel:     project.Services[i].Name,
				api.VersionLabel:     api.ComposeVersion,
				api.WorkingDirLabel:  project.WorkingDir,
				api.ConfigFilesLabel: strings.Join(project.ComposeFiles, ","),
				api.OneoffLabel:      "False", // default, will be overridden by `run` command
			}
		}
		project.Services[i].Build = nil
		if len(project.Services[i].Ports) != 0 {
			project.Services[i].Ports[0].Published = strconv.Itoa(port)
		}
	}

	timeout := 1 * time.Minute
	var consumer api.LogConsumer
	// consumer := formatter.NewLogConsumer(ctx, os.Stdout, false, false)

	err := svc.compose.Up(ctx, project, api.UpOptions{
		Create: api.CreateOptions{
			Services:             project.ServiceNames(),
			RemoveOrphans:        false,
			IgnoreOrphans:        true,
			Recreate:             api.RecreateDiverged,
			RecreateDependencies: api.RecreateDiverged,
			Inherit:              false,
			Timeout:              &timeout,
			QuietPull:            true,
		},
		Start: api.StartOptions{
			Project:      project,
			Attach:       consumer,
			AttachTo:     project.ServiceNames(),
			CascadeStop:  true,
			ExitCodeFrom: "",
			Wait:         true,
		},
	})
	if err != nil {
		return errors.Wrap(err, "up")
	}
	return nil
}

func (svc *Compose) Down(ctx context.Context, id string) error {
	timeout := 1 * time.Minute
	err := svc.compose.Down(ctx, id, api.DownOptions{
		RemoveOrphans: true,
		Project:       nil,
		Timeout:       &timeout,
		Images:        "",
		Volumes:       true,
	})
	if err != nil {
		return errors.Wrap(err, "down")
	}
	return nil
}
