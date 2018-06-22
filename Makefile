
clean:
	rm -rf debug
	rm -rf dict
deps:
	go get -u github.com/PuerkitoBio/goquery
	go get -u github.com/fatih/color
	go get -u github.com/boltdb/bolt/...

build: clean deps
	# -x 打印编译期间所用到的其它命令
	# -ldflags 里的  -s 去掉符号信息， -w 去掉DWARF调试信息
	go build -x -ldflags "-s -w" -o dict

install: build
	# 建立一个软链到usr/local/bin
	ln dict /usr/local/bin