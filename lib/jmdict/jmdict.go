package jmdict

import (
	"bufio"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"time"
)

// Structs for parsing XML data
type (
	Jmdict struct {
		XMLName xml.Name `xml:"JMdict"`
		Entries []JapEng `xml:"entry"`
	}

	JapEng struct {
		XMLName xml.Name `xml:"entry"`
		// Unique number sequence from original
		EntrySequence int              `xml:"ent_seq"`
		Kanji         []KanjiElement   `xml:"k_ele"`
		Reading       []ReadingElement `xml:"r_ele"`
		Sense         []SenseElement   `xml:"sense"`
	}
	ReadingElement struct {
		XMLName xml.Name `xml:"r_ele"`
		// Kana reading
		Reb string `xml:"reb"`
		// Gairaigo
		NoKanji xml.Name `xml:"re_nokanji"`
		// With kanji apply
		RestrictedTo []string `xml:"re_restr"`
		// Info
		ReInf []string `xml:"re_inf"`
		RePri []string `xml:"re_pri"`
	}
	KanjiElement struct {
		XMLName xml.Name `xml:"k_ele"`
		// Shorphrase represent
		Keb string `xml:"keb"`
		// Information irregular
		KeInf []string `xml:"ke_inf"`
		// Priority
		KePri []string `xml:"ke_pri"`
	}
	SenseElement struct {
		XMLName             *xml.Name `json:"-" xml:"sense"`
		RestrictedToKanji   []string  `json:"-" xml:"stagk"`
		RestrictedToReading []string  `xml:"stagr" json:"-"`
		CrossReference      []string  `xml:"xref" json:",omitempty"`
		Antonym             []string  `xml:"ant" json:",omitempty"`
		PartOfSpeech        []string  `xml:"pos" json:",omitempty"`
		Field               []string  `xml:"field" json:",omitempty"`
		Misc                []string  `xml:"misc" json:",omitempty"`
		Info                []string  `xml:"s_info" json:",omitempty"`
		Dialect             []string  `xml:"dial" json:",omitempty"`
		// Source from a loan word
		LSource []LangSource `xml:"lsource" json:",omitempty"`
		Gloss   []string     `xml:"gloss" json:",omitempty"`
	}
	LangSource struct {
		XMLName xml.Name `xml:"lsource"`
		// Part or full
		Lang  string `xml:"lang,attr"`
		Key   string `xml:",chardata"`
		LType string `xml:"ls_type,attr"`
		// Construct from words not actual phrase
		LWasei string `xml:"ls_wasei,attr"`
	}
)

// DTO struct
type (
	EntryCollect struct {
		Id       int64
		Sequence int
		// Query key - data indexes
		Keys       map[string][]int
		ReadingSet [][]Reb
		KanjiSet   [][]Keb
		SenseSet   [][]Sense
	}
	// Short phrase contain kanji
	Keb struct {
		Key string
		Pri []string `json:",omitempty"`
		Inf []string `json:",omitempty"`
	}
	// Short phrase contain only kana
	Reb struct {
		Key string
		Pri []string `json:",omitempty"`
		Inf []string `json:",omitempty"`
	}
	// Phrase's meaning & meta info
	Sense struct {
		Gloss []string
		Meta  map[string][]string `json:",omitempty"`
	}
)

// Helper functions
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Return file content with string by lines
func ReadTextByLines(filePath string) []string {
	var result []string
	absPath, _ := filepath.Abs(filePath)
	file, err := os.Open(absPath)
	CheckErr(err)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Loop each line and read
		result = append(result, scanner.Text())
	}
	return result
}

// Read jmdict xml file
func parseXMLDict(filePath string) Jmdict {
	b, err := ioutil.ReadFile(filePath)
	CheckErr(err)
	// Read XML
	s := Jmdict{}
	d := xml.NewDecoder(bytes.NewReader(b))
	d.Entity = map[string]string{}

	// Manual define custom terms
	//TODO(hails) move to config
	termList := ReadTextByLines("data/term_list.txt")
	var term string
	for index := 0; index < len(termList); index++ {
		term = termList[index]
		d.Entity[term] = term
	}
	err = d.Decode(&s)
	CheckErr(err)
	_, err = json.MarshalIndent(s, "", "  ")
	CheckErr(err)

	// Test
	for _, e := range s.Entries {
		if e.EntrySequence == 1000310 {
			fmt.Println(e.Reading[4].NoKanji)
		}

	}
	return s
}

// Add new set of nodes if not existed. Return set index
func combineSet(nodes []Node, updateSet *[][]Node) int {
	// Check if set exist & update
	for idx, set := range *updateSet {
		if reflect.DeepEqual(nodes, set) {
			return idx
		}
	}
	*updateSet = append(*updateSet, nodes)
	return len(*updateSet) - 1
}

