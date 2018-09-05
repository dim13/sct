set term png
set output "wp.png"
set title "White point"
set xlabel "Kelvin"
set grid
set key off
set yrange [-0.05:1.05]
plot	"wp.txt" using 1:2 title "red"   smooth csplines lc rgb "#ff0000", \
	"wp.txt" using 1:3 title "green" smooth csplines lc rgb "#00ff00", \
	"wp.txt" using 1:4 title "blue"  smooth csplines lc rgb "#0000ff"
