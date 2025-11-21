package makerworld

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/eduardooliveira/stLib/v2/library/downloader/tools"
	"github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/library/libfs"
	"github.com/eduardooliveira/stLib/v2/library/process"
	"github.com/eduardooliveira/stLib/v2/library/repo"
	"github.com/eduardooliveira/stLib/v2/utils"
	"golang.org/x/net/html"
)

type MakerWorldDownloader struct {
	ctx        context.Context
	l          *slog.Logger
	r          *repo.AssetRepo
	p          *process.Processor
	client     *http.Client
	userAgent  string
	asset      *entities.Asset
	metadata   *makerWorldMetaData
	fSys       libfs.LibFS
}

const failedToFetch3MF = "failed downloading 3mf"

func New(ctx context.Context, r *repo.AssetRepo, p *process.Processor, cookies []*http.Cookie, userAgent string) (*MakerWorldDownloader, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{
		Jar: jar,
	}

	// Set cookies if provided
	if len(cookies) > 0 {
		u, _ := url.Parse("https://makerworld.com")
		httpClient.Jar.SetCookies(u, cookies)
	}

	return &MakerWorldDownloader{
		ctx:       ctx,
		l:         slog.With("module", "makerworld"),
		r:         r,
		p:         p,
		client:    httpClient,
		userAgent: userAgent,
		fSys:      libfs.GetDefaultFS(),
	}, nil
}

func (mw *MakerWorldDownloader) Fetch(urlString string, parent entities.Asset) error {
	u, err := url.Parse(urlString)
	if err != nil {
		return fmt.Errorf("parsing URL: %w", err)
	}

	mw.l.Info("Fetching MakerWorld design", "url", urlString)

	// Fetch design details
	metadata, err := mw.fetchDetails(u)
	if err != nil {
		return fmt.Errorf("fetching details: %w", err)
	}
	mw.metadata = metadata

	// Create directory for the design
	designName := strings.ReplaceAll(metadata.Props.PageProps.Design.Title, "/", "-")
	path := filepath.Join(*parent.Path, designName)

	if err := mw.fSys.Mkdir(path); err != nil {
		return fmt.Errorf("creating folder: %w", err)
	}

	// Create asset for the design
	mw.asset = entities.NewAsset(mw.fSys.(entities.LibFS), path, true, &parent)
	mw.asset.Label = utils.Ptr(metadata.Props.PageProps.Design.Title)
	mw.asset.Description = utils.Ptr(metadata.Props.PageProps.Design.Summary)

	// Add tags
	for _, tag := range metadata.Props.PageProps.Design.Tags {
		mw.asset.Tags = append(mw.asset.Tags, entities.StringToTag(tag))
	}
	for _, cat := range metadata.Props.PageProps.Design.Categories {
		mw.asset.Tags = append(mw.asset.Tags, entities.StringToTag(cat.Name))
	}

	// Add external link
	mw.asset.Properties = entities.Properties{
		"externalLink": urlString,
	}

	// Save asset
	if err := mw.r.SaveAsset(*mw.asset); err != nil {
		return fmt.Errorf("saving asset: %w", err)
	}

	// Download cover image
	if err := mw.fetchCover(); err != nil {
		mw.l.Error("Failed to fetch cover", "error", err)
	}

	// Download model files
	if err := mw.fetchModels(); err != nil {
		mw.l.Error("Failed to fetch models", "error", err)
	}

	// Download instances (3MF files)
	if err := mw.fetchInstances(); err != nil {
		mw.l.Error("Failed to fetch instances", "error", err)
	}

	// Download pictures
	if err := mw.fetchPictures(); err != nil {
		mw.l.Error("Failed to fetch pictures", "error", err)
	}

	mw.l.Info("MakerWorld design downloaded successfully", "design", designName)
	return nil
}

func (mw *MakerWorldDownloader) fetchCover() error {
	coverURL := mw.metadata.Props.PageProps.Design.CoverURL
	if coverURL == "" {
		return errors.New("no cover URL found")
	}

	req, err := http.NewRequest("GET", coverURL, nil)
	if err != nil {
		return err
	}
	req.Header.Add("User-Agent", mw.userAgent)

	filename := path.Base(coverURL)
	if err := tools.DownloadFile(mw.fSys, filename, *mw.asset, mw.client, req); err != nil {
		return err
	}

	return mw.processFile(filename)
}

