package thingiverse

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"time"

	"github.com/eduardooliveira/stLib/v2/config"
	"github.com/eduardooliveira/stLib/v2/library/downloader/tools"
	"github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/library/libfs"
	"github.com/eduardooliveira/stLib/v2/library/process"
	"github.com/eduardooliveira/stLib/v2/library/repo"
	"github.com/eduardooliveira/stLib/v2/utils"
	"golang.org/x/sync/errgroup"
)

type ThingyDownloader struct {
	l     *slog.Logger
	token string
	r     *repo.AssetRepo
	p     *process.Processor
	asset *entities.Asset
	thing *Thing
	eg    *errgroup.Group
	http  *http.Client
	fSys  libfs.LibFS
}

var matcher = regexp.MustCompile(`thing:(\d+)`)

func New(r *repo.AssetRepo, p *process.Processor) (*ThingyDownloader, error) {
	rtn := &ThingyDownloader{
		l:     slog.With("module", "thingiverse"),
		token: config.Cfg.Integrations.Thingiverse.Token,
		r:     r,
		p:     p,
		eg:    &errgroup.Group{},
		http:  &http.Client{},
		thing: &Thing{},
		fSys:  libfs.GetDefaultFS(),
	}
	if rtn.token == "" {
		return nil, errors.New("thingiverse config not set")
	}
	return rtn, nil
}

func (t *ThingyDownloader) Fetch(url string, parent entities.Asset) error {

	matches := matcher.FindStringSubmatch(url)

	if len(matches) == 0 {
		return errors.New("url doesn't match thingiverse schema")
	}

	id := matches[1]
	t.l.Info("Processing", "thing", id)

	httpClient := &http.Client{}

	err := t.fetchThing(id, httpClient)
	if err != nil {
		return fmt.Errorf("fetching %v: %v", id, err)
	}
	path := filepath.Join(*parent.Path, fmt.Sprintf("%v-%s", t.thing.ID, t.thing.Name))
	err = t.fSys.Mkdir(path)
	if err != nil {
		return fmt.Errorf("creating folder: %v", err)
	}

	t.asset = entities.NewAsset(t.fSys, path, true, &parent)
	t.asset.Label = utils.Ptr(t.thing.Name)
	t.asset.Description = utils.Ptr(t.thing.Description)

	for _, tag := range t.thing.Tags {
		t.asset.Tags = append(t.asset.Tags, entities.StringToTag(tag.Name))
	}

	if err := t.r.SaveAsset(*t.asset); err != nil {
		return fmt.Errorf("creating asset: %v", err)
	}

	err = t.fetchFiles()
	if err != nil {
		return fmt.Errorf("fetching files: %v", err)
	}

	err = t.fetchImages()
	if err != nil {
		return fmt.Errorf("fetching files: %v", err)
	}
	return t.eg.Wait()
}

func (t ThingyDownloader) fetchThing(id string, httpClient *http.Client) error {
	u := &url.URL{Scheme: "https", Host: "api.thingiverse.com", Path: "/things/" + id}

	req := &http.Request{
		Method: "GET",
		URL:    u,
		Header: http.Header{
			"Authorization": []string{"Bearer " + t.token},
		},
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(t.thing); err != nil {
		return err
	}

	return nil
}

func (t ThingyDownloader) fetchFiles() error {
	req := &http.Request{
		Method: "GET",
		URL:    &url.URL{Scheme: "https", Host: "api.thingiverse.com", Path: fmt.Sprintf("/things/%d/files", t.thing.ID)},
		Header: http.Header{
			"Authorization": []string{"Bearer " + t.token},
		},
	}
	res, err := t.http.Do(req)
	if err != nil {
		return fmt.Errorf("fetching files: %v", err)
	}
	defer res.Body.Close()

	var files []*ThingFile
	if err := json.NewDecoder(res.Body).Decode(&files); err != nil {
		return fmt.Errorf("decoding files: %v", err)
	}

	req.Method = "GET"

	for _, file := range files {
		lReq := req.Clone(context.Background())
		lReq.URL, _ = url.Parse(file.DownloadURL)
		t.eg.Go(func() error {
			err := tools.DownloadFile(t.fSys, file.Name, *t.asset, t.http, lReq)
			if err != nil {
				return fmt.Errorf("downloading %v: %v", file.Name, err)
			}
			return t.processFile(file.Name)
		})
	}

	t.l.Info("Downloading files", "count", len(files))
	return nil
}

func (t ThingyDownloader) fetchImages() error {
	req := &http.Request{
		Method: "GET",
		URL:    &url.URL{Scheme: "https", Host: "api.thingiverse.com", Path: fmt.Sprintf("/things/%d/images", t.thing.ID)},
		Header: http.Header{
			"Authorization": []string{"Bearer " + t.token},
		},
	}
	res, err := t.http.Do(req)
	if err != nil {
		return fmt.Errorf("fetching files: %v", err)
	}
	defer res.Body.Close()

	var tImages []*ThingImage
	if err := json.NewDecoder(res.Body).Decode(&tImages); err != nil {
		return err
	}

	req.Method = "GET"

	for _, image := range tImages {
		for _, size := range image.Sizes {
			if size.Size == "large" && size.Type == "display" {
				lReq := req.Clone(context.Background())
				lReq.URL, _ = url.Parse(size.URL)
				t.eg.Go(func() error {
					err := tools.DownloadFile(t.fSys, image.Name, *t.asset, t.http, lReq)
					if err != nil {
						return fmt.Errorf("downloading %v: %v", image.Name, err)
					}

					return t.processFile(image.Name)
				})

			}
		}
	}

	t.l.Info("Downloading images", "count", len(tImages))
	return nil
}

func (t ThingyDownloader) processFile(name string) error {
	i := entities.NewAsset(t.fSys, filepath.Join(*t.asset.Path, name), false, t.asset)
	if err := t.r.SaveAsset(*i); err != nil {
		return fmt.Errorf("saving asset: %v", err)
	}
	return t.p.Process(i).Wait()
}

type ThingImage struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	URL   string `json:"url"`
	Sizes []struct {
		Type string `json:"type"`
		Size string `json:"size"`
		URL  string `json:"url"`
	} `json:"sizes"`
}

