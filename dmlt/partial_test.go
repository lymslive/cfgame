package dmlt

import (
	"log"
	"testing"

	"github.com/lymslive/gotoo/assert"
)

func TestPartial(t *testing.T) {
	assert.BeginTest(t)
	defer assert.EndTest()

	pe, err := ParseOper("")
	assert.True(err == nil, "ParseOper fails")

	log.Printf("files result %d group", len(pe))
	for key, files := range pe {
		log.Printf("pe[%s] = []\n", key)
		for _, file := range files {
			log.Printf("| %v ", file)
		}
		log.Printf("\n")
	}
}
