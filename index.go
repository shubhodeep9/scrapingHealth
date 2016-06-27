package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strings"
)

/*
{
    "Beans & Legumes":{
        "Baked Beans":[
            {
                "uri" : "http://www.fatsecret.com/calories-nutrition/usda/baked-beans",
                "calories": 382,
                "text":
                "totalNutrients": {
                	"FAT": {
      					"label": "Fat",
      					"quantity": 0.37910000000000005,
      					"unit": "g"
    				},
     			}
                "totalDaily":
            },
            ....
        ],
        ....
    },
    ....
}

*/

func mainPage(baseuri string) {
	doc, err := goquery.NewDocument(baseuri + "/calories-nutrition/")
	if err != nil {
		log.Fatal(err)
	}

	s := doc.Find(".generic .common")
	s.Find("a").Each(func(i int, sel *goquery.Selection) {
		if i%3 == 1 {
			fmt.Println(sel.Text())
		} else if i%3 == 0 {
			fmt.Println(sel.Attr("href"))
		}
	})
}

func secondPage(baseuri string) {
	doc, err := goquery.NewDocument(baseuri + "/calories-nutrition/group/beans-and-legumes")
	if err != nil {
		log.Fatal(err)
	}

	s := doc.Find(".secHolder")
	ge := make(map[string][]string)
	s.Find("h2").Each(func(i int, sel *goquery.Selection) {
		selNext := sel.Next()
		selNext.Find("a").Each(func(j int, sele *goquery.Selection) {
			ge[sel.Text()] = append(ge[sel.Text()], sele.Text())
		})
	})
	fmt.Println(ge)
}

func thirdPage(baseuri string) {
	doc, err := goquery.NewDocument(baseuri + "/calories-nutrition/usda/baked-beans")
	if err != nil {
		log.Fatal(err)
	}
	//text := doc.Find(".generic .spaced").Eq(1).Find("td").Text()
	doc.Find("td.label.borderTop").Each(func(i int, s *goquery.Selection) {
		if i >= 2 && i < doc.Find("td.label.borderTop").Length()-2 {
			hey := strings.Replace(strings.TrimSpace(s.Text()), string(9), "", -1)
			fmt.Println(i, hey)
		}
	})

}

func main() {
	baseuri := "http://www.fatsecret.com"
	thirdPage(baseuri)
}
