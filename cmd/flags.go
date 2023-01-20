package cmd

import "github.com/HakimMouchquelita/golang/src/fs"

func init() {
	RootCommand.PersistentFlags().BoolVarP(&fs.DEBUG, "debug", "d", fs.DEBUG, "enable debug mode for display tracing")
	RootCommand.PersistentFlags().BoolVarP(&fs.MULTI_THREADS, "multithread", "m", fs.MULTI_THREADS, "use all cpus (warning in some case its useless)")
	RootCommand.PersistentFlags().BoolVarP(&fs.PROGRESS, "progress", "p", fs.PROGRESS, "display progress bar")
}