func (mw *MakerWorldDownloader) fetchPictures() error {
	pictures := mw.metadata.Props.PageProps.Design.DesignExtension.DesignPictures
	mw.l.Info("Downloading pictures", "count", len(pictures))

	for _, pic := range pictures {
		req, err := http.NewRequest("GET", pic.URL, nil)
		if err != nil {
			mw.l.Warn("Failed to create request for picture", "error", err)
			continue
		}
		req.Header.Add("User-Agent", mw.userAgent)

		filename := pic.Name
		if filename == "" {
			filename = path.Base(pic.URL)
		}

		if err := tools.DownloadFile(mw.fSys, filename, *mw.asset, mw.client, req); err != nil {
			mw.l.Warn("Failed to download picture", "filename", filename, "error", err)
			continue
		}

		if err := mw.processFile(filename); err != nil {
			mw.l.Warn("Failed to process picture", "filename", filename, "error", err)
		}
	}

	return nil
}

func (mw *MakerWorldDownloader) fetchModels() error {
	models := mw.metadata.Props.PageProps.Design.DesignExtension.ModelFiles
	mw.l.Info("Downloading models", "count", len(models))

	for _, model := range models {
		req, err := http.NewRequest("GET", model.ModelURL, nil)
		if err != nil {
			mw.l.Warn("Failed to create request for model", "error", err)
			continue
		}
		req.Header.Add("User-Agent", mw.userAgent)

		filename := model.ModelName
		if filename == "" {
			filename = path.Base(model.ModelURL)
		}

		if err := tools.DownloadFile(mw.fSys, filename, *mw.asset, mw.client, req); err != nil {
			mw.l.Warn("Failed to download model", "filename", filename, "error", err)
			continue
		}

		if err := mw.processFile(filename); err != nil {
			mw.l.Warn("Failed to process model", "filename", filename, "error", err)
		}
	}

	return nil
}

func (mw *MakerWorldDownloader) fetchInstances() error {
	instances := mw.metadata.Props.PageProps.Design.Instances
	mw.l.Info("Downloading instances (3MF)", "count", len(instances))

	for _, instance := range instances {
		// Add random delay to avoid rate limiting (3-6 seconds)
		sleepDuration := time.Duration(rand.Intn(3000)+3000) * time.Millisecond
		mw.l.Debug("Rate limit delay", "duration", sleepDuration)
		time.Sleep(sleepDuration)

		// Fetch 3MF data
		mfData, err := mw.fetch3MFData(instance.ID)
		if err != nil {
			mw.l.Warn("Failed to fetch 3MF data", "instanceID", instance.ID, "error", err)
			continue
		}

		// Create unique filename
		ext := filepath.Ext(mfData.Name)
		name := strings.TrimSuffix(mfData.Name, ext)
		filename := fmt.Sprintf("%s-%d%s", name, instance.ID, ext)

		// Add another delay before downloading
		sleepDuration = time.Duration(rand.Intn(3000)+3000) * time.Millisecond
		time.Sleep(sleepDuration)

		// Download the file
		req, err := http.NewRequest("GET", mfData.URL, nil)
		if err != nil {
			mw.l.Warn("Failed to create request for 3MF", "error", err)
			continue
		}
		req.Header.Add("User-Agent", mw.userAgent)

		if err := tools.DownloadFile(mw.fSys, filename, *mw.asset, mw.client, req); err != nil {
			mw.l.Warn("Failed to download 3MF", "filename", filename, "error", err)
			continue
		}

		if err := mw.processFile(filename); err != nil {
			mw.l.Warn("Failed to process 3MF", "filename", filename, "error", err)
		}
	}

	return nil
}

func (mw *MakerWorldDownloader) fetch3MFData(id int) (*mf, error) {
	apiURL := fmt.Sprintf("https://makerworld.com/api/v1/design-service/instance/%d/f3mf?type=download", id)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", mw.userAgent)

	resp, err := mw.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: status %d", failedToFetch3MF, resp.StatusCode)
	}

	mfData := &mf{}
	if err := json.NewDecoder(resp.Body).Decode(mfData); err != nil {
		return nil, err
	}

	if mfData.Name == "" || mfData.URL == "" {
		return nil, errors.New(failedToFetch3MF + ": empty response")
	}

	return mfData, nil
}

func (mw *MakerWorldDownloader) fetchDetails(url *url.URL) (*makerWorldMetaData, error) {
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", mw.userAgent)

	resp, err := mw.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch page: status %d", resp.StatusCode)
	}

	metadata, err := parseMakerWorldPage(resp.Body)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

func (mw *MakerWorldDownloader) processFile(name string) error {
	asset := entities.NewAsset(mw.fSys.(entities.LibFS), filepath.Join(*mw.asset.Path, name), false, mw.asset)

	if err := mw.r.SaveAsset(*asset); err != nil {
		return fmt.Errorf("saving asset: %w", err)
	}

	return mw.p.Process(mw.ctx, asset).Wait()
}

