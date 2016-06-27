package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strconv"
	"strings"
)

/*
{
    "Beans & Legumes":{
        "Baked Beans":[
            {
        		"name" : "Baked Beans",
                "uri" : "http://www.fatsecret.com/calories-nutrition/usda/baked-beans",
                "calories": 382,
                "text":
                "totalNutrients": [
                {
      					"label": "Fat",
      					"quantity": 0.37910000000000005,
      					"unit": "g"
    			},
    			]
                "totalDaily": []
            },
            ....
        ],
        ....
    },
    ....
}

*/

type totNut struct {
	label    string
	quantity float64
	unit     string
}

type foodPar struct {
	name           string
	uri            string
	calories       int
	text           string
	totalNutrients []totNut
	totalDaily     []totNut
}

func mainPage(baseuri string) {
	doc, err := goquery.NewDocument(baseuri + "/calories-nutrition/")
	if err != nil {
		log.Fatal(err)
	}

	s := doc.Find(".generic.common")
	var name, uri string
	finalRet := make(map[string]map[string][]foodPar)
	s.Find("a").Each(func(i int, sel *goquery.Selection) {
		if i%3 == 1 {
			name = sel.Text()
		} else if i%3 == 0 {
			uri, _ = sel.Attr("href")
		} else {
			func() {
				doc1, err := goquery.NewDocument(baseuri + uri)
				if err != nil {
					log.Fatal(err)
				}

				s := doc1.Find(".secHolder")
				ge := make(map[string][]foodPar)
				s.Find("h2").Each(func(i int, sel *goquery.Selection) {
					selNext := sel.Next()
					selNext.Find("a").Each(func(j int, sele *goquery.Selection) {
						uri, _ = sele.Attr("href")
						ge[sel.Text()] = append(ge[sel.Text()], thirdPage(baseuri, uri))
					})
				})
				finalRet[name] = ge
			}()
		}
	})
	js, _ := json.Marshal(finalRet)
	fmt.Println(string(js))
}

func thirdPage(baseuri, uri string) foodPar {
	var theFood foodPar
	var theNut, theDaily totNut
	doc, err := goquery.NewDocument(baseuri + uri)
	if err != nil {
		log.Fatal(err)
	}
	theFood.name = doc.Find("h1").Text()
	theFood.text = doc.Find(".generic .spaced").Eq(1).Find("td").Text()
	theFood.uri = baseuri + uri
	theFood.calories, _ = strconv.Atoi(doc.Find("div.factValue").Eq(0).Text())
	doc.Find("td.label.borderTop").Each(func(i int, s *goquery.Selection) {
		if i >= 2 && i < doc.Find("td.label.borderTop").Length()-2 {
			hey := strings.Split(strings.Replace(strings.TrimSpace(s.Text()), string(9), "", -1), "\n")
			theNut.label = hey[0]
			var val float64
			if strings.Contains(hey[1], "mg") {
				val, _ = strconv.ParseFloat(strings.Replace(hey[1], "mg", "", -1), 64)
				val = float64(val / 1000)
			} else {
				if strings.Contains(hey[1], "-") {
					val = 0
				} else {
					val, _ = strconv.ParseFloat(strings.Replace(hey[1], "g", "", -1), 64)
				}
			}
			theNut.quantity = val
			theNut.unit = "g"
			theFood.totalNutrients = append(theFood.totalNutrients, theNut)
			//Now the daily percent
			theDaily.label = hey[0]
			theDaily.unit = "%"
			cont := s.Next().Text()
			if strings.Contains(cont, "%") {
				val, _ = strconv.ParseFloat(strings.Replace(strings.TrimSpace(cont), "%", "", -1), 64)
			} else {
				val = 0
			}
			theDaily.quantity = val
			theFood.totalDaily = append(theFood.totalDaily, theDaily)
		}
	})
	return theFood

}

func main() {
	baseuri := "http://www.fatsecret.com"
	mainPage(baseuri)
}
