#make && 6g libgogotest.go && 6l libgogotest.6 && ./6.out "libgogotest.go" && (make clean > /dev/null)
(make > /dev/null) && 6g libgogotest.go && 6l libgogotest.6 && ./6.out "libgogotest.go" && (make clean > /dev/null)
