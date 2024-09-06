package thumbnails

import (
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
)

/* Controls the no of concurrent processes used to generate thumbnails */
type GuardProcs struct {
	g   chan struct{}
	max int
}

func NewGuardProcs() GuardProcs {

	maxProcs := runtime.NumCPU()
	maxProcsEnv := os.Getenv("AVRP_MAX_PROCS")
	if maxProcsEnv != "" {
		mp, err := strconv.Atoi(maxProcsEnv)
		if err == nil {
			maxProcs = mp
		}
	}

	log.Printf("Using %v threads for thumbnail generation\n", maxProcs)

	return GuardProcs{
		g:   make(chan struct{}, maxProcs),
		max: maxProcs,
	}

}

func (p *GuardProcs) MaxProcs() int {
	return p.max
}

func (p *GuardProcs) In() {
	p.g <- struct{}{}
}

func (p *GuardProcs) Out() {
	<-p.g
}

type GuardFile struct {
	locks map[string]*sync.Mutex
}

func NewGuardFile() GuardFile {
	return GuardFile{
		locks: make(map[string]*sync.Mutex),
	}
}

func (p *GuardFile) Lock(key string) {
	l, exists := p.locks[key]
	if exists {
		l.Lock()
	} else {
		p.locks[key] = &sync.Mutex{}
		p.locks[key].Lock()
	}
}

func (p *GuardFile) Unlock(key string) {
	l, exists := p.locks[key]
	if exists {
		l.Unlock()
	} else {
		log.Printf("ERR: %s not found in GuardFile\n", key)
	}
}
