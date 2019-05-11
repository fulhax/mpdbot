package plugin

import (
	"log"
	"plugin"
)

type MpdMediaSource interface {
	IsSong(string) bool
	GetURI(string) (string, string, error) //URI, track name, error
	Auth()
}

func LoadPlugins(soFiles []string) *[]MpdMediaSource {
	sources := make([]MpdMediaSource, 0)
	for _, v := range soFiles {
		plug, err := plugin.Open(v)
		if err != nil {
			log.Printf("Failed to load plugin %s\n", v)
			continue
		}

		source, err := plug.Lookup("MpdMediaSource")
		if err != nil {
			log.Printf("%s", err.Error())
			return &sources
		}

		var songSource MpdMediaSource
		songSource, ok := source.(MpdMediaSource)
		if !ok {
			log.Printf("Failed to cast to MpdMediaSource")
			return &sources
		}

		songSource.Auth()
		log.Printf("%s", songSource)
		sources = append(sources, songSource)
	}

	return &sources
}

// type PluginHandler struct {
// 	loadedPlugins []Plugin
// }
//
// type Plugin struct {
// 	name    string
// 	version string
// 	plugin  string
// }

// func (p *PluginHandler) loadPlugin(s string) {
// 	fmt.Printf("Loading plugin %s\n", s)
//
// 	p, err := plugin.Open(str)
// 	if err != nil {
// 		fmt.Println("Error")
// 	}
// }
