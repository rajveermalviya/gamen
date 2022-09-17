package main

import (
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/exp/slices"
)

func main() {
	if len(os.Args) < 2 {
		panic("specify package names")
	}

	// serverUrl := getFedoraMirrorServer()
	// fmt.Println(serverUrl)
	serverUrl := "https://download-cc-rdu01.fedoraproject.org/pub/fedora/linux/releases/36/Everything/x86_64/os/repodata/repomd.xml"
	primaryXmlUrl := getPrimaryXmlUrl(serverUrl)
	metadata := getPackageXmlMetadata(primaryXmlUrl)
	packageUrls := getPackageUrls(
		serverUrl, metadata,
		os.Args[1:]...,
	)

	var wg sync.WaitGroup

	for pkg, pkgUrl := range packageUrls {
		pkg, pkgUrl := pkg, pkgUrl

		wg.Add(1)
		go func() {
			defer wg.Done()

			fmt.Println("Downloading: " + pkgUrl)

			res := mustV(http.Get(pkgUrl))
			defer res.Body.Close()

			f := mustV(os.Create(pkg + ".rpm"))
			defer f.Close()

			_ = mustV(io.Copy(f, res.Body))

			fmt.Println("Done downloading: " + pkgUrl)
		}()
	}

	wg.Wait()
}

func getFedoraMirrorServer() string {
	const FedoraReleaseEver = "36"
	const FedoraRepoUrl = "https://mirrors.fedoraproject.org/metalink?repo=fedora-" + FedoraReleaseEver + "&arch=x86_64"

	type rUrl struct {
		Text       string `xml:",chardata"`
		Protocol   string `xml:"protocol,attr"`
		Preference string `xml:"preference,attr"`
	}

	var r struct {
		XMLName xml.Name `xml:"metalink"`
		Files   struct {
			File []struct {
				Name      string `xml:"name,attr"`
				Resources struct {
					URL []rUrl `xml:"url"`
				} `xml:"resources"`
			} `xml:"file"`
		} `xml:"files"`
	}

	fmt.Println("Downloading: " + FedoraRepoUrl)
	res := mustV(http.Get(FedoraRepoUrl))
	defer res.Body.Close()
	fmt.Println("Done downloading: " + FedoraRepoUrl)

	must(xml.NewDecoder(res.Body).Decode(&r))

	urls := []rUrl{}
	for _, file := range r.Files.File {
		if file.Name == "repomd.xml" {
			for _, url := range file.Resources.URL {
				if url.Protocol == "https" {
					urls = append(urls, url)
				}
			}
			break
		}
	}
	sort.Slice(urls, func(i, j int) bool {
		return mustV(strconv.ParseFloat(urls[i].Preference, 64)) >
			mustV(strconv.ParseFloat(urls[j].Preference, 64))
	})

	if len(urls) == 0 {
		panic("getFedoraMirrorServer failed")
	}

	return strings.TrimSpace(urls[0].Text)
}

func getPrimaryXmlUrl(serverUrl string) string {
	fmt.Println("Downloading: " + serverUrl)
	res := mustV(http.Get(serverUrl))
	defer res.Body.Close()
	fmt.Println("Done downloading: " + serverUrl)

	var r struct {
		XMLName xml.Name `xml:"repomd"`
		Data    []struct {
			Type     string `xml:"type,attr"`
			Location struct {
				Text string `xml:",chardata"`
				Href string `xml:"href,attr"`
			} `xml:"location"`
		} `xml:"data"`
	}

	must(xml.NewDecoder(res.Body).Decode(&r))

	var primaryXmlUrl string
	for _, data := range r.Data {
		if data.Type == "primary" {
			primaryXmlUrl = strings.TrimSpace(data.Location.Href)
			break
		}
	}
	if primaryXmlUrl == "" {
		panic("getPrimaryXmlUrl failed")
	}

	return strings.TrimSuffix(serverUrl, "repodata/repomd.xml") + primaryXmlUrl
}

func getPackageUrls(serverUrl string, metadata PackageXmlMetadata, packages ...string) map[string]string {
	serverUrlPrefix := strings.TrimSuffix(serverUrl, "repodata/repomd.xml")

	packageUrls := map[string]string{}
	for _, pkg := range metadata.Package {
		if pkg.Type == "rpm" && (pkg.Arch == "x86_64" || pkg.Arch == "noarch") && slices.Contains(packages, pkg.Name) {
			pkgUrl := serverUrlPrefix + pkg.Location.Href
			packageUrls[pkg.Name] = pkgUrl
		}
	}

	if len(packageUrls) != len(packages) {
		notFound := []string{}
		for _, pkg := range packages {
			if _, ok := packageUrls[pkg]; !ok {
				notFound = append(notFound, pkg)
			}
		}

		panic(fmt.Sprintf("some packages not found: %v", notFound))
	}

	return packageUrls
}

type PackageXmlMetadata struct {
	XMLName  xml.Name `xml:"metadata"`
	Packages string   `xml:"packages,attr"`
	Package  []struct {
		Type     string `xml:"type,attr"`
		Name     string `xml:"name"`
		Arch     string `xml:"arch"`
		Location struct {
			Href string `xml:"href,attr"`
		} `xml:"location"`
	} `xml:"package"`
}

func getPackageXmlMetadata(primaryXmlUrl string) PackageXmlMetadata {
	fmt.Println("Downloading: " + primaryXmlUrl)
	res := mustV(http.Get(primaryXmlUrl))
	defer res.Body.Close()
	fmt.Println("Done downloading: " + primaryXmlUrl)

	var metadata PackageXmlMetadata

	if strings.HasSuffix(primaryXmlUrl, ".gz") {
		gzipReader := mustV(gzip.NewReader(res.Body))
		defer gzipReader.Close()

		must(
			xml.NewDecoder(gzipReader).
				Decode(&metadata),
		)
	} else {
		must(
			xml.NewDecoder(res.Body).
				Decode(&metadata),
		)
	}

	return metadata
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func mustV[T any](r T, err error) T {
	if err != nil {
		panic(err)
	}
	return r
}
