package godeps

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// -   github.com/klauspost/compress v1.16.5 // indirect
// -   golang.org/x/image v0.8.0
// +   golang.org/x/exp v0.0.0-20230713183714-613f0c0eb8a1 // indirect

type Version struct {
	Major uint8
	Minor uint8
	Patch uint8
	Date  time.Time
	Sha   string
}

func (v *Version) String() string {
	if v.Sha != "" {
		return fmt.Sprintf("%d.%d.%d-%s-%s", v.Major, v.Minor, v.Patch, v.Date.Format("20060102150405"), v.Sha)
	}
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

var prefixRegex = regexp.MustCompile(`^[+|-]\s+`)

func findNewVersions() (map[string]Upgrade, error) {

	readFile, err := os.Open("./diff")
	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	var dependencies = make(map[string]Upgrade)

	for fileScanner.Scan() {
		var line = fileScanner.Text()
		if strings.HasSuffix(line, "// indirect") {
			continue
		}

		var delete = strings.HasPrefix(line, "-")
		var add = strings.HasPrefix(line, "+")

		if strings.HasPrefix(line, "--- a/go.") || strings.HasPrefix(line, "+++ b/go.") {
			continue

		} else if delete || add {

			var dependency = strings.Split(line, " ")
			var repo = prefixRegex.ReplaceAllString(dependency[0], "")
			var ver, err = parseVersion(dependency[1])
			if err != nil {
				return nil, err
			}

			if exitsingVersion, exists := dependencies[repo]; exists {
				if exitsingVersion.From.IsNewer(ver) {
					if delete {
						dependencies[repo] = Upgrade{From: ver, To: exitsingVersion.To}
					} else {
						dependencies[repo] = Upgrade{From: exitsingVersion.From, To: ver}
					}
				}
			} else {
				if delete {
					dependencies[repo] = Upgrade{From: ver}
				} else {
					dependencies[repo] = Upgrade{To: ver}
				}
			}
		}
	}

	if err := readFile.Close(); err != nil {
		return nil, err
	}

	return dependencies, err
}

func parseVersion(versionStr string) (Version, error) {
	var v Version

	versionStr = strings.TrimSpace(versionStr)
	if string(versionStr[0]) != "v" {
		return Version{}, fmt.Errorf("version does not start with v: %s", versionStr)
	}

	var verArr = strings.Split(versionStr[1:], ".")

	var num, err = strconv.Atoi(string(verArr[0]))
	if err != nil {
		return Version{}, err
	}
	v.Major = uint8(num)

	num, err = strconv.Atoi(string(verArr[1]))
	if err != nil {
		return Version{}, err
	}
	v.Minor = uint8(num)

	// patch is not as simple as the rest as it can contain a timestamp and sha
	// e.g. v0.0.0-20230522175609-2e198f4a06a1
	if strings.Contains(string(verArr[2]), "-") {
		var patchArr = strings.Split(string(verArr[2]), "-")
		num, err = strconv.Atoi(patchArr[0])
		if err != nil {
			return Version{}, err
		}
		v.Patch = uint8(num)

		v.Date, err = time.Parse("20060102150405", patchArr[1])
		if err != nil {
			return Version{}, err
		}

		v.Sha = patchArr[2]
	} else {
		var patchArr = strings.Split(string(verArr[2]), "-")
		num, err = strconv.Atoi(patchArr[0])
		if err != nil {
			return Version{}, err
		}
		v.Patch = uint8(num)
	}

	return v, nil
}

func (v *Version) IsNewer(versionTwo Version) bool {
	if versionTwo.Major > v.Major {
		return true
	} else if versionTwo.Minor > v.Minor {
		return true
	} else if versionTwo.Patch > v.Patch {
		return true
	} else if versionTwo.Date.After(v.Date) {
		return true
	}

	return false
}
