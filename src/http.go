/*http server*/

package main

import (
	"encoding/json"
	"io"
	"net/http"
	//"sync"
	//"time"
)

type http_svc struct {
	svc_t  *http.Server
	addr   string //http server addr
	out_ch chan int
}

/*
Start listen and serve
*/
func (svc *http_svc) StartHttpSvc() {
	WAIT_GROUP.Add(1)
	defer WAIT_GROUP.Done()

	go func() {
		svc.svc_t.ListenAndServe()
	}()

	<-svc.out_ch
	svc.svc_t.Shutdown(nil)
}

/*
server handle
*/
func (svc *http_svc) handle(w http.ResponseWriter, r *http.Request) {
	//io.WriteString(w, "Hello, world!\n")
	//fmt.Printf("%+v\n", r)
	//var ret bool

	tmp := ret_format{
		Ret:     false,
		Message: "Parameter error!",
	}
	data, _ := json.Marshal(tmp)
	data_str := string(data)

	/*解析URL参数*/
	r.ParseForm()
	var id_num, room_num string
	if len(r.Form) > 0 {
		room_num = r.Form.Get("room_num")
		id_num = r.Form.Get("id_num")
		LOG_TRAC.Println("room_num: ", room_num, ", Id_num: ", id_num)
	} else {
		jsonp_ret(w, data_str)
		return
	}

	if len(id_num) != 6 || len(room_num) == 0 {
		jsonp_ret(w, data_str)
		return
	}

	user := NewUser(id_num, room_num)
	ret := user.check_user()

	jsonp_ret(w, ret)

}

func jsonp_ret(w http.ResponseWriter, str string) {
	io.WriteString(w, "getSMSResult("+str+")")
}
