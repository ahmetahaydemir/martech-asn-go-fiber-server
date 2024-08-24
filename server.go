package main

import (
	"fmt"
	"log"
	"net/netip"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/oschwald/maxminddb-golang/v2"
)

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":3000"
	} else {
		port = ":" + port
	}

	return port
}
func main() {
	app := fiber.New()
	app.Use(cors.New())
	//
	app.Get("/asn", func(c *fiber.Ctx) error {
		db, err := maxminddb.Open("asn.mmdb")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		//
		query := c.IP()
		addr, addErr := netip.ParseAddr(query)
		if addErr != nil {
			return c.SendString("Invalid Protocol - 1")
			// panic(addErr)
		}
		//
		var record struct {
			Domain string `maxminddb:"domain"`
			ASN    string `maxminddb:"asn"`
			Name   string `maxminddb:"name"`
		}
		err = db.Lookup(addr).Decode(&record)
		if err != nil {
			return c.SendString("Invalid Structure - 2")
			// log.Panic(err)
		}
		//
		if record.ASN == "" {
			return c.SendString("Invalid Lookup - 3")
			// log.Panic(err)
		}
		fmt.Println(record.ASN)
		fmt.Println(record.Domain)
		fmt.Println(record.Name)
		//
		return c.SendString(c.IP() + "|" + record.ASN + "|" + record.Name + "|" + record.Domain)
	})
	app.Get("/asn/:value", func(c *fiber.Ctx) error {
		db, err := maxminddb.Open("asn.mmdb")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		//
		query := c.Params("value")
		addr, addErr := netip.ParseAddr(query)
		if addErr != nil {
			return c.SendString("Invalid Protocol - 1")
			// panic(addErr)
		}
		//
		var record struct {
			Domain string `maxminddb:"domain"`
			ASN    string `maxminddb:"asn"`
			Name   string `maxminddb:"name"`
		}
		err = db.Lookup(addr).Decode(&record)
		if err != nil {
			return c.SendString("Invalid Structure - 2")
			// log.Panic(err)
		}
		//
		if record.ASN == "" {
			return c.SendString("Invalid Lookup - 3")
			// log.Panic(err)
		}
		fmt.Println(record.ASN)
		fmt.Println(record.Domain)
		fmt.Println(record.Name)
		//
		return c.SendString(query + "|" + record.ASN + "|" + record.Name + "|" + record.Domain)
	})
	//
	app.Listen(getPort())
}
