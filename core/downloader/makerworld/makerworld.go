package makerworld

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/eduardooliveira/stLib/core/data/database"
	"github.com/eduardooliveira/stLib/core/downloader/tools"
	"github.com/eduardooliveira/stLib/core/entities"
	"github.com/eduardooliveira/stLib/core/processing/types"
	"github.com/eduardooliveira/stLib/core/utils"
	"golang.org/x/net/html"
)

type mwClient struct {
	client    *http.Client
	userAgent string
	project   *entities.Project
	metadata  *makerWorldMetaData
}

const failedToFech3MF = "failed downloading 3mf"

func Fetch(urlString string, cookies []*http.Cookie, userAgent string) error {

	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}

	httpClient := &http.Client{
		Jar: jar,
	}

	u, err := url.Parse(urlString)
	if err != nil {
		return err
	}
	httpClient.Jar.SetCookies(u, cookies)

	project := entities.NewProject("CHANGE ME")

	mwc := &mwClient{
		client:    httpClient,
		userAgent: userAgent,
		project:   project,
	}

	metadata, err := mwc.fetchDetails(u)
	if err != nil {
		return err
	}
	mwc.metadata = metadata

	if err = utils.CreateFolder(utils.ToLibPath(project.FullPath())); err != nil {
		log.Println("error creating project folder")
		return err
	}

	if err = utils.CreateAssetsFolder(project.UUID); err != nil {
		log.Println("error creating assets folder")
		return err
	}

	assets := make([]*types.ProcessableAsset, 0)
	as, err := mwc.fetchCover()
	if err != nil {
		log.Println("error fetching cover")
		return err
	}

	for _, a := range as {
		if a.Asset != nil {
			project.DefaultImageID = a.Asset.ID
		}
	}

	assets = append(assets, as...)

	as, err = mwc.fetchModels()
	if err != nil {
		log.Println("error fetching models")
		return err
	}
	assets = append(assets, as...)

	_, err = mwc.fetchInstances()
	if err != nil {
		log.Println("error fetching models")
		return err
	}

	as, err = mwc.fetchPictures()
	if err != nil {
		log.Println("error fetching pictures")
		return err
	}
	assets = append(assets, as...)

	if project.DefaultImageID == "" {
		for _, a := range assets {
			if a.Asset != nil && a.Asset.AssetType == "image" {
				project.DefaultImageID = a.Asset.ID
				break
			}
		}
	}

	return database.InsertProject(project)
}

