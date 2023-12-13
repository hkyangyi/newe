package ftpmag

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/jlaffaye/ftp"
)

// ftp管理器
type FtpMag struct {
	Conn map[string]*FtpObj
}

type FtpObj struct {
	Key            string
	Url            string
	Username       string
	Password       string
	Status         bool
	LoginStatus    bool
	ErrorMsg       string
	Client         *ftp.ServerConn
	Mutex          *sync.Mutex
	Send           chan *UpFile // 上传队列
	WorkPath       string
	WorkPathStatus bool
	AuthMakeDir    bool
}

type UpFile struct {
	Path     string
	ToPath   string
	FileName string
	Statuc   chan bool
	ErrMsg   string
}

// 启动管理器
func NewFtpNag() *FtpMag {
	var conn = make(map[string]*FtpObj)
	var ftpmag = FtpMag{
		Conn: conn,
	}
	go ftpmag.RunMag()
	return &ftpmag
}

// FTP链接管理器
func (a *FtpMag) RunMag() {
	for {
		for _, v := range a.Conn {
			if !v.Status {
				err := v.Clent()
				if err != nil {
					v.ErrorMsg = err.Error()
					continue
				}
			}
			if !v.LoginStatus {
				err := v.Login()
				if err != nil {
					v.ErrorMsg = err.Error()
					continue
				}
			}
			if !v.WorkPathStatus {
				err := v.GoToWorkPath()
				if err != nil {
					v.ErrorMsg = err.Error()
					continue
				}
			}
		}
		time.Sleep(10 * time.Second)
	}
}

// 上传文件监听
func (a *FtpObj) UpFile() {
	for {
		select {
		case file, ok := <-a.Send:
			if ok {
				//检测链接状态
				if !a.Status {
					err := a.Clent()
					if err != nil {
						file.ErrMsg = err.Error()
						file.UpError()
						continue
						//链接失败
					}
				}
				//检测登录状态
				if !a.LoginStatus {
					err := a.Login()
					if err != nil {
						file.ErrMsg = err.Error()
						file.UpError()
						continue
					}
				}

				//检测工作目录
				if a.WorkPathStatus {
					err := a.UpLoad(file)
					if err != nil {
						file.ErrMsg = err.Error()
						file.UpError()
						a.Refresh()
						continue
					}
					file.UpSuccess()
				} else {
					file.ErrMsg = "无法找到远程目录"
					file.UpError()
					continue
				}
			}
		}
	}
}

// 添加FTP
func (a *FtpMag) Create(url, user, pass, path string, AuthMakeDir int) *FtpObj {
	key := fmt.Sprintf("%s&%s&%s&%s", url, user, pass, path)
	var b bool
	if AuthMakeDir == 1 {
		b = true
	} else {
		b = false
	}
	if _, ok := a.Conn[key]; ok {
		a.Conn[key].Del()
	}

	var ftpitem = FtpObj{
		Key:         key,
		Url:         url,
		Username:    user,
		Password:    pass,
		WorkPath:    path,
		Status:      false,
		Send:        make(chan *UpFile, 1), //队列
		Mutex:       &sync.Mutex{},
		AuthMakeDir: b,
	}
	a.Conn[key] = &ftpitem
	go ftpitem.UpFile()
	return &ftpitem
}

// 获取FTP
func (a *FtpMag) GetFtp(key string) (*FtpObj, error) {
	if _, ok := a.Conn[key]; ok {
		return a.Conn[key], nil
	}
	return nil, errors.New("未找到FTP")
}

func (a *FtpMag) Del(key string) error {
	if _, ok := a.Conn[key]; ok {
		a.Conn[key].Client.Quit()
		close(a.Conn[key].Send)
		delete(a.Conn, key)
		return nil
	}
	return nil
}

// ftp链接
func (a *FtpObj) Clent() error {
	a.Mutex.Lock()
	if a.Status {
		a.Mutex.Unlock()
		return nil
	}
	client, err := ftp.Dial(a.Url, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		a.ErrorMsg = err.Error()
		a.Mutex.Unlock()
		return err
	}
	a.Client = client
	a.Status = true
	a.Mutex.Unlock()
	return nil
}

func (a *FtpObj) Close() {
	a.Client.Quit()
	a.Status = false
}

func (a *FtpObj) Del() {
	if a.Status == true {
		a.Client.Quit()
		a.Status = false
	}

}

// 连FTP
func (a *FtpObj) Refresh() {
	a.Client.Quit()
	a.Status = false
	a.LoginStatus = false
	a.WorkPathStatus = false
}

// ftp登录
func (a *FtpObj) Login() error {
	a.Mutex.Lock()
	if a.LoginStatus {
		a.Mutex.Unlock()
		return nil
	}
	err := a.Client.Login(a.Username, a.Password)
	if err != nil {
		a.Mutex.Unlock()
		return err
	}
	a.LoginStatus = true
	a.Mutex.Unlock()
	return nil
}

// 检测工作目录
func (a *FtpObj) GoToWorkPath() error {
	a.Mutex.Lock()
	if a.WorkPathStatus {
		a.Mutex.Unlock()
		return nil
	}

	err := a.Client.ChangeDir(a.WorkPath)
	if err != nil {
		//在这里写创建远程目录
		if a.AuthMakeDir {
			err = a.Client.MakeDir(a.WorkPath)
			if err != nil {
				a.Status = false
				a.LoginStatus = false
				a.WorkPathStatus = false
				a.Client.Quit()
				a.ErrorMsg = err.Error()
				a.Mutex.Unlock()
				return err
			}

		} else {
			a.Status = false
			a.LoginStatus = false
			a.WorkPathStatus = false
			a.Client.Quit()
			a.ErrorMsg = err.Error()
			a.Mutex.Unlock()
			return err
		}
	}
	a.WorkPathStatus = true
	a.Mutex.Unlock()
	return nil
}

// 上传文件
func (a *FtpObj) UpSend(path, topath, fileName string) *UpFile {
	var file = UpFile{
		Path:     path,
		ToPath:   topath,
		FileName: fileName,
		Statuc:   make(chan bool, 1),
		ErrMsg:   "",
	}
	a.Send <- &file
	return &file
}

func (a *FtpObj) UpLoad(file *UpFile) error {
	f, err := os.Open(file.Path + "/" + file.FileName)
	if err != nil {
		return err
	}
	defer f.Close()
	err = a.Client.Stor(file.ToPath+"/"+file.FileName, f)
	if err != nil {
		return err
	}
	return nil
}

// 文件上传失败
func (a *UpFile) UpError() {
	a.Statuc <- false
}

// 文件上传恒功
func (a *UpFile) UpSuccess() {
	a.Statuc <- true
}
