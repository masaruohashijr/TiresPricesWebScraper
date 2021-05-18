package main

import (
	c "4tiresWebScraper/config"
	m "4tiresWebScraper/models"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/sclevine/agouti"
)

var driver *agouti.WebDriver

var taxonomy m.Taxonomy

func main() {
	StartDriver()
	page, err := driver.NewPage()
	c.InitConfig()
	tiresCfg := c.ReadConfig()
	if err != nil {
		log.Fatal("Failed to open page:", err)
	}
	if err := page.Navigate(tiresCfg.TiresPricesURL); err != nil {
		log.Fatal("Failed to navigate:", err)
	}
	time.Sleep(1 * time.Second)
	fields := make(map[string]string)
	fields["WIDTH"] = "select2-tires-search-select-largeur-container"
	fields["PROFILE"] = "select2-tires-search-select-rapport-container"
	fields["DIAMETER"] = "select2-tires-search-select-diametre-container"
	println("********WIDTH*********")
	fmt.Printf("O campo eh %v.\n\n", fields["WIDTH"])
	fWidth := page.FindByID(fields["WIDTH"])
	fWidth.Click()
	time.Sleep(1 * time.Second)
	html, _ := page.HTML()
	widthIDs := GetAllOptionsIds(html)

	for i := 0; i < len(widthIDs); i++ {
		if i > 0 {
			html, _ := page.HTML()
			widthIDs = GetAllOptionsIds(html)
		}
		println("********WIDTH*********")
		v := widthIDs[i]
		opt := page.FindByID(v)
		optionW, _ := opt.Text()
		fmt.Printf("A chave é %v com conteudo %s\n\n", v, optionW)
		// Width 145
		opt.Click()
		time.Sleep(1 * time.Second)
		// PROFILE is set.
		println("******PROFILE********")
		fmt.Printf("O campo eh %v.\n\n", fields["PROFILE"])
		pWidth := page.FindByID(fields["PROFILE"])
		pWidth.Click()
		time.Sleep(1 * time.Second)
		html, _ := page.HTML()
		profileIDs := GetAllOptionsIds(html)
		// Select PROFILE Option
		for j := 0; j < len(profileIDs); j++ {
			v := profileIDs[j]
			fmt.Printf("Profile id %v\n", v)
			opt := page.FindByID(v)
			optionP, _ := opt.Text()
			println(j, optionP)
			opt.Click()
			time.Sleep(1 * time.Second)
			// DIAMETER is set.
			println("*****DIAMETER*********")
			fmt.Printf("O campo eh %v.\n\n", fields["DIAMETER"])
			dWidth := page.FindByID(fields["DIAMETER"])
			dWidth.Click()
			time.Sleep(500 * time.Millisecond)
			html, _ := page.HTML()
			newHtml := extract(html, "select2-tires-search-select-diametre-results")
			diametersIDs := GetAllOptionsIds(newHtml)
			// Select DIAMETER Option
			for m := 0; m < len(diametersIDs); m++ {
				v := diametersIDs[j]
				fmt.Printf("Diameter id %v\n", v)
				opt := page.FindByID(v)
				optionD, _ := opt.Text()
				println(j, optionD)
				opt.Click()
				time.Sleep(500 * time.Millisecond)
				searchButton := page.FindByClass("tires-search-form-buttons")
				searchButton.Click()
				time.Sleep(1 * time.Second)
				// Scrap Tires Prices.
				tires := scrapTires(page)
				page.Back()
				time.Sleep(1 * time.Second)
				optsD := make(map[string][]Tire)
				optsD[optionD] = tires
				optsP := make(map[string]map[string][]Tire)
				optsP[optionP] = optsD
				optsW := make(map[string]map[string]map[string][]Tire)
				optsW[optionW] = optsP
			}
		}
		fWidth = page.FindByID(fields["WIDTH"])
		fWidth.Click()
		time.Sleep(1 * time.Second)
	}
	time.Sleep(1 * time.Second)
	println("******FIM**********")
	// StopDriver()
}

func scrapTires(page *agouti.Page) (tires []Tire) {
	html, _ := page.HTML()
	// <span class="produits_list_item_header_b">
	/*	<h3 class="dimensions">
		<span>215</span><span class="dimen_r">/</span><span>40</span><span class="dimen_r">R</span><span>18</span>
		<span class="xl">XL</span>
		<span class="speed-ratings">89W</span>
		</h3>
	*/
	//<span id="item_price_19016" class="header_value_price">173.36$</span>
	pos := strings.Index(html, "<span class=\"produits_list_item_header_b\">")
	for pos > -1 {
		cursor := 0
		fmt.Printf("\nposicao %d: ", pos)
		snippet := html[pos:]
		posAng := strings.Index(snippet, ">")
		cursor += posAng
		snippet = snippet[posAng+1:]
		openAng := strings.Index(snippet, "<")
		cursor += openAng
		name := snippet[:openAng]
		fmt.Println("=============================")
		fmt.Printf("Name: %s\n", name)
		// Speed Rating
		posSR := strings.Index(snippet, "speed-ratings")
		cursor += posSR
		snippet = snippet[posSR:]
		posAng = strings.Index(snippet, ">")
		cursor += posAng
		snippet = snippet[posAng+1:]
		openAng = strings.Index(snippet, "<")
		speedRating := snippet[:openAng]
		cursor += openAng
		fmt.Printf("SR: %s\n", speedRating)
		// Price
		posPrice := strings.Index(snippet, "class=\"header_value_price\">")
		cursor += posPrice
		snippet = snippet[posPrice:]
		posAng = strings.Index(snippet, ">")
		cursor += posAng
		snippet = snippet[posAng+1:]
		openAng = strings.Index(snippet, "<")
		price := snippet[:openAng]
		cursor += openAng
		fmt.Printf("\nPrice: %s\n", price)
		// <span class="speed-ratings">
		tire := Tire{
			Name:        name,
			Price:       price,
			SpeedRating: speedRating,
		}
		tires = append(tires, tire)
		fmt.Printf("\ncursor %d: ", cursor)
		html = html[pos+cursor:]
		fmt.Printf("\nHTML: %s", html[0:100])
		pos = strings.Index(html, "<span class=\"produits_list_item_header_b\">")
	}
	return tires
}

func extract(html string, stop string) (snippet string) {
	pos := strings.Index(html, stop)
	if pos == -1 {
		println("ERRO")
		return html
	}
	snippet = html[pos:]
	return snippet
}

func GetAllOptionsIds(html string) (ids map[int]string) {
	set := make(map[int]string)

	posLI := strings.Index(html, "li class=\"select2-results__option")
	i := 0
	for posLI > -1 {
		posFinal := posLI + 300
		if posFinal > len(html) {
			posFinal = len(html) - 1
		}
		snippet := html[posLI:posFinal]
		posID := strings.Index(snippet, "id=\"")
		snippetID := snippet[posID+4:]
		posClosingId := strings.Index(snippetID, "\" ")
		ID := snippetID[:posClosingId]
		set[i] = ID
		html = html[posID:]
		posLI = strings.Index(html, "li class=\"select2-results__option\"")
		i++
	}
	return set
}

func StartDriver() {
	driver = agouti.ChromeDriver()
	if err := driver.Start(); err != nil {
		log.Fatal("Failed to start driver:", err)
	}
}

func StopDriver() {
	if err := driver.Stop(); err != nil {
		log.Fatal("Failed to close pages and stop WebDriver:", err)
	}
}
