package utils

import (
	"regexp"
	"github.com/jfrogdev/jfrog-cli-go/utils/cliutils"
	"strings"
	"errors"
)

// Extract folder path form searchResult by matching the deletePattern.
// At first we replace each * with (.*)
// After finding the matching groups we replace each occurrence of (.*) with the correlated group.
func WilcardToDirsPath(deletePattern, searchResult string) (string, error) {
	if !strings.HasSuffix(deletePattern, "/") {
		return "", errors.New("Delete pattern must end with \"/\"")
	}
	tempPattern := strings.Replace(deletePattern, "**", "*", -1)
	for deletePattern != tempPattern {
		deletePattern = tempPattern
		tempPattern = strings.Replace(deletePattern, "**", "*", -1)
	}

	splitedDeletePattern := strings.Split(deletePattern, "/")
	splitedLen := len(splitedDeletePattern)

	if splitedLen > 1 && strings.Contains(splitedDeletePattern[splitedLen - 2], "*") && strings.Contains(splitedDeletePattern[splitedLen - 3], "*") {
		removedPath := splitedDeletePattern[splitedLen - 2:]
		newDeletePattern := strings.Join(splitedDeletePattern[:splitedLen - 2], "/") + "/"
		del, e := WilcardToDirsPath(newDeletePattern, searchResult)
		if e != nil {
			return "", e
		}
		return WilcardToDirsPath(del + removedPath[0] + "/", searchResult)
	}

	regexpPattern := cliutils.PathToRegExp(deletePattern)
	regexpPattern = strings.Replace(regexpPattern, ".*", "(.*)", -1)
	r, err := regexp.Compile(regexpPattern)
	cliutils.CheckError(err)
	if err != nil {
		return "", err
	}
	groups := r.FindStringSubmatch(searchResult)
	result := deletePattern
	for i := 1; i < len(groups) - 1; i++ {
		// In case the deletePattern ends with * like a/a/a/*/
		// We only need the first level of the matching pattern
		// for example: if the matching result is b/c/d/e/ the c/d/e/ path is redundant.
		if i == (len(groups) - 2) {
			splited := strings.Split(result, "/")
			if strings.Contains(splited[len(splited) - 2], "*") {
				forReplace := strings.Split(groups[i], "/")
				result = strings.Replace(result, "*", forReplace[0], 1)
				continue
			}
			result = strings.Replace(result, "*", groups[i], 1)
		}
		result = strings.Replace(result, "*", groups[i], 1)
	}

	return result, err
}
