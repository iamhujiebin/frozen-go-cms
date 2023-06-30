package vap_r

import (
	"bytes"
	"frozen-go-cms/hilo-common/mycontext"
	"github.com/abema/go-mp4"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
)

type VapVapcReq struct {
	Mp4 string `form:"mp4"`
}

// @Tags VAP
// @Summary 获取vapc
// @Param Authorization header string true "token"
// @Param mp4 formData string true "mp4"
// @Success 200 {object} []Music
// @Router /v1_0/vap/vapc [post]
func VapVapc(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	var param VapVapcReq
	if err := c.ShouldBindJSON(&param); err != nil {
		return myCtx, err
	}
	var vapc string
	downloadFile(param.Mp4)
	vapc = do2()
	c.Writer.Write([]byte(vapc))
	//resp.ResponseOk(c, vapc)
	return myCtx, nil
}

func BoxTypeVapc() mp4.BoxType { return mp4.StrToBoxType("vapc") }

func init() {
	mp4.AddBoxDef(&Vapc{}, 0)
}

type Vapc struct {
	mp4.FullBox `mp4:"0,extend"`
	UI32        uint32 `mp4:"1,size=32"`
	ByteArray   []byte `mp4:"2,size=8,len=dynamic"`
}

func (*Vapc) GetType() mp4.BoxType {
	return BoxTypeVapc()
}

func do2() string {
	// 打开mp4文件
	file, err := os.Open("example.mp4")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	//var vapc string
	buf := bytes.NewBuffer(nil)
	_, err = mp4.ReadBoxStructure(file, func(h *mp4.ReadHandle) (interface{}, error) {
		if h.BoxInfo.Type.String() == "vapc" {
			_, _ = h.ReadData(buf)
			return h.Expand()
		}
		return nil, nil
	})
	return buf.String()
}

func downloadFile(url string) *os.File {
	response, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer response.Body.Close()

	file, err := os.Create("example.mp4")
	if err != nil {
		return nil
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return nil
	}
	//file.Seek(0, 0)
	return file
}
