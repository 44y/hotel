/*file process*/

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	Update_dir_period int = 10
)

type user_st struct {
	id_num       string
	room_num     string
	id_type      string
	online_time  *time.Time
	offline_time *time.Time
}

type onoff_st struct {
	online_time  *time.Time
	offline_time *time.Time
	file_name    []string //whole name
}

type ret_format struct {
	Ret       bool
	Auth_id   string //room num
	Auth_name string // id num and id type
	Message   string //extra message
}

var (
	File_Mtx sync.Mutex
	user_map map[string]onoff_st //key:room_num and id_num; value :on off line struct
)

/*
查询用户是否在线
ret :
    login       -- true
    no login    -- false
*/
func (user *user_st) check_user() string {

	user.load_log()
	LOG_TRAC.Println(user)

	tmp := new(ret_format)

	/*user not online*/
	if user.online_time == nil {
		LOG_TRAC.Println("user not online")
		tmp.Ret = false
		goto JSONandRet
		/*user online*/
	} else if user.offline_time == nil {
		LOG_TRAC.Println("user online")
		tmp.Ret = true
		tmp.Auth_id = user.room_num
		tmp.Auth_name = user.id_num + "_" + user.id_type
		goto JSONandRet
	}

	/*user online then offline*/
	if user.offline_time.After(*user.online_time) {
		LOG_TRAC.Println("user online then offline")
		tmp.Ret = false
		goto JSONandRet
		/*user offline then online*/
	} else {
		LOG_TRAC.Println("user offline then online")
		tmp.Ret = true
		tmp.Auth_id = user.room_num
		tmp.Auth_name = user.id_num + "_" + user.id_type
		goto JSONandRet
	}

JSONandRet:
	data, _ := json.Marshal(*tmp)
	return string(data)
}

/*
加载用户日志
*/
func (user *user_st) load_log() {
	File_Mtx.Lock()
	defer File_Mtx.Unlock()

	files, err := ioutil.ReadDir(DIR)
	if err != nil {
		LOG_ERR.Println("ReanDir error!")
		return
	}

	/*get system pathseparator*/
	sep := string(os.PathSeparator)
	/*for each file*/
	for _, tmp_file := range files {
		/*jump directory*/
		if tmp_file.IsDir() == false {

			whole_name := DIR + sep + tmp_file.Name()
			data, err := ioutil.ReadFile(whole_name)
			if err != nil {
				LOG_ERR.Println(err)
				return
			}
			data_str := string(data)
			data_slice := strings.Split(data_str, ";")
			/*删除格式不正确的文件*/
			if len(data_slice) < 6 {
				LOG_ERR.Println(whole_name, "file format wrong!")
				os.Remove(whole_name)
				LOG_TRAC.Println("delete wrong format file : ", whole_name)
				continue
			}
			if strings.ToUpper(data_slice[0]) == strings.ToUpper(user.room_num) && strings.HasSuffix(strings.ToUpper(data_slice[2]), strings.ToUpper(user.id_num)) {
				user.id_num = data_slice[2]
				/*get online time*/
				if len(data_slice[3]) != 0 {
					user.online_time, err = parseTime(data_slice[3])
					if err != nil {
						LOG_ERR.Println(err)
						/*删除格式不正确的文件*/
						os.Remove(whole_name)
						LOG_TRAC.Println("delete wrong format file : ", whole_name)
						continue
					}
				}
				/*get offline time*/
				if len(data_slice[4]) != 0 {
					user.offline_time, err = parseTime(data_slice[4])
					if err != nil {
						LOG_ERR.Println(err)
						/*删除格式不正确的文件*/
						os.Remove(whole_name)
						LOG_TRAC.Println("delete wrong format file : ", whole_name)
						continue
					}
				}
				user.id_type = data_slice[6]
				if len(user.id_type) != 2 {
					/*去掉结尾的\r\n*/
					user.id_type = string([]byte(user.id_type)[:2])
				}
			}
		}
	}
}

/*
创建新用户
*/
func NewUser(id, room string) *user_st {
	user := &user_st{
		id_num:   id,
		room_num: room,
	}
	return user
}

/*
周期性删除下线用户数据
*/
func UpdateDir(out_ch chan int) {
	WAIT_GROUP.Add(1)
	defer WAIT_GROUP.Done()
	time_up := time.After(time.Second)
	for {
		select {
		case <-time_up:
			LOG_INFO.Println("Start updating dir")
			updateHandle()
			time_up = time.After(time.Duration(Update_dir_period) * time.Second)
		case <-out_ch:
			LOG_INFO.Println("Got out channel,bye")
		}
	}
}

