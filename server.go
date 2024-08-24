package main

import (
	"fmt"
	"log"
	"net"
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
func IsPublicIP(IP net.IP) bool {
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := IP.To4(); ip4 != nil {
		switch {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}
	return false
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
		queries := c.IPs()
		for _, ip := range queries {
			IP := net.ParseIP(ip)
			fmt.Println(IP, ":", IsPublicIP(IP))
			if IsPublicIP(IP) {
				query = ip
			}
		}

		addr, addErr := netip.ParseAddr(query)
		if addErr != nil {
			return c.SendString("Invalid Protocol|1|" + query)
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
			return c.SendString("Invalid Structure|2|" + query)
			// log.Panic(err)
		}
		//
		if record.ASN == "" {
			return c.SendString("Invalid Lookup|3|" + query)
			// log.Panic(err)
		}
		//
		return c.SendString(query + "|" + record.ASN + "|" + record.Name + "|" + record.Domain)
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
			return c.SendString("Invalid Protocol|1|" + query)
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
			return c.SendString("Invalid Structure|2|" + query)
			// log.Panic(err)
		}
		//
		if record.ASN == "" {
			return c.SendString("Invalid Lookup|3|" + query)
			// log.Panic(err)
		}
		// fmt.Println(record.ASN)
		// fmt.Println(record.Domain)
		// fmt.Println(record.Name)
		//
		return c.SendString(query + "|" + record.ASN + "|" + record.Name + "|" + record.Domain)
	})
	//
	app.Listen(getPort())
}
