package imagebuilder

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/compose-spec/compose-go/types"
	"github.com/containerd/containerd/platforms"
	"github.com/docker/buildx/build"
	"github.com/docker/buildx/driver"
	_ "github.com/docker/buildx/driver/docker"
	"github.com/docker/buildx/util/buildflags"
	"github.com/docker/buildx/util/progress"
	"github.com/docker/cli/cli/config"
	"github.com/docker/cli/cli/config/configfile"
	"github.com/docker/compose/v2/pkg/api"
	moby "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	bclient "github.com/moby/buildkit/client"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/auth/authprovider"
	"github.com/moby/moby/pkg/jsonmessage"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
	"github.com/theoremoon/kosenctfx/scoreserver/task"
	"github.com/theoremoon/kosenctfx/scoreserver/task/registry"
)

const driverName = "default"

type ImageBuilder interface {
	BuildAndPush(ctx context.Context, t *task.TaskDefinition, compose *types.Project) error
	CleanContainerName(compose *types.Project) error
}

type imageBuilder struct {
	dockerClient *client.Client
	registry     *registry.RegistryConfig
}

func New(registry *registry.RegistryConfig) (ImageBuilder, error) {
	c, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	return &imageBuilder{
		dockerClient: c,
		registry:     registry,
	}, nil
}

// container nameを指定されていると別々に起動できなくて都合が悪い
// docker buildとかするとcontainer_nameを勝手に埋められてしまうので上書きする
func (b *imageBuilder) CleanContainerName(compose *types.Project) error {
	for i, _ := range compose.AllServices() {
		compose.Services[i].ContainerName = ""
	}
	return nil
}

/// build docker images by docker-compose and push it into container registery
/// this changes the `compose` argument
func (b *imageBuilder) BuildAndPush(ctx context.Context, t *task.TaskDefinition, compose *types.Project) error {
	configFile := loadConfigFile()

	// prepare driver
	d, err := driver.GetDriver(ctx, driverName, nil, b.dockerClient, configFile, nil, nil, nil, nil, nil, compose.WorkingDir)
	if err != nil {
		return errors.Wrap(err, "get driver")
	}
	driverInfo := []build.DriverInfo{
		{
			Name:   driverName,
			Driver: d,
		},
	}

	// prepare writer
	progressCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	progressWriter := progress.NewPrinter(progressCtx, os.Stdout, os.Stdout, progress.PrinterModeAuto)

	// ここでdocker-composeを書き換えて、imageの指定を行う
	shouldPushList := make(map[string]struct{})
	for i, service := range compose.Services {
		if compose.Services[i].Image == "" {
			image := b.getFullImageName(t, service)
			compose.Services[i].Image = image
			shouldPushList[image] = struct{}{} // 自分でimage名を指定したヤツはpush対象
		}
	}

	// build
	opts, err := b.makeBuildOptions(t, compose)
	if err != nil {
		return errors.Wrap(err, "make options")
	}
	resp, err := build.Build(ctx, driverInfo, opts, nil, filepath.Dir(configFile.Filename), progressWriter)
	if err != nil {
		return err
	}
	if err := progressWriter.Wait(); err != nil {
		return errors.Wrap(err, "progresswriter")
	}

	// なんかこのあたりは使ってないけど値を入れておかないとnilで困ったりするので
	// デフォルトっぽい値を入れている
	for i, service := range compose.Services {
		if compose.Services[i].CustomLabels == nil {
			compose.Services[i].CustomLabels = types.Labels{}
		}
		digest, ok := resp[service.Image]
		if ok {
			compose.Services[i].CustomLabels[api.ImageDigestLabel] = digest.ExporterResponse["containeraimage.digest"]
		}
	}

	// docker registryの認証情報をセットする
	// registry毎の認証情報をmapにしておいて、imageタグの値から合致する認証情報が使われる
	// 現実的には複数のregistryに分散してpushしたいことはないはずなので指定されている一個だけの認証情報を使う
	authes := make(map[string]string)
	authes[b.registry.GetRegistry()] = b.registry.GetRegistryAuth()

	// don't use system default config
	// for key, aconf := range configFile.GetAuthConfigs() {
	// 	blob, _ := json.Marshal(aconf)
	// 	authes[key] = base64.URLEncoding.EncodeToString(blob)
	// }

	progressWriter = progress.NewPrinter(progressCtx, os.Stdout, os.Stdout, progress.PrinterModeAuto)
	for _, service := range compose.AllServices() {
		// push対象になっているimageだけ処理する
		if _, exist := shouldPushList[service.Image]; !exist {
			continue
		}

		// 認証情報を読み込む
		var auth string
		for k, v := range authes {
			if strings.HasPrefix(service.Image, k) {
				auth = v
				break
			}
		}
		if auth == "" {
			log.Printf("skip [%s] to push\n", service.Image)
			continue
		}

		err = b.pushDockerImage(service, auth)
		if err != nil {
			return err
		}
	}

	return nil
}

