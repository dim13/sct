package main

type Whitepoint struct {
	R, G, B float64
}

type Whitepoints []Whitepoint

var Points = Whitepoints{
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

func (p Whitepoints) avg(temp int) (r, g, b float64) {
	temp -= 1000
	ratio := float64(temp%500) / 500.0
	i, j := temp/500, temp/500+1
	r = p[i].R*(1-ratio) + p[j].R*ratio
	g = p[i].G*(1-ratio) + p[j].G*ratio
	b = p[i].B*(1-ratio) + p[j].B*ratio
	return r, g, b
}

func (p Whitepoints) Gamma(size, temp int) (r, g, b []uint16) {
	gammar, gammag, gammab := p.avg(temp)
	r = make([]uint16, size)
	g = make([]uint16, size)
	b = make([]uint16, size)
	for i := 0; i < size; i++ {
		gamma := 65535.0 * float64(i) / float64(size)
		r[i] = uint16(gamma * gammar)
		g[i] = uint16(gamma * gammag)
		b[i] = uint16(gamma * gammab)
	}
	return r, g, b
}