type ThingFile struct {
	ID            int           `json:"id"`
	Name          string        `json:"name"`
	Size          int           `json:"size"`
	URL           string        `json:"url"`
	PublicURL     string        `json:"public_url"`
	DownloadURL   string        `json:"download_url"`
	ThreejsURL    string        `json:"threejs_url"`
	Thumbnail     string        `json:"thumbnail"`
	DefaultImage  interface{}   `json:"default_image"`
	Date          string        `json:"date"`
	FormattedSize string        `json:"formatted_size"`
	MetaData      []interface{} `json:"meta_data"`
	DownloadCount int           `json:"download_count"`
	DirectURL     string        `json:"direct_url"`
}

type Thing struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Thumbnail string `json:"thumbnail"`
	URL       string `json:"url"`
	PublicURL string `json:"public_url"`
	Creator   struct {
		ID               int    `json:"id"`
		Name             string `json:"name"`
		FirstName        string `json:"first_name"`
		LastName         string `json:"last_name"`
		URL              string `json:"url"`
		PublicURL        string `json:"public_url"`
		Thumbnail        string `json:"thumbnail"`
		CountOfFollowers int    `json:"count_of_followers"`
		CountOfFollowing int    `json:"count_of_following"`
		CountOfDesigns   int    `json:"count_of_designs"`
		AcceptsTips      bool   `json:"accepts_tips"`
		IsFollowing      bool   `json:"is_following"`
		Location         string `json:"location"`
		Cover            string `json:"cover"`
	} `json:"creator"`
	Added        time.Time   `json:"added"`
	Modified     time.Time   `json:"modified"`
	IsPublished  int         `json:"is_published"`
	IsWip        int         `json:"is_wip"`
	IsFeatured   interface{} `json:"is_featured"`
	IsNsfw       bool        `json:"is_nsfw"`
	LikeCount    int         `json:"like_count"`
	IsLiked      bool        `json:"is_liked"`
	CollectCount int         `json:"collect_count"`
	IsCollected  bool        `json:"is_collected"`
	CommentCount int         `json:"comment_count"`
	IsWatched    bool        `json:"is_watched"`
	DefaultImage struct {
		ID    int    `json:"id"`
		URL   string `json:"url"`
		Name  string `json:"name"`
		Sizes []struct {
			Type string `json:"type"`
			Size string `json:"size"`
			URL  string `json:"url"`
		} `json:"sizes"`
		Added time.Time `json:"added"`
	} `json:"default_image"`
	Description      string `json:"description"`
	Instructions     string `json:"instructions"`
	DescriptionHTML  string `json:"description_html"`
	InstructionsHTML string `json:"instructions_html"`
	Details          string `json:"details"`
	DetailsParts     []struct {
		Type     string `json:"type"`
		Name     string `json:"name"`
		Required string `json:"required,omitempty"`
		Data     []struct {
			Content string `json:"content"`
		} `json:"data,omitempty"`
	} `json:"details_parts"`
	EduDetails        string      `json:"edu_details"`
	EduDetailsParts   interface{} `json:"edu_details_parts"`
	License           string      `json:"license"`
	AllowsDerivatives bool        `json:"allows_derivatives"`
	FilesURL          string      `json:"files_url"`
	ImagesURL         string      `json:"images_url"`
	LikesURL          string      `json:"likes_url"`
	AncestorsURL      string      `json:"ancestors_url"`
	DerivativesURL    string      `json:"derivatives_url"`
	TagsURL           string      `json:"tags_url"`
	Tags              []struct {
		Name        string `json:"name"`
		Tag         string `json:"tag"`
		URL         string `json:"url"`
		Count       int    `json:"count"`
		ThingsURL   string `json:"things_url"`
		AbsoluteURL string `json:"absolute_url"`
	} `json:"tags"`
	CategoriesURL     string      `json:"categories_url"`
	FileCount         int         `json:"file_count"`
	LayoutCount       int         `json:"layout_count"`
	LayoutsURL        string      `json:"layouts_url"`
	IsPrivate         int         `json:"is_private"`
	IsPurchased       int         `json:"is_purchased"`
	InLibrary         bool        `json:"in_library"`
	PrintHistoryCount int         `json:"print_history_count"`
	AppID             interface{} `json:"app_id"`
	DownloadCount     int         `json:"download_count"`
	ViewCount         int         `json:"view_count"`
	Education         struct {
		Grades   []interface{} `json:"grades"`
		Subjects []interface{} `json:"subjects"`
	} `json:"education"`
	RemixCount       int           `json:"remix_count"`
	MakeCount        int           `json:"make_count"`
	AppCount         int           `json:"app_count"`
	RootCommentCount int           `json:"root_comment_count"`
	Moderation       string        `json:"moderation"`
	IsDerivative     bool          `json:"is_derivative"`
	Ancestors        []interface{} `json:"ancestors"`
	CanComment       bool          `json:"can_comment"`
	TypeName         string        `json:"type_name"`
}
