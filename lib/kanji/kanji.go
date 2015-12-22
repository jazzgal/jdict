package kanji

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

// Helper functions
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Structs for parsing XML data
type (
	Kanjidict struct {
		XMLName xml.Name    `xml:"kanjidic2"`
		Entries []Character `xml:"character"`
	}
	Character struct {
		XMLName   xml.Name         `xml:"character"`
		Literal   string           `xml:"literal"`
		Radical   RadElement       `xml:"radical"`
		Refs      []DictRef        `xml:"dic_number>dic_ref"`
		Strokes   int              `xml:"misc>stroke_count"`
		Grade     int              `xml:"misc>grade"`
		Freq      int              `xml:"misc>freq"`
		Jlpt      int              `xml:"misc>jlpt"`
		QueryCode []CodeElement    `xml:"query_code>q_code"`
		Meanings  []MeaningElement `xml:"reading_meaning>rmgroup>meaning"`
		Readings  []ReadingElement `xml:"reading_meaning>rmgroup>reading"`
		// Readings that are now only associated with names
		Nanori []string `xml:"reading_meaning>nanori"`
	}
	// Query code for finding kanji
	CodeElement struct {
		Code  string `xml:"qc_type,attr"`
		Value string `xml:",chardata"`
	}
	// Meaning
	MeaningElement struct {
		Value string `xml:",chardata"`
		Lang  string `xml:"m_lang,attr"`
	}
	ReadingElement struct {
		Value string `xml:",chardata"`
		Type  string `xml:"r_type,attr"`
	}
	// Reference to published books, page
	RefElement struct {
		XMLName xml.Name  `xml:"dic_number"`
		Info    []DictRef `xml:"dic_ref"`
	}
	DictRef struct {
		XMLName xml.Name `xml:"dic_ref"`
		Src     string   `xml:"dr_type,attr"`
		Index   string   `xml:",chardata"`
		Page    int      `xml:"m_page,attr"`
		Vol     int      `xml:"m_vol,attr"`
	}
	// Radical values
	RadElement struct {
		XMLName  xml.Name   `xml:"radical"`
		Radattrs []RadValue `xml:"rad_value"`
	}
	RadValue struct {
		XMLName xml.Name `xml:"rad_value"`
		Type    string   `xml:"rad_type,attr"`
		Value   string   `xml:",chardata"`
	}
)

func parseXMLDict(filePath string) Kanjidict {
	b, err := ioutil.ReadFile(filePath)
	CheckErr(err)
	// Read XML
	s := Kanjidict{}
	d := xml.NewDecoder(bytes.NewReader(b))
	err = d.Decode(&s)
	CheckErr(err)
	_, err = json.MarshalIndent(s, "", "  ")
	CheckErr(err)
	return s
}

func syncKanjiData(kanjidict Kanjidict) {
	fmt.Println(kanjidict.Entries[0])
}

func PopulateData(filePath string) {
	fmt.Println("Reading kanji xml  at ...", filePath)
	kanjis := parseXMLDict(filePath)
	fmt.Printf("Parsed %d Kanji entries\n", len(kanjis.Entries))
	fmt.Println("Push data to database ...")
	syncKanjiData(kanjis)
	fmt.Println("Done!")
}