func (mwc *mwClient) fetchCover() ([]*types.ProcessableAsset, error) {
	req, err := http.NewRequest("GET", mwc.metadata.Props.PageProps.Design.CoverURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", mwc.userAgent)

	return tools.DownloadAsset(path.Base(mwc.metadata.Props.PageProps.Design.CoverURL), mwc.project, mwc.client, req)
}

func (mwc *mwClient) fetchPictures() ([]*types.ProcessableAsset, error) {
	assets := make([]*types.ProcessableAsset, 0)
	for _, p := range mwc.metadata.Props.PageProps.Design.DesignExtension.DesignPictures {
		req, err := http.NewRequest("GET", p.URL, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Add("User-Agent", mwc.userAgent)

		a, err := tools.DownloadAsset(p.Name, mwc.project, mwc.client, req)
		if err != nil {
			log.Println("Error fetchig image, skiping: ", err)
			continue
		}
		assets = append(assets, a...)
	}

	return assets, nil
}

func (mwc *mwClient) fetchModels() ([]*types.ProcessableAsset, error) {
	assets := make([]*types.ProcessableAsset, 0)
	for _, m := range mwc.metadata.Props.PageProps.Design.DesignExtension.ModelFiles {
		req, err := http.NewRequest("GET", m.ModelURL, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Add("User-Agent", mwc.userAgent)

		a, err := tools.DownloadAsset(m.ModelName, mwc.project, mwc.client, req)
		if err != nil {
			log.Println("Error fetchig model, skiping: ", err)
			continue
		}
		assets = append(assets, a...)
	}

	return assets, nil
}

func (mwc *mwClient) fetchInstances() ([]*types.ProcessableAsset, error) {
	assets := make([]*types.ProcessableAsset, 0)

	for _, m := range mwc.metadata.Props.PageProps.Design.Instances {
		sl := rand.Intn(6000-3000) + 3000
		log.Println("sleeping: ", sl)
		time.Sleep(time.Duration(sl) * time.Millisecond)
		mfData, err := mwc.fetch3MFData(m.ID)
		if err != nil {
			log.Println("Failed to download 3MFData, skipping")
			continue
		}
		req, err := http.NewRequest("GET", mfData.URL, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Add("User-Agent", mwc.userAgent)

		ext := filepath.Ext(mfData.Name)
		name := strings.TrimSuffix(mfData.Name, filepath.Ext(mfData.Name))
		name = fmt.Sprintf("%s%d%s", name, m.ID, ext)

		sl = rand.Intn(6000-3000) + 3000
		log.Println("sleeping: ", sl)
		time.Sleep(time.Duration(sl) * time.Millisecond)

		a, err := tools.DownloadAsset(name, mwc.project, mwc.client, req)
		if err != nil {
			log.Println("Failed to download 3MF File, skipping")
			continue
		}
		assets = append(assets, a...)
	}

	return assets, nil
}

func (mwc *mwClient) fetch3MFData(id int) (*mf, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://makerworld.com/api/v1/design-service/instance/%d/f3mf?type=download", id), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", mwc.userAgent)
	qwe, _ := httputil.DumpRequestOut(req, true)
	log.Println(string(qwe))
	resp, err := mwc.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(failedToFech3MF)
	}

	mfData := &mf{}
	if err = json.NewDecoder(resp.Body).Decode(mfData); err != nil {
		return nil, err
	}

	if mfData.Name == "" {
		return nil, errors.New(failedToFech3MF)
	}

	return mfData, nil
}

func (mwc *mwClient) fetchDetails(url *url.URL) (*makerWorldMetaData, error) {
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", mwc.userAgent)

	resp, err := mwc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	metadata, err := parseMakerWorldPage(resp.Body)
	if err != nil {
		return nil, err
	}

	mwc.project.Name = metadata.Props.PageProps.Design.Title
	mwc.project.Description = metadata.Props.PageProps.Design.Summary
	mwc.project.Tags = entities.StringsToTags(metadata.Props.PageProps.Design.Tags)
	for _, c := range metadata.Props.PageProps.Design.Categories {
		mwc.project.Tags = append(mwc.project.Tags, entities.StringToTag(c.Name))
	}
	mwc.project.ExternalLink = url.String()

	return metadata, nil
}

func parseMakerWorldPage(body io.ReadCloser) (*makerWorldMetaData, error) {
	doc, err := html.Parse(body)
	if err != nil {
		fmt.Println(err)
		return nil, err
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
					metaDataStr = n.FirstChild.Data
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			search(c)
		}
	}
	search(doc)

	if metaDataStr == "" {
		return nil, errors.New("metadata not found")
	}

	metaData := &makerWorldMetaData{}
	if err := json.Unmarshal([]byte(metaDataStr), metaData); err != nil {
		return nil, err
	}
	return metaData, nil
}

type mf struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

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
					Extention struct {
						ModelInfo struct {
							AuxiliaryPictures []struct {
								Name string `json:"name"`
								URL  string `json:"url"`
							} `json:"auxiliaryPictures"`
							AuxiliaryBom   []any `json:"auxiliaryBom"`
							AuxiliaryGuide []any `json:"auxiliaryGuide"`
							AuxiliaryOther []any `json:"auxiliaryOther"`
							Compatibility  struct {
								DevModelName   string  `json:"devModelName"`
								DevProductName string  `json:"devProductName"`
								NozzleDiameter float64 `json:"nozzleDiameter"`
							} `json:"compatibility"`
							OtherCompatibility []struct {
								DevModelName   string  `json:"devModelName"`
								DevProductName string  `json:"devProductName"`
								NozzleDiameter float64 `json:"nozzleDiameter"`
							} `json:"otherCompatibility"`
							Plates []struct {
								Index     int    `json:"index"`
								Name      string `json:"name"`
								Thumbnail struct {
									Name string `json:"name"`
									URL  string `json:"url"`
								} `json:"thumbnail"`
								TopPicture struct {
									Name string `json:"name"`
									Dir  string `json:"dir"`
									URL  string `json:"url"`
								} `json:"top_picture"`
								PickPicture struct {
									Name string `json:"name"`
									Dir  string `json:"dir"`
									URL  string `json:"url"`
								} `json:"pick_picture"`
								Objects            any  `json:"objects"`
								SkippedObjects     any  `json:"skipped_objects"`
								LabelObjectEnabled bool `json:"label_object_enabled"`
								Prediction         int  `json:"prediction"`
								Weight             int  `json:"weight"`
								Filaments          []struct {
									ID    string `json:"id"`
									Type  string `json:"type"`
									Color string `json:"color"`
									UsedM string `json:"usedM"`
									UsedG string `json:"usedG"`
								} `json:"filaments"`
								Warning any `json:"warning"`
							} `json:"plates"`
						} `json:"modelInfo"`
						InstanceSetting struct {
							SubmitAsPrivate        bool   `json:"submitAsPrivate"`
							IsPrinterPresetChanged bool   `json:"isPrinterPresetChanged"`
							IsPrinterTested        bool   `json:"isPrinterTested"`
							IsDonateToAuthor       bool   `json:"isDonateToAuthor"`
							AuthorsChoice          bool   `json:"authorsChoice"`
							MakerLab               string `json:"makerLab"`
							MakerLabVersion        string `json:"makerLabVersion"`
						} `json:"instanceSetting"`
					} `json:"extention"`
					Prediction      int `json:"prediction"`
					Weight          int `json:"weight"`
					InstanceCreator struct {
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
					} `json:"instanceCreator"`
					Cover             string `json:"cover"`
					MaterialCnt       int    `json:"materialCnt"`
					IsDefault         bool   `json:"isDefault"`
					MaterialColorCnt  int    `json:"materialColorCnt"`
					NeedAms           bool   `json:"needAms"`
					InstanceFilaments []struct {
						Type  string `json:"type"`
						Color string `json:"color"`
						UsedM string `json:"usedM"`
						UsedG string `json:"usedG"`
					} `json:"instanceFilaments"`
					RatingScoreTotal int     `json:"ratingScoreTotal"`
					RatingCount      int     `json:"ratingCount"`
					Score            float64 `json:"score"`
					DownloadCount    int     `json:"downloadCount"`
					PrintCount       int     `json:"printCount"`
					IsOfficial       bool    `json:"isOfficial"`
					AppCanPrint      bool    `json:"appCanPrint"`
				} `json:"instances"`
				License         string `json:"license"`
				Nsfw            bool   `json:"nsfw"`
				Originals       []any  `json:"originals"`
				ModelSource     int    `json:"modelSource"`
				DesignExtension struct {
					DesignSetting struct {
						AllowOthersProfile bool `json:"allowOthersProfile"`
						SubmitAsPrivate    bool `json:"submitAsPrivate"`
						DesignUpdateMsg    struct {
							NoticeInstanceAuthor bool   `json:"noticeInstanceAuthor"`
							NoticeDownloadUser   bool   `json:"noticeDownloadUser"`
							Content              string `json:"content"`
							Images               any    `json:"images"`
						} `json:"designUpdateMsg"`
					} `json:"design_setting"`
					DesignPictures []struct {
						Name string `json:"name"`
						URL  string `json:"url"`
					} `json:"design_pictures"`
					DesignBom   []any `json:"design_bom"`
					DesignGuide []any `json:"design_guide"`
					DesignOther []any `json:"design_other"`
					ModelFiles  []struct {
						ThumbnailName   string `json:"thumbnailName"`
						ThumbnailSize   int    `json:"thumbnailSize"`
						ThumbnailURL    string `json:"thumbnailUrl"`
						ModelName       string `json:"modelName"`
						ModelSize       int    `json:"modelSize"`
						ModelURL        string `json:"modelUrl"`
						ModelType       string `json:"modelType"`
						Note            string `json:"note"`
						IsDir           bool   `json:"isDir"`
						DirName         string `json:"dirName"`
						IsAutoGenerated bool   `json:"isAutoGenerated"`
						Children        any    `json:"children"`
						ModelFileName   string `json:"modelFileName"`
						Unikey          string `json:"unikey"`
					} `json:"model_files"`
				} `json:"designExtension"`
				Status            int    `json:"status"`
				DefaultInstanceID int    `json:"defaultInstanceId"`
				IsStaffPicked     bool   `json:"isStaffPicked"`
				PickReason        string `json:"pickReason"`
				IsPrintable       bool   `json:"isPrintable"`
				IsOfficial        bool   `json:"isOfficial"`
				IsPointRedeemable bool   `json:"isPointRedeemable"`
				PointRedeemDetail struct {
					Price  int    `json:"price"`
					Sku    string `json:"sku"`
					Status int    `json:"status"`
				} `json:"pointRedeemDetail"`
				IsAlreadyRedeemed bool `json:"isAlreadyRedeemed"`
				IsExclusive       bool `json:"isExclusive"`
				Contest           struct {
					ContestID   int    `json:"contestId"`
					Rank        int    `json:"rank"`
					Status      int    `json:"status"`
					ContestName string `json:"contestName"`
				} `json:"contest"`
			} `json:"design"`
			RemixedDesigns struct {
				Total int   `json:"total"`
				Hits  []any `json:"hits"`
			} `json:"remixedDesigns"`
		} `json:"pageProps"`
		NSsp bool `json:"__N_SSP"`
	} `json:"props"`
}
