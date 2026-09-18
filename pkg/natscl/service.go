package natscl

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/mrscorpio/uahelper/pkg/tagdata"
	"github.com/nats-io/nats.go"
)

type NatsCl struct {
	C         *nats.Conn
	OnlineBuf []float32
	TimeBuf   time.Time
	ClCmd     []byte
	Id        string
	BarBuf    []float64
}

func NewNats(addr string) (*NatsCl, error) {
	var ncw WasmNatsConnectionWrapper
	cl, err := nats.Connect(addr, nats.SetCustomDialer(ncw))
	if err != nil {
		cl = nil
	}
	clCmd := make([]byte, 6)
	src := rand.NewSource(time.Now().Unix())
	var id string
	for i := 1; i < len(clCmd); i++ {
		clCmd[i] = byte(src.Int63() % 256)
		id += fmt.Sprint(clCmd[i])
	}
	log.Println(id)

	return &NatsCl{C: cl, ClCmd: clCmd, Id: id}, err
}

func (nc *NatsCl) SendCurrent() error {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, nc.OnlineBuf)
	if err != nil {
		return err
	}
	err = nc.C.Publish("online", buf.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func (nc *NatsCl) GetCurrent() error {
	var buf bytes.Buffer
	recBuf := make([]int32, len(nc.OnlineBuf))

	_, err := nc.C.Subscribe("online", func(msg *nats.Msg) {
		buf.Write(msg.Data)
		err := binary.Read(&buf, binary.BigEndian, recBuf)
		if err != nil {
			log.Println(err)
		}
		for i, v := range recBuf {
			nc.OnlineBuf[i] = float32(v) / 1000
		}
		nc.BarBuf[1] = float64(nc.OnlineBuf[1] / 100)
		nc.TimeBuf, err = time.Parse("2006-01-02 15:04:05.000", msg.Header.Get("crtm"))
		if err != nil {
			log.Println(err)
		}
	})
	if err != nil {
		return err
	}

	return nil
}

func (nc *NatsCl) SendCmd(cmd byte, d *tagdata.AllTags) error {
	nc.ClCmd[0] = cmd
	err := nc.C.Publish("cmd", nc.ClCmd)
	if err != nil {
		return err
	}
	if cmd == 66 {
		_, err := nc.C.Subscribe(nc.Id, func(msg *nats.Msg) {
			log.Println("Receved bytes:", msg.Size())

			r, err := zip.NewReader(bytes.NewReader(msg.Data), int64(msg.Size()))
			if err != nil {
				log.Println(err)
			}

			buf := new(bytes.Buffer)
			for _, f := range r.File {
				rc, err := f.Open()
				if err != nil {
					log.Println(err)
				}
				n, err := buf.ReadFrom(rc)
				if err != nil {
					log.Println(err)
				}
				fmt.Println(n)
				rc.Close()
			}
			d.Mu.Lock()
			err = json.Unmarshal(buf.Bytes(), d)
			nc.OnlineBuf = make([]float32, len(d.Tag))
			nc.BarBuf = make([]float64, len(d.Tag))
			d.Mu.Unlock()
			if err != nil {
				log.Println(err)
			}
			log.Println(len(d.Tt))
		})
		if err != nil {
			return err
		}
	}
	return nil
}
