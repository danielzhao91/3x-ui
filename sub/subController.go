package sub

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/mhsanaei/3x-ui/v2/config"

	"github.com/gin-gonic/gin"
)

// SUBController handles HTTP requests for subscription links and JSON configurations.
type SUBController struct {
	subTitle       string
	subPath        string
	subJsonPath    string
	subClashPath   string
	jsonEnabled    bool
	clashEnabled   bool
	subEncrypt     bool
	updateInterval string

	subService     *SubService
	subJsonService *SubJsonService
	subClashService *SubClashService
}

// NewSUBController creates a new subscription controller with the given configuration.
func NewSUBController(
	g *gin.RouterGroup,
	subPath string,
	jsonPath string,
	clashPath string,
	jsonEnabled bool,
	clashEnabled bool,
	encrypt bool,
	showInfo bool,
	rModel string,
	update string,
	jsonFragment string,
	jsonNoise string,
	jsonMux string,
	jsonRules string,
	clashRuleSet string,
	subTitle string,
) *SUBController {
	sub := NewSubService(showInfo, rModel)
	a := &SUBController{
		subTitle:       subTitle,
		subPath:        subPath,
		subJsonPath:    jsonPath,
		subClashPath:   clashPath,
		jsonEnabled:    jsonEnabled,
		clashEnabled:   clashEnabled,
		subEncrypt:     encrypt,
		updateInterval: update,

		subService:      sub,
		subJsonService:  NewSubJsonService(jsonFragment, jsonNoise, jsonMux, jsonRules, sub),
		subClashService: NewSubClashService(clashRuleSet, sub),
	}
	a.initRouter(g)
	return a
}

// initRouter registers HTTP routes for subscription links and JSON endpoints
// on the provided router group.
func (a *SUBController) initRouter(g *gin.RouterGroup) {
	gLink := g.Group(a.subPath)
	gLink.GET(":subid", a.subs)
	if a.jsonEnabled {
		gJson := g.Group(a.subJsonPath)
		gJson.GET(":subid", a.subJsons)
	}
	if a.clashEnabled {
		gJson := g.Group(a.subClashPath)
		gJson.GET(":subid", a.subClashes)
	}
}

// subs handles HTTP requests for subscription links, returning either HTML page or base64-encoded subscription data.
func (a *SUBController) subs(c *gin.Context) {
	subId := c.Param("subid")
	scheme, host, hostWithPort, hostHeader := a.subService.ResolveRequest(c)
	subs, lastOnline, traffic, err := a.subService.GetSubs(subId, host)
	if err != nil || len(subs) == 0 {
		c.String(400, "Error!")
	} else {
		result := ""
		for _, sub := range subs {
			result += sub + "\n"
		}

		// If the request expects HTML (e.g., browser) or explicitly asked (?html=1 or ?view=html), render the info page here
		accept := c.GetHeader("Accept")
		if strings.Contains(strings.ToLower(accept), "text/html") || c.Query("html") == "1" || strings.EqualFold(c.Query("view"), "html") {
			// Build page data in service
			subURL, subJsonURL, subClashURL := a.subService.BuildURLs(scheme, hostWithPort, a.subPath, a.subJsonPath, a.subClashPath, subId)
			if !a.jsonEnabled {
				subJsonURL = ""
			}
			if !a.clashEnabled {
				subClashURL = ""
			}
			// Get base_path from context (set by middleware)
			basePath, exists := c.Get("base_path")
			if !exists {
				basePath = "/"
			}
			// Add subId to base_path for asset URLs
			basePathStr := basePath.(string)
			if basePathStr == "/" {
				basePathStr = "/" + subId + "/"
			} else {
				// Remove trailing slash if exists, add subId, then add trailing slash
				basePathStr = strings.TrimRight(basePathStr, "/") + "/" + subId + "/"
			}
			page := a.subService.BuildPageData(subId, hostHeader, traffic, lastOnline, subs, subURL, subJsonURL, subClashURL, basePathStr)
			c.HTML(200, "subpage.html", gin.H{
				"title":        "subscription.title",
				"cur_ver":      config.GetVersion(),
				"host":         page.Host,
				"base_path":    page.BasePath,
				"sId":          page.SId,
				"download":     page.Download,
				"upload":       page.Upload,
				"total":        page.Total,
				"used":         page.Used,
				"remained":     page.Remained,
				"expire":       page.Expire,
				"lastOnline":   page.LastOnline,
				"datepicker":   page.Datepicker,
				"downloadByte": page.DownloadByte,
				"uploadByte":   page.UploadByte,
				"totalByte":    page.TotalByte,
				"subUrl":       page.SubUrl,
				"subJsonUrl":   page.SubJsonUrl,
				"subClashUrl":   page.SubClashUrl,
				"result":       page.Result,
			})
			return
		}

		// Add headers
		header := fmt.Sprintf("upload=%d; download=%d; total=%d; expire=%d", traffic.Up, traffic.Down, traffic.Total, traffic.ExpiryTime/1000)
		a.ApplyCommonHeaders(c, header, a.updateInterval, a.subTitle)

		if a.subEncrypt {
			c.String(200, base64.StdEncoding.EncodeToString([]byte(result)))
		} else {
			c.String(200, result)
		}
	}
}

// subJsons handles HTTP requests for JSON subscription configurations.
func (a *SUBController) subJsons(c *gin.Context) {
	subId := c.Param("subid")
	_, host, _, _ := a.subService.ResolveRequest(c)
	jsonSub, header, err := a.subJsonService.GetJson(subId, host)
	if err != nil || len(jsonSub) == 0 {
		c.String(400, "Error!")
	} else {

		// Add headers
		a.ApplyCommonHeaders(c, header, a.updateInterval, a.subTitle)

		c.String(200, jsonSub)
	}
}

func (a *SUBController) subClashes(c *gin.Context) {
	subId := c.Param("subid")
	_, host, _, _ := a.subService.ResolveRequest(c)
	clashSub, header, err := a.subClashService.GetClash(subId, host)
	if err != nil || len(clashSub) == 0 {
		c.String(400, "Error!")
	} else {

		// Add headers
		a.ApplyCommonHeaders(c, header, a.updateInterval, a.subTitle)

		c.String(200, clashSub)
	}
}

// ApplyCommonHeaders sets common HTTP headers for subscription responses including user info, update interval, and profile title.
func (a *SUBController) ApplyCommonHeaders(c *gin.Context, header, updateInterval, profileTitle string) {
	c.Writer.Header().Set("Subscription-Userinfo", header)
	c.Writer.Header().Set("Profile-Update-Interval", updateInterval)
	c.Writer.Header().Set("Profile-Title", "base64:"+base64.StdEncoding.EncodeToString([]byte(profileTitle)))
}
