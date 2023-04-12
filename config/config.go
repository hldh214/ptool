package config

import (
	"os"
	"sync"

	"github.com/sagan/ptool/utils"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const (
	DEFAULT_SITE_TIMEZONE                         = "Asia/Shanghai"
	DEFAULT_CLIENT_BRUSH_MIN_DISK_SPACE           = int64(5 * 1024 * 1024 * 1024)
	DEFAULT_CLIENT_BRUSH_SLOW_UPLOAD_SPEED_TIER   = int64(100 * 1024)
	DEFAULT_CLIENT_BRUSH_MAX_DOWNLOADING_TORRENTS = int64(6)
	DEFAULT_CLIENT_BRUSH_MAX_TORRENTS             = int64(500)
	DEFAULT_CLIENT_BRUSH_MIN_RATION               = float64(0.2)
	DEFAULT_SITE_TORRENT_UPLOAD_SPEED_LIMIT       = int64(10 * 1024 * 1024)
	VERSION                                       = "0.0.1"
)

type ClientConfigStruct struct {
	Name                          string  `yaml:"name"`
	Disabled                      bool    `yaml:"disabled"`
	Type                          string  `yaml:"type"`
	Url                           string  `yaml:"url"`
	Username                      string  `yaml:"username"`
	Password                      string  `yaml:"password"`
	BrushMinDiskSpace             string  `yaml:"brushMinDiskSpace"`
	BrushSlowUploadSpeedTier      string  `yaml:"brushSlowUploadSpeedTier"`
	BrushMaxDownloadingTorrents   int64   `yaml:"brushMaxDownloadingTorrents"`
	BrushMaxTorrents              int64   `yaml:"brushMaxTorrents"`
	BrushMinRatio                 float64 `yaml:"brushMinRatio"`
	BrushMinDiskSpaceValue        int64
	BrushSlowUploadSpeedTierValue int64
}

type SiteConfigStruct struct {
	Name                         string `yaml:"name"`
	Disabled                     bool   `yaml:"disabled"`
	Type                         string `yaml:"type"`
	Url                          string `yaml:"url"`
	BrushUrl                     string `yaml:"brushUrl"`
	Cookie                       string `yaml:"cookie"`
	TorrentUploadSpeedLimit      string `yaml:"uploadSpeedLimit"`
	GlobalHnR                    bool   `yaml:"globalHnR"`
	Timezone                     string `yaml:"timezone"`
	TorrentUploadSpeedLimitValue int64
}

type ConfigStruct struct {
	IyuuToken        string               `yaml:"iyuutoken"`
	BrushEnableStats bool                 `yaml:"brushEnableStats"`
	Clients          []ClientConfigStruct `yaml:"clients"`
	Sites            []SiteConfigStruct   `yaml:"sites"`
}

var (
	VerboseLevel               = 0
	ConfigDir                  = ""
	ConfigFile                 = ""
	ConfigLoaded               = false
	Config       *ConfigStruct = &ConfigStruct{}
	mu           sync.Mutex
)

func init() {

}

func Get() *ConfigStruct {
	if !ConfigLoaded {
		mu.Lock()
		if !ConfigLoaded {
			log.Debugf("Read config file %s", ConfigFile)
			file, err := os.ReadFile(ConfigFile)
			if err == nil {
				err = yaml.Unmarshal(file, &Config)
				if err != nil {
					log.Print(err)
				}
			}
			for i, client := range Config.Clients {
				v, err := utils.RAMInBytes(client.BrushMinDiskSpace)
				if err != nil || v <= 0 {
					v = DEFAULT_CLIENT_BRUSH_MIN_DISK_SPACE
				}
				Config.Clients[i].BrushMinDiskSpaceValue = v

				v, err = utils.RAMInBytes(client.BrushSlowUploadSpeedTier)
				if err != nil || v <= 0 {
					v = DEFAULT_CLIENT_BRUSH_SLOW_UPLOAD_SPEED_TIER
				}
				Config.Clients[i].BrushSlowUploadSpeedTierValue = v

				if client.BrushMaxDownloadingTorrents == 0 {
					Config.Clients[i].BrushMaxDownloadingTorrents = DEFAULT_CLIENT_BRUSH_MAX_DOWNLOADING_TORRENTS
				}

				if client.BrushMaxTorrents == 0 {
					Config.Clients[i].BrushMaxTorrents = DEFAULT_CLIENT_BRUSH_MAX_TORRENTS
				}

				if client.BrushMinRatio == 0 {
					Config.Clients[i].BrushMinRatio = DEFAULT_CLIENT_BRUSH_MIN_RATION
				}

				if client.Name == "" {
					Config.Clients[i].Name = client.Type
				}
			}
			for i, site := range Config.Sites {
				v, err := utils.RAMInBytes(site.TorrentUploadSpeedLimit)
				if err != nil || v <= 0 {
					v = DEFAULT_SITE_TORRENT_UPLOAD_SPEED_LIMIT
				}
				Config.Sites[i].TorrentUploadSpeedLimitValue = v

				if site.Name == "" {
					Config.Sites[i].Name = site.Type
				}

				if site.Timezone == "" {
					Config.Sites[i].Timezone = DEFAULT_SITE_TIMEZONE
				}
			}
			ConfigLoaded = true
		}
		Config.Clients = utils.Filter(Config.Clients, func(c ClientConfigStruct) bool {
			return !c.Disabled
		})
		Config.Sites = utils.Filter(Config.Sites, func(s SiteConfigStruct) bool {
			return !s.Disabled
		})
		mu.Unlock()
	}
	return Config
}

func GetClientConfig(name string) *ClientConfigStruct {
	for _, client := range Get().Clients {
		if client.Name == name {
			return &client
		}
	}
	return nil
}

func GetSiteConfig(name string) *SiteConfigStruct {
	for _, site := range Get().Sites {
		if site.Name == name {
			return &site
		}
	}
	return nil
}
