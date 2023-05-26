package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/biogo/hts/bam"
	"github.com/biogo/hts/bgzf"
	"github.com/biogo/hts/sam"
	"github.com/liserjrqlxue/goUtil/osUtil"
	"github.com/liserjrqlxue/goUtil/simpleUtil"
)

var (
	input = flag.String(
		"i",
		"",
		"input bam",
	)
	out1 = flag.String(
		"1",
		"",
		"output read1",
	)
	out2 = flag.String(
		"2",
		"",
		"output read2",
	)
)

func main() {
	flag.Parse()
	if *input == "" || *out1 == "" || *out2 == "" {
		flag.Usage()
		log.Fatal("-i/-o required!")
	}

	// input
	var f = osUtil.Open(*input)
	defer simpleUtil.DeferClose(f)
	if !simpleUtil.HandleError(bgzf.HasEOF(f)).(bool) {
		log.Printf("file %q has no bgzf magic block: may be truncated", *input)
	}
	var br = simpleUtil.HandleError(bam.NewReader(f, 1)).(*bam.Reader)

	// output
	var o1 = osUtil.Create(*out1)
	defer simpleUtil.DeferClose(o1)
	var o2 = osUtil.Create(*out2)
	defer simpleUtil.DeferClose(o2)

	// PE
	var (
		read1 = make(map[string]*sam.Record, 1e6)
		read2 = make(map[string]*sam.Record, 1e6)
	)
	br2pe(br, read1, read2)

	for s, r1 := range read1 {
		var r2, ok = read2[s]
		if !ok {
			continue
		}
		simpleUtil.HandleError(fmt.Fprint(o1, record2fq(r1)))
		simpleUtil.HandleError(fmt.Fprint(o2, record2fq(r2)))
	}
}

func br2pe(br *bam.Reader, r1, r2 map[string]*sam.Record) {
	for {
		var (
			r, err = br.Read()
		)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("fail to read BAM record: [%v]", err)
		}
		if r.Flags&sam.Read1 == sam.Read1 {
			r1[r.Name] = r
		}
		if r.Flags&sam.Read2 == sam.Read2 {
			r2[r.Name] = r
		}
	}
}

// from https://forum.golangbridge.org/t/easy-way-for-letter-substitution-reverse-complementary-dna-sequence/20101
// from https://go.dev/play/p/IXI6PY7XUXN
var dnaComplement = strings.NewReplacer(
	"A", "T",
	"T", "A",
	"G", "C",
	"C", "G",
	"a", "t",
	"t", "a",
	"g", "c",
	"c", "g",
)

func record2fq(r *sam.Record) string {
	var (
		name = r.Name
		seq  = formatSeq(r.Seq)
		note = "+"
		qual = formatQual(r.Qual)
	)
	if r.Flags&sam.Read1 == sam.Read1 {
		name += "/1"
	}
	if r.Flags&sam.Read2 == sam.Read2 {
		name += "/2"
	}
	if r.Flags&sam.Reverse == sam.Reverse {
		note = "-"
		seq = Reverse(seq)
		qual = Reverse(qual)
		return fmt.Sprintf("%s\n%s\n%s\n%s\n", name, dnaComplement.Replace(string(seq)), note, qual)
	} else {
		return fmt.Sprintf("%s\n%s\n%s\n%s\n", name, seq, note, qual)
	}
}

// from https://github.com/biogo/hts/blob/master/sam/record.go
func formatSeq(s sam.Seq) []byte {
	if s.Length == 0 {
		return []byte{'*'}
	}
	return s.Expand()
}
func formatQual(q []byte) []byte {
	for _, v := range q {
		if v != 0xff {
			a := make([]byte, len(q))
			for i, p := range q {
				a[i] = p + 33
			}
			return a
		}
	}
	return []byte{'*'}
}

// Reverse returns its argument string reversed rune-wise left to right.
// from https://github.com/golang/example/blob/master/stringutil/reverse.go
func Reverse(r []byte) []byte {
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return r
}
