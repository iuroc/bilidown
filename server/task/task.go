package task

import (
	"bilidown/bilibili"
	"bilidown/util"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

// TaskInitOption 创建任务时的初始数据
type TaskInitOption struct {
	Bvid   string               `json:"bvid"`
	Cid    int                  `json:"cid"`
	Format bilibili.MediaFormat `json:"format"`
	Title  string               `json:"title"`
	Owner  string               `json:"owner"`
	Cover  string               `json:"cover"`
	Status TaskStatus           `json:"status"`
	Folder string               `json:"folder"`
}

// done | waiting | running | error
type TaskStatus string

type Task struct {
	TaskInitOption
	ID            int64   `json:"id"`
	AudioProgress float64 `json:"audioProgress"`
	VideoProgress float64 `json:"videoProgress"`
}

var GlobalTaskList []Task
var GlobalTaskMux = &sync.Mutex{}
var GlobalDownloadSem = util.NewSemaphore(3)

func (task *Task) Create(db *sql.DB) error {
	result, err := db.Exec(`INSERT INTO "task" ("bvid", "cid", "format", "title", "owner", "cover", "status", "folder")
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, task.Bvid, task.Cid, task.Format, task.Title, task.Owner, task.Cover, task.Status, task.Folder)
	if err != nil {
		return err
	}

	task.ID, err = result.LastInsertId()
	return err
}

// Create 创建任务，并将任务加入全局任务列表
func (task *Task) Start(db *sql.DB) error {
	defer db.Close()
	GlobalTaskMux.Lock()
	GlobalTaskList = append(GlobalTaskList, *task)
	GlobalTaskMux.Unlock()
	sessdata, err := bilibili.GetSessdata(db)
	if err != nil {
		fmt.Println(err.Error())
		task.UpdateStatus(db, "error")
		return nil
	}
	client := &bilibili.BiliClient{SESSDATA: sessdata}

	playInfo, err := client.GetPlayInfo(task.Bvid, task.Cid)
	if err != nil {
		task.UpdateStatus(db, "error")
		return nil
	}

	videoURL := ""
	audioURL := ""

	for _, item := range playInfo.Dash.Video {
		if item.ID == task.Format {
			videoURL = item.BaseURL
			break
		}
	}

	var maxAudioID bilibili.MediaFormat
	for _, item := range playInfo.Dash.Audio {
		if item.ID > maxAudioID {
			maxAudioID = item.ID
			audioURL = item.BaseURL
		}
	}
	GlobalDownloadSem.Acquire()
	task.UpdateStatus(db, "running")
	err = DownloadMedia(client, audioURL, task, "audio")
	if err != nil {
		task.UpdateStatus(db, "error")
		return nil
	}
	err = DownloadMedia(client, videoURL, task, "video")
	if err != nil {
		task.UpdateStatus(db, "error")
		return nil
	}
	GlobalDownloadSem.Release()
	task.UpdateStatus(db, "done")
	fmt.Printf("任务 %d 完成\n", task.ID)
	return nil
}

func (task *Task) UpdateStatus(db *sql.DB, status TaskStatus) error {
	_, err := db.Exec(`UPDATE "task" SET "status" = ? WHERE "id" = ?`, status, task.ID)
	return err
}

func GetSize(client *bilibili.BiliClient, _url string) (int64, error) {
	client2 := &http.Client{}
	req, err := http.NewRequest("HEAD", _url, nil)
	if err != nil {
		return 0, err
	}
	req.Header = client.MakeHeader()

	resp, err := client2.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	fmt.Println(_url, resp.ContentLength)
	client2.CloseIdleConnections()
	return resp.ContentLength, nil
}

func DownloadMedia(client *bilibili.BiliClient, _url string, task *Task, mediaType string) error {
	fmt.Printf("开始下载 %d_%s\n", task.ID, mediaType)
	resp, err := client.SimpleGET(_url, nil)
	if err != nil {
		return err
	}

	filename := strconv.FormatInt(task.ID, 10) + "." + mediaType
	filepath := filepath.Join(task.Folder, filename)

	progress := newProgressBar(resp.ContentLength)

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := io.TeeReader(resp.Body, file)
	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		progress.add(n)
		GlobalTaskMux.Lock()
		if mediaType == "video" {
			task.VideoProgress = progress.percent()
		} else {
			task.AudioProgress = progress.percent()
		}
		GlobalTaskMux.Unlock()
	}
	return nil
}

type progressBar struct {
	total   int64
	current int64
}

func (p *progressBar) add(n int) {
	p.current += int64(n)
}

func (p *progressBar) percent() float64 {
	return float64(p.current) / float64(p.total)
}

func newProgressBar(total int64) *progressBar {
	return &progressBar{
		total: total,
	}
}

func GetCurrentFolder(db *sql.DB) (string, error) {
	var filepath string
	err := db.QueryRow(`SELECT "value" FROM "field" WHERE "name" = "download_folder"`).Scan(&filepath)
	if err != nil {
		return "", err
	}
	return filepath, nil
}

// SaveDownloadFolder 保存下载路径，不存在则自动创建
func SaveDownloadFolder(db *sql.DB, downloadFolder string) error {
	_, err := os.Stat(downloadFolder)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(downloadFolder, os.ModePerm)
			if err != nil {
				return err
			}
		}
		return err
	}
	_, err = db.Exec(`INSERT OR REPLACE INTO "field" ("name", "value") VALUES ("download_folder", ?)`, downloadFolder)
	return err
}
