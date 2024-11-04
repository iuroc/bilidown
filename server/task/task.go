package task

import (
	"bilidown/bilibili"
	"bilidown/common"
	"bilidown/util"
	"bufio"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// TaskInitOption 创建任务时需要从 POST 请求获取的参数
type TaskInitOption struct {
	Bvid     string             `json:"bvid"`
	Cid      int                `json:"cid"`
	Format   common.MediaFormat `json:"format"`
	Title    string             `json:"title"`
	Owner    string             `json:"owner"`
	Cover    string             `json:"cover"`
	Status   TaskStatus         `json:"status"`
	Folder   string             `json:"folder"`
	Audio    string             `json:"audio"`
	Video    string             `json:"video"`
	Duration int                `json:"duration"`
}

// TaskInDB 任务数据库中的数据
type TaskInDB struct {
	TaskInitOption
	ID       int64     `json:"id"`
	CreateAt time.Time `json:"createAt"`
}

// done | waiting | running | error
type TaskStatus string

type Task struct {
	TaskInDB
	AudioProgress float64 `json:"audioProgress"`
	VideoProgress float64 `json:"videoProgress"`
	MergeProgress float64 `json:"mergeProgress"`
}

var GlobalTaskList = []*Task{}
var GlobalTaskMux = &sync.Mutex{}
var GlobalDownloadSem = util.NewSemaphore(3)
var GlobalMergeSem = util.NewSemaphore(3)

func (task *Task) Create(db *sql.DB) error {
	result, err := db.Exec(`INSERT INTO "task" ("bvid", "cid", "format", "title", "owner", "cover", "status", "folder", "duration")
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		task.Bvid,
		task.Cid,
		task.Format,
		task.Title,
		task.Owner,
		task.Cover,
		task.Status,
		task.Folder,
		task.Duration,
	)
	if err != nil {
		return err
	}

	task.ID, err = result.LastInsertId()
	task.CreateAt = time.Now()
	return err
}

// Create 创建任务，并将任务加入全局任务列表
func (task *Task) Start() {
	GlobalTaskMux.Lock()
	GlobalTaskList = append(GlobalTaskList, task)
	GlobalTaskMux.Unlock()
	db := util.GetDB()
	defer db.Close()
	sessdata, err := bilibili.GetSessdata(db)
	if err != nil {
		task.UpdateStatus(db, "error", fmt.Errorf("bilibili.GetSessdata: %v", err))
		return
	}
	client := &bilibili.BiliClient{SESSDATA: sessdata}

	GlobalDownloadSem.Acquire()
	task.UpdateStatus(db, "running")
	err = DownloadMedia(client, task.Audio, task, "audio")
	if err != nil {
		GlobalDownloadSem.Release()
		task.UpdateStatus(db, "error", fmt.Errorf("DownloadMedia: %v", err))
		return
	}
	err = DownloadMedia(client, task.Video, task, "video")
	if err != nil {
		GlobalDownloadSem.Release()
		task.UpdateStatus(db, "error", fmt.Errorf("DownloadMedia: %v", err))
		return
	}
	GlobalDownloadSem.Release()

	outputPath := filepath.Join(task.Folder,
		fmt.Sprintf("%s %s.mp4", util.FilterFileName(task.Title),
			strings.Replace(base64.StdEncoding.EncodeToString([]byte(strconv.FormatInt(task.ID, 10))), "=", "", -1),
		),
	)

	videoPath := filepath.Join(task.Folder, strconv.FormatInt(task.ID, 10)+".video")
	audioPath := filepath.Join(task.Folder, strconv.FormatInt(task.ID, 10)+".audio")
	GlobalMergeSem.Acquire()
	err = task.MergeMedia(outputPath, videoPath, audioPath)
	if err != nil {
		GlobalMergeSem.Release()
		task.UpdateStatus(db, "error", fmt.Errorf("task.MergeMedia: %v", err))
		return
	}
	err = os.Remove(videoPath)
	if err != nil {
		GlobalMergeSem.Release()
		task.UpdateStatus(db, "error", fmt.Errorf("os.Remove: %v", err))
		return
	}
	err = os.Remove(audioPath)
	if err != nil {
		GlobalMergeSem.Release()
		task.UpdateStatus(db, "error", fmt.Errorf("os.Remove: %v", err))
		return
	}
	GlobalMergeSem.Release()
	task.UpdateStatus(db, "done")
}

// 合并音视频
func (task *Task) MergeMedia(outputPath string, inputPaths ...string) error {
	inputs := []string{}
	for _, path := range inputPaths {
		inputs = append(inputs, "-i", path)
	}
	cmd := exec.Command("ffmpeg", append(inputs, "-c:v", "copy", "-c:a", "copy", "-progress", "pipe:1", outputPath)...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	scanner := bufio.NewScanner(stdout)

	progress := newProgressBar(int64(task.Duration))
	outTimeRegex := regexp.MustCompile(`out_time_ms=(\d+)`) // 毫秒

	for scanner.Scan() {
		line := scanner.Text()
		match := outTimeRegex.FindStringSubmatch(line)
		if len(match) == 2 {
			outTime, err := strconv.ParseInt(match[1], 10, 64)
			if err != nil {
				return err
			}
			progress.current = outTime / 1000000
			task.MergeProgress = progress.percent()
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}
	task.MergeProgress = 1
	return nil
}

func GetVideoURL(medias []bilibili.Media, format common.MediaFormat) (string, error) {
	for _, code := range []int{12, 7, 13} {
		for _, item := range medias {
			if item.ID == format && item.Codecid == code {
				return item.BaseURL, nil
			}
		}
	}
	return "", errors.New("未找到对应视频分辨率格式")
}

func GetAudioURL(dash *bilibili.Dash) string {
	if dash.Flac != nil {
		return dash.Flac.Audio.BaseURL
	}
	var maxAudioID common.MediaFormat
	var audioURL string
	for _, item := range dash.Audio {
		if item.ID > maxAudioID {
			maxAudioID = item.ID
			audioURL = item.BaseURL
		}
	}
	return audioURL
}

func (task *Task) UpdateStatus(db *sql.DB, status TaskStatus, errs ...error) error {
	_, err := db.Exec(`UPDATE "task" SET "status" = ? WHERE "id" = ?`, status, task.ID)
	for _, err := range errs {
		if err != nil {
			err = util.CreateLog(db, fmt.Sprintf("Task-%d-Error: %v", task.ID, err))
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
	task.Status = status
	return err
}

func DownloadMedia(client *bilibili.BiliClient, _url string, task *Task, mediaType string) error {
	var resp *http.Response
	var err error
	for i := 0; i < 5; i++ {
		resp, err = client.SimpleGET(_url, nil)
		if err == nil {
			break
		}
	}

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

// GetCurrentFolder 获取数据库中的下载保存路径，如果不存在则将默认路径保存到数据库
func GetCurrentFolder(db *sql.DB) (string, error) {
	var filepath string
	err := db.QueryRow(`SELECT "value" FROM "field" WHERE "name" = 'download_folder'`).Scan(&filepath)
	if err != nil {
		if err == sql.ErrNoRows {
			folder, err := util.GetDefaultDownloadFolder()
			if err != nil {
				return "", err
			}
			err = SaveDownloadFolder(db, folder)
			if err != nil {
				return "", err
			}
			return folder, nil
		}
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
	_, err = db.Exec(`INSERT OR REPLACE INTO "field" ("name", "value") VALUES ('download_folder', ?)`, downloadFolder)
	return err
}

func GetTaskList(db *sql.DB, page int, pageSize int) ([]TaskInDB, error) {
	tasks := []TaskInDB{}

	rows, err := db.Query(`SELECT
		"id", "bvid", "cid", "format", "title",
		"owner", "cover", "status", "folder", "create_at"
	FROM "task" ORDER BY "id" DESC LIMIT ?, ?`,
		page*pageSize, pageSize,
	)
	if err != nil {
		return nil, err
	}

	createAt := ""

	for rows.Next() {
		task := TaskInDB{}
		err = rows.Scan(
			&task.ID,
			&task.Bvid,
			&task.Cid,
			&task.Format,
			&task.Title,
			&task.Owner,
			&task.Cover,
			&task.Status,
			&task.Folder,
			&createAt,
		)
		if err != nil {
			return nil, err
		}
		task.CreateAt, err = time.Parse("2006-01-02 15:04:05", createAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// InitHistoryTask 初始化上一次程序退出时未完成的任务，将进度变为 error
func InitHistoryTask(db *sql.DB) error {
	_, err := db.Exec(`UPDATE "task" SET "status" = 'error' WHERE "status" IN ('waiting', 'running')`)
	return err
}
