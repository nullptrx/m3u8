package main

import (
	"flag"
	"fmt"
	"github.com/nullptrx/v2/common"
	"github.com/nullptrx/v2/dl"
	nurl "net/url"
	"os"
	"strings"
)

var (
	url      string
	output   string
	chanSize int
	verbose  bool
	key      string
	merge    bool
	direct   bool
	proxy    string
)

func init() {
	flag.StringVar(&url, "u", "", "URL, required")
	flag.IntVar(&chanSize, "c", 10, "Maximum number of occurrences")
	flag.StringVar(&output, "o", "", "Output folder, required")
	flag.BoolVar(&verbose, "v", false, "Verbose log, optional")
	flag.StringVar(&key, "k", "", "Key path, optional")
	flag.BoolVar(&merge, "m", false, "Merge files, optional")
	flag.BoolVar(&direct, "d", false, "Enable direct connect. no proxy if enabled.")
	flag.StringVar(&proxy, "p", "socks5://127.0.0.1:7890", "Proxy url (such as socks://127.0.0.1:1080, http://127.0.0.1:1080), optional")
}

func main() {
	flag.Parse()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("[error]", r)
			os.Exit(-1)
		}
	}()
	if !direct {
		common.Proxy = proxy
	}
	if !strings.HasPrefix(url, "http") {
		if len(flag.Args()) > 0 {
			for _, arg := range flag.Args() {
				if strings.HasPrefix(arg, "http") {
					url = arg
					break
				}
			}
		}
		if url == "" {
			panicParameter("u")
		}
	}
	//if output == "" {
	//	panicParameter("o")
	//}
	if chanSize <= 0 {
		panic("parameter 'c' must be greater than 0")
	}

	u, err := nurl.Parse(url)
	if err != nil {
		panicParameter("u")
	}

	isM3u8 := strings.HasSuffix(u.Path, ".m3u8")
	if isM3u8 {
		downloader, err := dl.NewTask(output, url, verbose, key)
		if err != nil {
			panic(err)
		}
		if merge {
			if err := downloader.Merge(); err != nil {
				panic(err)
			}
		} else {
			if err := downloader.Start(chanSize); err != nil {
				panic(err)
			}
		}
		fmt.Println("Done!")
	} else {
		dl.DirectDownload(output, url, chanSize, verbose)
	}

}

func panicParameter(name string) {
	panic("parameter '" + name + "' is required")
}