func parseMakerWorldPage(body io.ReadCloser) (*makerWorldMetaData, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return nil, fmt.Errorf("parsing HTML: %w", err)
	}

	var metaDataStr string
	var search func(n *html.Node)
	search = func(n *html.Node) {
		if metaDataStr != "" {
			return
		}
		if n.Type == html.ElementNode && n.Data == "script" {
			for _, a := range n.Attr {
				if a.Key == "id" && a.Val == "__NEXT_DATA__" {
					if n.FirstChild != nil {
						metaDataStr = n.FirstChild.Data
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			search(c)
		}
	}
	search(doc)

	if metaDataStr == "" {
		return nil, errors.New("metadata not found in page")
	}

	metaData := &makerWorldMetaData{}
	if err := json.Unmarshal([]byte(metaDataStr), metaData); err != nil {
		return nil, fmt.Errorf("parsing metadata JSON: %w", err)
	}

	return metaData, nil
}

// mf represents 3MF file metadata
type mf struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// makerWorldMetaData represents the page metadata from MakerWorld
type makerWorldMetaData struct {
	Props struct {
		PageProps struct {
			Design struct {
				ID                        int      `json:"id"`
				Title                     string   `json:"title"`
				CoverURL                  string   `json:"coverUrl"`
				Summary                   string   `json:"summary"`
				LikeCount                 int      `json:"likeCount"`
				CollectionCount           int      `json:"collectionCount"`
				ShareCount                int      `json:"shareCount"`
				PrintCount                int      `json:"printCount"`
				CommentCount              int      `json:"commentCount"`
				DownloadCount             int      `json:"downloadCount"`
				RawModelFileDownloadCount int      `json:"rawModelFileDownloadCount"`
				ReadCount                 int      `json:"readCount"`
				Tags                      []string `json:"tags"`
				DesignCreator             struct {
					UID          int64  `json:"uid"`
					Name         string `json:"name"`
					Avatar       string `json:"avatar"`
					FanCount     int    `json:"fanCount"`
					FollowCount  int    `json:"followCount"`
					IsFollowed   bool   `json:"isFollowed"`
					Certificated bool   `json:"certificated"`
					Handle       string `json:"handle"`
					Level        int    `json:"level"`
					GradeType    int    `json:"gradeType"`
				} `json:"designCreator"`
				Categories []struct {
					ID      int    `json:"id"`
					Name    string `json:"name"`
					PicURL  string `json:"picUrl"`
					Desc    string `json:"desc"`
					DescPic string `json:"descPic"`
				} `json:"categories"`
				ModelID    string `json:"modelId"`
				HasLike    bool   `json:"hasLike"`
				HasDislike bool   `json:"hasDislike"`
				HasCollect bool   `json:"hasCollect"`
				Instances  []struct {
					ID        int    `json:"id"`
					ProfileID int    `json:"profileId"`
					Status    int    `json:"status"`
					Title     string `json:"title"`
					Summary   string `json:"summary"`
				} `json:"instances"`
				License         string `json:"license"`
				Nsfw            bool   `json:"nsfw"`
				DesignExtension struct {
					DesignSetting struct {
						AllowOthersProfile bool `json:"allowOthersProfile"`
						SubmitAsPrivate    bool `json:"submitAsPrivate"`
					} `json:"design_setting"`
					DesignPictures []struct {
						Name string `json:"name"`
						URL  string `json:"url"`
					} `json:"design_pictures"`
					DesignBom   []interface{} `json:"design_bom"`
					DesignGuide []interface{} `json:"design_guide"`
					DesignOther []interface{} `json:"design_other"`
					ModelFiles  []struct {
						ThumbnailName   string      `json:"thumbnailName"`
						ThumbnailSize   int         `json:"thumbnailSize"`
						ThumbnailURL    string      `json:"thumbnailUrl"`
						ModelName       string      `json:"modelName"`
						ModelSize       int         `json:"modelSize"`
						ModelURL        string      `json:"modelUrl"`
						ModelType       string      `json:"modelType"`
						Note            string      `json:"note"`
						IsDir           bool        `json:"isDir"`
						DirName         string      `json:"dirName"`
						IsAutoGenerated bool        `json:"isAutoGenerated"`
						Children        interface{} `json:"children"`
						ModelFileName   string      `json:"modelFileName"`
						Unikey          string      `json:"unikey"`
					} `json:"model_files"`
				} `json:"designExtension"`
				Status            int  `json:"status"`
				DefaultInstanceID int  `json:"defaultInstanceId"`
				IsStaffPicked     bool `json:"isStaffPicked"`
				PickReason        string `json:"pickReason"`
				IsPrintable       bool   `json:"isPrintable"`
				IsOfficial        bool   `json:"isOfficial"`
			} `json:"design"`
		} `json:"pageProps"`
		NSsp bool `json:"__N_SSP"`
	} `json:"props"`
}
