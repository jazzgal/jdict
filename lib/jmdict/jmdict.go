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

// Process data before push into bucket
func syncJMData(jmdict Jmdict) {
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
