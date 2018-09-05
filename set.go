package main

import (
	"log"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xproto"
)

type Gammer interface {
	Gamma(size, temp int) (r, g, b []uint16)
}

type X struct {
	conn  *xgb.Conn
	root  xproto.Window
	crtcs []randr.Crtc
	g     Gammer
}

func NewX(g Gammer) (*X, error) {
	conn, err := xgb.NewConn()
	if err != nil {
		return nil, err
	}
	if err := randr.Init(conn); err != nil {
		return nil, err
	}

	root := xproto.Setup(conn).DefaultScreen(conn).Root

	res, err := randr.GetScreenResourcesCurrent(conn, root).Reply()
	if err != nil {
		return nil, err
	}

	return &X{conn: conn, root: root, crtcs: res.Crtcs, g: g}, nil
}

func (x *X) Close() {
	x.conn.Close()
}

func (x *X) Set(temp int) error {
	log.Println("set", temp)
	for _, crtc := range x.crtcs {
		size, err := randr.GetCrtcGammaSize(x.conn, crtc).Reply()
		if err != nil {
			return err
		}
		r, g, b := x.g.Gamma(int(size.Size), temp)
		if err := randr.SetCrtcGammaChecked(x.conn, crtc, size.Size, r, g, b).Check(); err != nil {
			return err
		}
	}
	return nil
}
