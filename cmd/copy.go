package cmd

import (
	"fmt"
	"math"
	"os"
	"runtime"

	"github.com/HakimMouchquelita/GO-F22ssh/src/fs"
	"github.com/HakimMouchquelita/GO-F22ssh/src/multithread"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var RootCommand = &cobra.Command{
	Use:   "copy <source> <target>",
	Short: "Allows you to synchronize 2 directories recursively",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		origin := args[0]
		target := args[1]

		folders, files, err := fs.Scan(origin, target)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if fs.PROGRESS {
			fs.PROGRESS_BAR = progressbar.Default(int64(len(files) + len(folders)))
		}

		if fs.MULTI_THREADS {
			total := float64(len(files) + len(folders))
			cpus := float64(runtime.NumCPU())
			chunck := int(math.Ceil((total / cpus)))
			tasks := make(chan bool)

			// FOLDER
			missions := multithread.SplitFiles(folders, chunck)
			for _, f := range missions {
				go multithread.CopyFolder(tasks, f)
			}
			var missionsComplete int = 0
			for range tasks {
				missionsComplete++
				if missionsComplete == len(missions) {
					break
				}
			}

			// FILE
			missions = multithread.SplitFiles(files, chunck)
			for _, f := range missions {
				go multithread.CopyFile(tasks, f)
			}
			missionsComplete = 0
			for range tasks {
				missionsComplete++
				if missionsComplete == len(missions) {
					break
				}
			}
		} else {
			if err := fs.Copy(origin, target); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		os.Exit(0)
	},
}