// Sort out index from reb, keb key. Return indexes & entry collect
func makeIndexes(entry JapEng) ([]string, EntryCollect) {
	graph, rGraph := buildGraphs(entry)

	var kanjiSet, readingSet, senseSet [][]Node
	indexes := map[string][]int{}

	// Traverse graph to make nodes combination
	for _, keb := range entry.Kanji {
		n := makeKNode(keb)
		dSet := DFS(n, graph)
		// Make reading & sense combination set
		indexes[n.id] = make([]int, 3)
		indexes[n.id][0] = -1
		for w, nodes := range dSet {
			sort.Sort(XSortablePoints{nodes})
			var updateSet *[][]Node
			if w == 2 {
				updateSet = &readingSet
			} else {
				updateSet = &senseSet
			}
			idx := combineSet(nodes, updateSet)
			indexes[n.id][w-1] = idx
		}
	}

	// Make reb query mapping
	for _, reb := range entry.Reading {
		n := makeRNode(reb)
		sSet := DFS(n, graph)
		kSet := DFS(n, rGraph)
		fmt.Println("------")
		fmt.Println(reb.Reb)
		indexes[n.id] = make([]int, 3)
		indexes[n.id][1] = -1
		for w, nodes := range kSet {
			// Must be level 1
			sort.Sort(XSortablePoints{nodes})
			idx := combineSet(nodes, &kanjiSet)
			indexes[n.id][w-1] = idx
		}
		for w, nodes := range sSet {
			// Must be level 3
			sort.Sort(XSortablePoints{nodes})
			idx := combineSet(nodes, &senseSet)
			indexes[n.id][w-1] = idx
		}
	}

	// Make sense query mapping
	for _, sense := range entry.Sense {
		n := makeSNode(sense)
		dSet := DFS(n, rGraph)
		// Make kanji & reading combination set
		for w, nodes := range dSet {
			sort.Sort(XSortablePoints{nodes})
			var updateSet *[][]Node
			if w == 2 {
				updateSet = &readingSet
			} else {
				updateSet = &kanjiSet
			}
			combineSet(nodes, updateSet)
			//TODO(hails) Sense query
		}
	}

	collect := makeEntryCollect(indexes, kanjiSet, readingSet, senseSet,
		entry.EntrySequence)
	// Return query keys only
	queries := make([]string, 0, len(indexes))
	for k := range indexes {
		queries = append(queries, k)
	}

	return queries, collect
}

// Create entry collection from node & indexes combination
func makeEntryCollect(indexes map[string][]int, kanjiSet [][]Node,
	readingSet [][]Node, senseSet [][]Node, sequence int) EntryCollect {
	// Make collect entry to be stored
	collect := EntryCollect{}
	collect.Sequence = sequence
	collect.Keys = indexes
	for _, set := range kanjiSet {
		var setDTO []Keb
		for _, node := range set {
			// Cast back
			data, _ := node.data.(KanjiElement)
			el := Keb{data.Keb, data.KePri, data.KeInf}
			setDTO = append(setDTO, el)
		}
		collect.KanjiSet = append(collect.KanjiSet, setDTO)
	}
	for _, set := range readingSet {
		var setDTO []Reb
		for _, node := range set {
			// Cast back
			data, _ := node.data.(ReadingElement)
			el := Reb{data.Reb, data.RePri, data.ReInf}
			setDTO = append(setDTO, el)
		}
		collect.ReadingSet = append(collect.ReadingSet, setDTO)
	}
	for _, set := range senseSet {
		var setDTO []Sense
		for _, node := range set {
			// Cast back
			data, _ := node.data.(SenseElement)
			el := Sense{data.Gloss, nil}
			// Make meta mapping
			data.Gloss = nil
			metaRaw, _ := json.Marshal(data)
			err := json.Unmarshal(json.RawMessage(metaRaw), &el.Meta)
			CheckErr(err)
			setDTO = append(setDTO, el)
		}
		collect.SenseSet = append(collect.SenseSet, setDTO)
	}

	return collect
}

// Process data before push into bucket
func syncJMData(jmdict Jmdict) {
	//db, err := bolt.Open("jdict.db", 0644, nil)
	//CheckErr(err)
	//defer db.Close()
	//var entry JapEng
	for _, entry := range jmdict.Entries {
		//0. Indices query key
		//Test entry
		if entry.EntrySequence == 1000390 {
			fmt.Println(entry)
			makeIndexes(entry)
		}
		//1. Store data to 'entry' bucket
		//2. Update 'index' bucket
	}

}

// Push data from source file into database
func PopulateData(filePath string) {
	fmt.Println("Reading xml ...")
	jmdict := parseXMLDict(filePath)
	fmt.Printf("Parsed %d entries\n", len(jmdict.Entries))
	fmt.Println("Push data to database ...")
	syncJMData(jmdict)
	fmt.Println("Done!")
}
