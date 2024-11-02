package task

import (
	"bilidown/bilibili"
	"bilidown/util"
	"database/sql"
	"fmt"
	"io"
	"os"
	"os/exec"
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
var GlobalMergeMux = &sync.Mutex{}

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

	videoURL := GetVideoURL(playInfo.Dash.Video, task.Format)
	audioURL := GetAudioURL(playInfo.Dash.Audio)

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

	outputPath := filepath.Join(task.Folder, fmt.Sprintf("%s_%s.mp4", task.Title, util.RandomString(6)))

	videoPath := filepath.Join(task.Folder, strconv.FormatInt(task.ID, 10)+".video")
	audioPath := filepath.Join(task.Folder, strconv.FormatInt(task.ID, 10)+".audio")
	GlobalMergeMux.Lock()
	err = MergeMedia(outputPath, videoPath, audioPath)
	GlobalMergeMux.Unlock()
	if err != nil {
		task.UpdateStatus(db, "error")
		fmt.Println(err.Error())
		return nil
	}
	task.UpdateStatus(db, "done")
	return nil
}

// 合并音视频
func MergeMedia(outputPath string, inputPaths ...string) error {
	inputs := []string{}
	for _, path := range inputPaths {
		inputs = append(inputs, "-i", path)
	}
	err := exec.Command("ffmpeg", append(inputs, "-c:v", "copy", "-c:a", "copy", outputPath)...).Run()
	if err != nil {
		return err
	}
	return nil
}

func GetVideoURL(medias []bilibili.Media, format bilibili.MediaFormat) string {
	for _, code := range []int{12, 7, 13} {
		for _, item := range medias {
			if item.ID == format && item.Codecid == code {
				return item.BaseURL
			}
		}
	}
	return ""
}

func GetAudioURL(medias []bilibili.Media) string {
	var maxAudioID bilibili.MediaFormat
	var audioURL string
	for _, item := range medias {
		if item.ID > maxAudioID {
			maxAudioID = item.ID
			audioURL = item.BaseURL
		}
	}
	return audioURL
}

func (task *Task) UpdateStatus(db *sql.DB, status TaskStatus) error {
	_, err := db.Exec(`UPDATE "task" SET "status" = ? WHERE "id" = ?`, status, task.ID)
	return err
}

func DownloadMedia(client *bilibili.BiliClient, _url string, task *Task, mediaType string) error {
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