// docker build するときの設定を作るヘルパ
// キャッシュを使いまわすかどうかなどの設定を入れる
func (b *imageBuilder) makeBuildOptions(t *task.TaskDefinition, compose *types.Project) (map[string]build.Options, error) {
	opts := make(map[string]build.Options)
	for _, service := range compose.AllServices() {
		if service.Build == nil {
			continue
		}

		imageName := b.getFullImageName(t, service)
		cacheFrom, err := buildflags.ParseCacheEntry(service.Build.CacheFrom)
		if err != nil {
			return nil, err
		}
		cacheTo, _ := buildflags.ParseCacheEntry(service.Build.CacheFrom)
		if err != nil {
			return nil, err
		}

		plats := make([]specs.Platform, 0)
		if service.Platform != "" {
			p, _ := platforms.Parse(service.Platform)
			plats = append(plats, p)
		}

		buildArgs := service.Build.Args.Resolve(func(s string) (string, bool) {
			s, ok := compose.Environment[s]
			return s, ok
		})
		argsMap := make(map[string]string)
		for k, v := range buildArgs {
			if v != nil {
				argsMap[k] = *v
			}
		}

		opts[imageName] = build.Options{
			Inputs: build.Inputs{
				ContextPath:    getContextPath(compose, service),
				DockerfilePath: getDockerfilePath(compose, service),
			},
			CacheFrom:   cacheFrom,
			CacheTo:     cacheTo,
			NoCache:     service.Build.NoCache,
			Pull:        service.Build.Pull,
			BuildArgs:   argsMap,
			Tags:        []string{imageName},
			Target:      service.Build.Target,
			Exports:     []bclient.ExportEntry{{Type: "image", Attrs: map[string]string{}}},
			Platforms:   plats,
			Labels:      service.Build.Labels,
			NetworkMode: service.Build.Network,
			ExtraHosts:  service.Build.ExtraHosts.AsList(),
			Session: []session.Attachable{
				authprovider.NewDockerAuthProvider(os.Stderr),
			},
		}
	}

	return opts, nil
}

func (b *imageBuilder) pushDockerImage(service types.ServiceConfig, auth string) error {
	// do push
	log.Printf("push %s\n", service.Image)
	stream, err := b.dockerClient.ImagePush(context.Background(), service.Image, moby.ImagePushOptions{
		RegistryAuth: auth,
	})
	if err != nil {
		return err
	}

	dec := json.NewDecoder(stream)
	lastStatus := ""

	for {
		var jm jsonmessage.JSONMessage
		if err := dec.Decode(&jm); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		if jm.ID != "" && lastStatus != jm.Progress.String() {
			log.Printf("pushing %s: %s\n", service.Image, jm.Progress.String())
			lastStatus = jm.Progress.String()
		}

		if jm.Error != nil {
			return errors.Wrap(errors.New(jm.Error.Message), "push")
		}
	}
	return nil
}

func loadConfigFile() *configfile.ConfigFile {
	return config.LoadDefaultConfigFile(io.Discard)
}

// image名が決まってなかったら作って返す
func getImageName(t *task.TaskDefinition, service types.ServiceConfig) string {
	if service.Image != "" {
		return service.Image
	}
	return t.ID + "_" + service.Name
}

func (b *imageBuilder) getFullImageName(t *task.TaskDefinition, service types.ServiceConfig) string {
	if service.Image != "" {
		return service.Image
	}

	return filepath.Join(b.registry.GetRegistry(), t.ID+"_"+service.Name)
}

func getContextPath(compose *types.Project, service types.ServiceConfig) string {
	return filepath.Join(compose.WorkingDir, service.Build.Context)
}

func getDockerfilePath(compose *types.Project, service types.ServiceConfig) string {
	return filepath.Join(compose.WorkingDir, service.Build.Context, service.Build.Dockerfile)
}
