package main

import (
	gzip_org "compress/gzip"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/biogo/hts/bam"
	"github.com/biogo/hts/bgzf"
	"github.com/biogo/hts/sam"
	gzip "github.com/klauspost/compress/gzip"
	"github.com/liserjrqlxue/goUtil/osUtil"
	"github.com/liserjrqlxue/goUtil/simpleUtil"
)

func TestMain(m *testing.M) {
	f = osUtil.Open("test.bam")
	defer simpleUtil.DeferClose(f)
	if !simpleUtil.HandleError(bgzf.HasEOF(f)).(bool) {
		log.Printf("file %q has no bgzf magic block: may be truncated", "test.bam")
	}
	br = simpleUtil.HandleError(bam.NewReader(f, 1)).(*bam.Reader)
	read1 = make(map[string]*sam.Record)
	read2 = make(map[string]*sam.Record)
	os.Exit(m.Run())
}

func BenchmarkBr2pe(b *testing.B) {
	for i := 0; i < b.N; i++ {
		br2pe(br, read1, read2)
	}
}

func BenchmarkRecord2fq(b *testing.B) {
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

func BenchmarkFmtFprint(b *testing.B) {
	var dir, err = os.MkdirTemp("", "fsdemo")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var o1 = osUtil.Create(filepath.Join(dir, "test.1.fq"))
		var o2 = osUtil.Create(filepath.Join(dir, "test.2.fq"))
		for s, r1 := range read1 {
			var r2, ok = read2[s]
			if !ok {
				continue
			}
			simpleUtil.HandleError(fmt.Fprint(o1, record2fq(r1)))
			simpleUtil.HandleError(fmt.Fprint(o2, record2fq(r2)))
			simpleUtil.HandleError(o2.WriteString(record2fq(r2)))
		}
	}
}

func BenchmarkWriteString(b *testing.B) {
	var dir, err = os.MkdirTemp("", "fsdemo")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var o1 = osUtil.Create(filepath.Join(dir, "test.1.fq"))
		var o2 = osUtil.Create(filepath.Join(dir, "test.2.fq"))
		for s, r1 := range read1 {
			var r2, ok = read2[s]
			if !ok {
				continue
			}
			simpleUtil.HandleError(o1.WriteString(record2fq(r1)))
			simpleUtil.HandleError(o2.WriteString(record2fq(r2)))
		}
		simpleUtil.CheckErr(o1.Close())
		simpleUtil.CheckErr(o2.Close())
	}
}

func BenchmarkWritePE(b *testing.B) {
	var dir, err = os.MkdirTemp("", "fsdemo")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		o1 = osUtil.Create(filepath.Join(dir, "test.1.fq"))
		o2 = osUtil.Create(filepath.Join(dir, "test.2.fq"))
		writePE(o1, o2, read1, read2)
		simpleUtil.CheckErr(o1.Close())
		simpleUtil.CheckErr(o2.Close())
	}
}

func BenchmarkWritePE_gzip_org(b *testing.B) {
	var dir, err = os.MkdirTemp("", "fsdemo")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		o1 = osUtil.Create(filepath.Join(dir, "test.1.fq.gz"))
		o2 = osUtil.Create(filepath.Join(dir, "test.2.fq.gz"))
		var zw1 = gzip_org.NewWriter(o1)
		var zw2 = gzip_org.NewWriter(o1)
		writePE(zw1, zw2, read1, read2)
		simpleUtil.CheckErr(zw1.Close())
		simpleUtil.CheckErr(zw2.Close())
		simpleUtil.CheckErr(o1.Close())
		simpleUtil.CheckErr(o2.Close())
	}
}

func BenchmarkWritePE_klauspos_gzip(b *testing.B) {
	var dir, err = os.MkdirTemp("", "fsdemo")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		o1 = osUtil.Create(filepath.Join(dir, "test.1.fq.gz"))
		o2 = osUtil.Create(filepath.Join(dir, "test.2.fq.gz"))
		var zw1 = gzip.NewWriter(o1)
		var zw2 = gzip.NewWriter(o1)
		writePE(zw1, zw2, read1, read2)
		simpleUtil.CheckErr(zw1.Close())
		simpleUtil.CheckErr(zw2.Close())
		simpleUtil.CheckErr(o1.Close())
		simpleUtil.CheckErr(o2.Close())
	}
}