func updateHandle() {
	File_Mtx.Lock()
	defer File_Mtx.Unlock()

	files, err := ioutil.ReadDir(DIR)
	if err != nil {
		LOG_ERR.Println("ReanDir error!")
		return
	}

	/*get system pathseparator*/
	sep := string(os.PathSeparator)
	/*for each file*/

	user_map := make(map[string]onoff_st) //key:room_num and id_num; value :on off line struct

	for _, tmp_file := range files {
		/*jump directory*/
		if tmp_file.IsDir() == false {

			whole_name := DIR + sep + tmp_file.Name()
			data, err := ioutil.ReadFile(whole_name)
			if err != nil {
				LOG_ERR.Println(err)
				return
			}
			data_str := string(data)
			data_slice := strings.Split(data_str, ";")
			/*删除格式不正确的文件*/
			if len(data_slice) < 6 {
				LOG_ERR.Println(whole_name, "file format wrong!")
				os.Remove(whole_name)
				LOG_TRAC.Println("delete wrong format file : ", whole_name)
				continue
			}

			user_num := data_slice[0] + data_slice[2]

			/*get on off line*/
			var tmp_st onoff_st
			if len(data_slice[3]) != 0 {
				tmp_st.online_time, err = parseTime(data_slice[3])
				if err != nil {
					LOG_ERR.Println(err)
					/*删除格式不正确的文件*/
					os.Remove(whole_name)
					LOG_TRAC.Println("delete wrong format file : ", whole_name)
					continue
				}
			}

			if len(data_slice[4]) != 0 {
				tmp_st.offline_time, err = parseTime(data_slice[4])
				if err != nil {
					LOG_ERR.Println(err)
					/*删除格式不正确的文件*/
					os.Remove(whole_name)
					LOG_TRAC.Println("delete wrong format file : ", whole_name)
					continue
				}
			}

			tmp_st.file_name = append(tmp_st.file_name, whole_name)
			user_map[user_num] = tmp_st
		}
	}
	//LOG_TRAC.Println("user_map: ", user_map)
	DeleteUnuseFile(&user_map)
}

/*
	获取上下线时间
	RET error：
		nil 成功
		not nil 失败
*/
func parseTime(ori_time string) (*time.Time, error) {
	parsed_time, err := time.Parse("2006-01-02 15:04:05", ori_time)
	if err != nil {
		/*航天日志可能会有不同格式*/
		parsed_time, err = time.Parse("2006-01-02", ori_time)
		if err != nil {
			return nil, err
		}
	}

	return &parsed_time, nil
}

func DeleteUnuseFile(um *map[string]onoff_st) {
	for _, v := range *um {
		/*user offline, delete files*/
		if v.offline_time != nil && v.offline_time.After(*v.online_time) {
			for _, file := range v.file_name {
				os.Remove(file)
				LOG_TRAC.Println("delete file : ", file)
			}
		}
	}
}

func getCurrDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err == nil {
		return strings.Replace(dir, "\\", "/", -1)
	} else {
		panic(err)
	}

}

/*解析json配置文件，返回配置和文件路径*/
func parseConfig() (*conf_st, string) {
	cur_dir := getCurrDir()

	sep := string(os.PathSeparator)
	conf_path := cur_dir + sep + "conf.json"

	fi, err := os.Open(conf_path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()

	br := bufio.NewReader(fi)

	data, _, end := br.ReadLine()
	if end == io.EOF {
		panic("config file err")
	}
	data, _, end = br.ReadLine()
	if end == io.EOF {
		panic("config file err")
	}

	/*替换目录中的"\"*/
	data_str := string(data)
	new_data_str := strings.Replace(data_str, "\\", "/", -1)

	tmp_conf := &conf_st{}
	json.Unmarshal([]byte(new_data_str), tmp_conf)
	return tmp_conf, conf_path
}

/*更新新的配置到json文件*/
func updateConfig(new_conf *conf_st, conf_path string) bool {
	data, _ := json.Marshal(new_conf)

	old_f, err := ioutil.ReadFile(conf_path)
	if err != nil {
		fmt.Println(err)
		return false
	}
	lines := strings.Split(string(old_f), "\n")

	for i, _ := range lines {
		if i == 1 {
			lines[i] = string(data)
			break
		}
	}

	//把lines用"\n"连接起来组成一条新的string
	new_f := strings.Join(lines, "\n")

	ioutil.WriteFile(conf_path, []byte(new_f), 0666)
	return true
}

func isDirExist(name string) bool {
	fi, err := os.Stat(name)

	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}
}

/*判断当前端口是否可用*/
func isPortAvl(port int) bool {
	cmd := exec.Command("netstat", "-nao")
	out_byt, _ := cmd.Output()
	out_str := string(out_byt)
	return strings.Contains(out_str, strconv.Itoa(port)) == false
}

/*找到一个可用的端口*/
func getAvlPort() int {
	for new_port := 40000; new_port < 65535; new_port++ {
		if isPortAvl(new_port) {
			return new_port
		}
	}
	return 0
}

var (
	searchDir = [...]string{
		"C:/",
		//"C:/Users/jere/Desktop",
		"D:/",
		"E:/",
		"F:/"}
	log_dir_name = "outputtxt"
	//log_dir_found = false
)

/*搜索航天金盾用户log路径*/
func searchUserLog() (bool, string) {
	log_dir_found := false
	var found_path string
	for _, val := range searchDir {
		if log_dir_found {
			return true, found_path
		} else if isDirExist(val) {
			err := filepath.Walk(val, func(path string, f os.FileInfo, err error) error {
				if f != nil && f.IsDir() && f.Name() == log_dir_name {
					log_dir_found = true
					found_path = strings.Replace(path, "\\", "/", -1)
				}
				return nil
			})
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return false, ""
}
