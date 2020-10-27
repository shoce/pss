/*
history:
20/290 v1
20/293 option "0" to show the tree of processes starting with with pid=0
20/301 first arg is pid to specify the root process

release:
rm -f pss.linux.gz pss.macos.gz
GOOS=linux GOARCH=amd64 go build -trimpath -o pss.linux
GOOS=darwin GOARCH=amd64 go build -trimpath -o pss.macos
gzip pss.linux pss.macos
open .

version=`{it vv}
git init
git add pss.go readme.text
git commit -am $version
git branch -M main
git remote add origin git@github.com:shoce/pss.git
git push -f -u origin main
echo git push --delete origin v
echo rm -rf .git
echo https://github.com/shoce/pss/releases
echo $version
echo rm -f pss.linux.gz pss.macos.gz

GoFmt GoBuildNull GoBuild GoRun
*/

package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	ps "github.com/mitchellh/go-ps"
)

type Process struct {
	Pid   int
	Ppid  int
	Ppids []int
	Name  string
}

var (
	pp []Process
)

func main() {
	var err error
	var rootpid int

	if len(os.Args) > 1 {
		rootpid, err = strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "First arg `%s`: %v\n", os.Args[1], err)
			os.Exit(1)
		}
	} else {
		rootpid = 1
	}

	pp0, err := ps.Processes()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	for _, p0 := range pp0 {
		p := Process{Pid: p0.Pid(), Ppid: p0.PPid(), Name: p0.Executable()}
		pp = append(pp, p)
	}

	sort.Slice(pp, func(i, j int) bool {
		if pp[i].Ppid < pp[j].Ppid {
			return true
		}
		if pp[i].Ppid > pp[j].Ppid {
			return false
		}
		return pp[i].Pid < pp[j].Pid
	})

	for i, p := range pp {
		if p.Pid == p.Ppid {
			continue
		}
		if p.Ppid > 0 {
			pp[i].Ppids = []int{p.Ppid}
		}
		for _, q := range pp {
			if q.Pid == p.Ppid {
				pp[i].Ppids = append(q.Ppids, pp[i].Ppids...)
			}
		}
	}

	for _, p := range pp {
		ppidss := ""
		rootpidinppids := false

		for _, ppid := range p.Ppids {
			if ppid == rootpid {
				rootpidinppids = true
			}
			ppidss += fmt.Sprintf("%d\t", ppid)
		}
		if rootpid == 0 || p.Pid == rootpid || rootpidinppids {
			fmt.Printf("%s%d\t%s\n", ppidss, p.Pid, p.Name)
		}
	}
}
