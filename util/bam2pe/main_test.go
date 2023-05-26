package main

import (
	"github.com/biogo/hts/bam"
	"github.com/biogo/hts/bgzf"
	"github.com/biogo/hts/sam"
	"github.com/liserjrqlxue/goUtil/osUtil"
	"github.com/liserjrqlxue/goUtil/simpleUtil"
	"log"
	"testing"
)

func BenchmarkBr2pe(b *testing.B) {
	var f = osUtil.Open("test.bam")
	defer simpleUtil.DeferClose(f)
	if !simpleUtil.HandleError(bgzf.HasEOF(f)).(bool) {
		log.Printf("file %q has no bgzf magic block: may be truncated", "test.bam")
	}
	var br = simpleUtil.HandleError(bam.NewReader(f, 1)).(*bam.Reader)

	// PE
	var (
		read1 = make(map[string]*sam.Record)
		read2 = make(map[string]*sam.Record)
	)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		br2pe(br, read1, read2)
	}
}

func BenchmarkRecord2fq(b *testing.B) {
	var f = osUtil.Open("test.bam")
	defer simpleUtil.DeferClose(f)
	if !simpleUtil.HandleError(bgzf.HasEOF(f)).(bool) {
		log.Printf("file %q has no bgzf magic block: may be truncated", "test.bam")
	}
	var br = simpleUtil.HandleError(bam.NewReader(f, 1)).(*bam.Reader)

	// PE
	var (
		read1 = make(map[string]*sam.Record)
		read2 = make(map[string]*sam.Record)
	)
	br2pe(br, read1, read2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for s, r1 := range read1 {
			var r2, ok = read2[s]
			if !ok {
				continue
			}
			record2fq(r1)
			record2fq(r2)
		}
	}
}
