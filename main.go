package main

import (
	"flag"
	"log"
	"sort"
	"strings"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xproto"
)

const (
	minTemp = 1000
	maxTemp = 10000
	defTemp = 6500
)

var preset = map[string]int{
	"candle":      2300,
	"tungsten":    2700,
	"halogen":     3400,
	"fluorescent": 4200,
	"daylight":    5000,
}

func usage(m map[string]int) string {
	var s []string
	for k := range m {
		s = append(s, k)
	}
	sort.Strings(s)
	return strings.Join(s, ", ")
}

type Whitepoints []struct{ R, G, B float64 }

var wp = Whitepoints{
	{1.00000000, 0.18172716, 0.00000000}, /* 1000K */
	{1.00000000, 0.42322816, 0.00000000},
	{1.00000000, 0.54360078, 0.08679949},
	{1.00000000, 0.64373109, 0.28819679},
	{1.00000000, 0.71976951, 0.42860152},
	{1.00000000, 0.77987699, 0.54642268},
	{1.00000000, 0.82854786, 0.64816570},
	{1.00000000, 0.86860704, 0.73688797},
	{1.00000000, 0.90198230, 0.81465502},
	{1.00000000, 0.93853986, 0.88130458},
	{1.00000000, 0.97107439, 0.94305985},
	{1.00000000, 1.00000000, 1.00000000}, /* 6500K */
	{0.95160805, 0.96983355, 1.00000000},
	{0.91194747, 0.94470005, 1.00000000},
	{0.87906581, 0.92357340, 1.00000000},
	{0.85139976, 0.90559011, 1.00000000},
	{0.82782969, 0.89011714, 1.00000000},
	{0.80753191, 0.87667891, 1.00000000},
	{0.78988728, 0.86491137, 1.00000000}, /* 10000K */
	{0.77442176, 0.85453121, 1.00000000},
}

func (wp Whitepoints) Gamma(temp int) (r, g, b float64) {
	temp -= 1000
	ratio := float64(temp%500) / 500.0
	i, j := temp/500, temp/500+1
	r = wp[i].R*(1-ratio) + wp[j].R*ratio
	g = wp[i].G*(1-ratio) + wp[j].G*ratio
	b = wp[i].B*(1-ratio) + wp[j].B*ratio
	return
}

func Gamma(size int, temp int) (r, g, b []uint16) {
	gammar, gammag, gammab := wp.Gamma(temp)
	r = make([]uint16, size)
	g = make([]uint16, size)
	b = make([]uint16, size)
	for i := 0; i < size; i++ {
		gamma := 65535.0 * float64(i) / float64(size)
		r[i] = uint16(gamma * gammar)
		g[i] = uint16(gamma * gammag)
		b[i] = uint16(gamma * gammab)
	}
	return
}

func SetTemp(temp int) error {
	conn, err := xgb.NewConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := randr.Init(conn); err != nil {
		return err
	}

	root := xproto.Setup(conn).DefaultScreen(conn).Root
	res, err := randr.GetScreenResourcesCurrent(conn, root).Reply()
	if err != nil {
		return err
	}

	for _, crtc := range res.Crtcs {
		size, err := randr.GetCrtcGammaSize(conn, crtc).Reply()
		if err != nil {
			return err
		}
		r, g, b := Gamma(int(size.Size), temp)
		err = randr.SetCrtcGammaChecked(conn, crtc, size.Size, r, g, b).Check()
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	t := flag.Int("temp", defTemp, "color temperature")
	p := flag.String("preset", "", usage(preset))
	flag.Parse()

	if *t < minTemp || *t > maxTemp {
		*t = defTemp
	}
	if pr, ok := preset[*p]; ok {
		*t = pr
	}

	log.Printf("Set %vK", *t)
	if err := SetTemp(*t); err != nil {
		log.Fatal(err)
	}
}
