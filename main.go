package main

import (
	"github.com/im-auld/waypoint/cmd"
)

func main() {
	// f, err := os.Create("cpu.prof")
	// if err != nil {
	// 	log.Fatal("could not create CPU profile: ", err)
	// }
	// if err := pprof.StartCPUProfile(f); err != nil {
	// 	log.Fatal("could not start CPU profile: ", err)
	// }
	// defer pprof.StopCPUProfile()

	cmd.Execute()

	// f, err = os.Create("mem.prof")
	// if err != nil {
	// 	log.Fatal("could not create memory profile: ", err)
	// }
	// runtime.GC() // get up-to-date statistics
	// if err := pprof.WriteHeapProfile(f); err != nil {
	// 	log.Fatal("could not write memory profile: ", err)
	// }
	// f.Close()
}
