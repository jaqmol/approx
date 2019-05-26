package main

import (
	"log"

	"github.com/jaqmol/approx/conf"
)

func main() {
	f := conf.ReadFormation()
	varsMap := make(map[string]bool)
	varsSlc := make([]string, 0)

	log.Println("Public configs:")
	for _, p := range f.PublicConfs {
		log.Printf("  %v\n", p.Name())
		log.Println("    Inputs:")
		for _, i := range p.Inputs() {
			log.Printf("      %v\n", i)
		}
		log.Println("    Outputs:")
		for _, o := range p.Outputs() {
			log.Printf("      %v\n", o)
		}
		if len(p.Required()) > 0 {
			log.Println("    Required:")
			for k, v := range p.Required() {
				log.Printf("      %v: %v\n", k, v)
			}
		}
	}

	log.Println("Private configs:")
	for _, p := range f.PrivateConfs {
		log.Printf("  %v\n", p.Name())
		log.Println("    Inputs:")
		for _, i := range p.Inputs() {
			log.Printf("      %v\n", i)
		}
		log.Println("    Outputs:")
		for _, o := range p.Outputs() {
			log.Printf("      %v\n", o)
		}
		if len(p.Required()) > 0 {
			log.Println("    Required:")
			for k, v := range p.Required() {
				log.Printf("      %v: %v\n", k, v)
			}
		}
	}

	if len(varsMap) > 0 {
		log.Println("Formation-local vars:")
		for v := range varsMap {
			log.Printf("  %v", v)
		}
	}
	if len(varsSlc) > 0 {
		log.Println("All occurencies of formation-local vars:")
		for _, v := range varsSlc {
			log.Printf("  %v", v)
		}
	}
}
