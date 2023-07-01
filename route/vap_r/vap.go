package vap_r

import (
	"bytes"
	"frozen-go-cms/hilo-common/mycontext"
	"github.com/abema/go-mp4"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
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
	client := http.Client{}
	resp, err := client.Get(param.Mp4)
	if err != nil {
		return myCtx, err
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return myCtx, err
	}

	file := bytes.NewReader(buf.Bytes())
	vapcBuf := bytes.NewBuffer(nil)
	_, err = mp4.ReadBoxStructure(file, func(h *mp4.ReadHandle) (interface{}, error) {
		if h.BoxInfo.Type.String() == "vapc" {
			_, _ = h.ReadData(vapcBuf)
			return h.Expand()
		}
		return nil, nil
	})
	c.Writer.Write(vapcBuf.Bytes())
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
