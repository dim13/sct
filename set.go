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
	sizes []uint16
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

	sizes := make([]uint16, len(res.Crtcs))
	for i, crtc := range res.Crtcs {
		size, err := randr.GetCrtcGammaSize(conn, crtc).Reply()
		if err != nil {
			return nil, err
		}
		sizes[i] = size.Size
	}

	return &X{
		conn:  conn,
		root:  root,
		crtcs: res.Crtcs,
		sizes: sizes,
		g:     g,
	}, nil
}

func (x *X) Close() {
	x.conn.Close()
}

func (x *X) Set(temp int) error {
	log.Println("set", temp)
	for i, crtc := range x.crtcs {
		size := x.sizes[i]
		r, g, b := x.g.Gamma(int(size), temp)
		if err := randr.SetCrtcGammaChecked(x.conn, crtc, size, r, g, b).Check(); err != nil {
			return err
		}
	}
	return nil
}
