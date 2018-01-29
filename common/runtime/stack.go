//
// Because go language does not have industrial level of exception handing mechanism,
// using the information of calling state is the only way to expose secret in code.
//
// Obtain CallerInfo
//
// CallerInfo is the main struct which holds essential information of detail on code.
//
// You can obtain CallerInfo by various functions:
//
// 	GetCallerInfo() - Obtains the caller(the previous calling point to current function)
// 	GetCallerInfoWithDepth() - Obtains the caller of caller by numeric depth
//
package runtime

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"
)

// Information of the position in file and where line number is targeted
type CallerInfo struct {
	PackageName  string
	FileName     string
	Line         int
	FunctionName string

	rawFile string
}

func (c *CallerInfo) String() string {
	return fmt.Sprintf("%s[%s]:%d:%s", c.FileName, c.FunctionName, c.Line, c.PackageName)
}

type CallerStack []*CallerInfo

func (s CallerStack) AsStringStack() []string {
	callerStackString := make([]string, 0)
	for _, caller := range s {
		callerStackString = append(callerStackString, caller.String())
	}

	return callerStackString
}
func (s CallerStack) ConcatStringStack(sep string) string {
	return strings.Join(s.AsStringStack(), sep)
}

// Gets stack of caller info
//
// 0 - The current function
func GetCallerInfoStack(startDepth int, endDepth int) CallerStack {
	callers := make([]*CallerInfo, 0)

	for i := startDepth + 1; i < endDepth+2; i++ {
		callerInfo := GetCallerInfoWithDepth(i)
		if callerInfo == nil {
			break
		}

		callers = append(callers, callerInfo)
	}

	return callers
}

// Gets caller info from current function
func GetCallerInfo() *CallerInfo {
	// Skips
	// 1) this function
	// 2) the function calls this function
	return GetCallerInfoWithDepth(2)
}

func GetCurrentFuncInfo() *CallerInfo {
	// Skips this function
	return GetCallerInfoWithDepth(1)
}

// Gets caller info with number of skipping frames.
//
// 0 - means the current function
// N - means the Nth caller.
func GetCallerInfoWithDepth(countOfSkips int) *CallerInfo {
	pc := make([]uintptr, 1)

	// Skips
	// 1) this function
	// 2) the caller of this function
	n := runtime.Callers(2+countOfSkips, pc)

	if n == 0 {
		return nil
	}

	frame, _ := runtime.CallersFrames(pc).Next()
	return toCallerInfo(&frame)
}

var packageFromFunc = regexp.MustCompile("(.+/(?:\\w|-)+)\\.(.*)$")
var fileFromPath = regexp.MustCompile("[^/]+\\.go$")

func toCallerInfo(frame *runtime.Frame) *CallerInfo {
	finalInfo := &CallerInfo{
		PackageName:  "<N/A>",
		FunctionName: "<N/A>",
		Line:         -1,
		FileName:     "<N/A>",
	}

	if frame.Line > 0 {
		finalInfo.Line = frame.Line
	}

	/**
	 * 1. Extracts package name
	 * 2. Reduction for "/vendor/"
	 */
	matchPackage := packageFromFunc.FindStringSubmatch(frame.Function)
	if len(matchPackage) == 3 {
		finalInfo.PackageName = matchPackage[1]
		finalInfo.FunctionName = matchPackage[2]
	}

	indexOfVendor := strings.Index(finalInfo.PackageName, "/vendor/")
	if indexOfVendor >= 0 {
		finalInfo.PackageName = finalInfo.PackageName[indexOfVendor+8:]
	}
	// :~)

	/**
	 * Extracts file name
	 */
	fileNameMatch := fileFromPath.FindString(frame.File)
	if fileNameMatch != "" {
		finalInfo.FileName = fileNameMatch
	}
	// :~)

	finalInfo.rawFile = frame.File

	return finalInfo
}
