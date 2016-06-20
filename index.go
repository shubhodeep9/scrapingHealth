package main

import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
)

/*
{
    "Beans & Legumes":{
        "Baked Beans":[
            {
                "Name": "Vegetarian Baked Beans",
                "Nutrition summary":{
                    "Calories":"382",
                    "Fat":"13.03",
                    "Carbs":"54.12",
                    "Protein":"14.02",
                    "Description":"There are 239 calories in 1 cup of Vegetarian Baked Beans.
                                Calorie breakdown: 3% fat, 79% carbs, 18% protein."
                },
                "Common Serving Sizes":[
                    {
                        "Size":"1 oz",
                        "Calories":"27"
                    },
                    ....
                ],

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

	nutritionSummary := make(map[string]string)
	nutritionSummary["Calories"] = doc.Find(".fact").Eq(0).Find(".factValue").Text()
	nutritionSummary["Fat"] = doc.Find(".fact").Eq(1).Find(".factValue").Text()
	nutritionSummary["Carbs"] = doc.Find(".fact").Eq(2).Find(".factValue").Text()
	nutritionSummary["Protein"] = doc.Find(".fact").Eq(3).Find(".factValue").Text()
	nutritionSummary["Description"] = doc.Find(".generic .spaced").Eq(1).Find("td").Text()
	fmt.Println(nutritionSummary)
}

func main() {
	baseuri := "http://www.fatsecret.com"
	thirdPage(baseuri)
}
