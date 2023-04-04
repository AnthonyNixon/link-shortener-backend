package link

import (
	"errors"
	data "github.com/anthonynixon/link-shortener-backend/internal/cloud"
	"github.com/anthonynixon/link-shortener-backend/internal/shortcode"
	"github.com/anthonynixon/link-shortener-backend/internal/types"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func AddLinkV1(router *gin.Engine) {
	router.GET("/:short", RedirectToLink)
	router.GET("/link/:short", GetLongLink)
	router.POST("/link", CreateShortLink)
}

func getLinkDetails(short string) (link types.Link, err error) {
	link, err = data.GetLink(short)
	return
}

func GetLongLink(c *gin.Context) {
	short := c.Param("short")
	link, err := getLinkDetails(short)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, link)
}

func RedirectToLink(c *gin.Context) {
	short := c.Param("short")
	short = strings.ToUpper(short)
	link, err := getLinkDetails(short)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "That link doesn't exist"})
	}

	c.Redirect(http.StatusFound, link.Long)
}

func CreateShortLink(c *gin.Context) {
	var newLink types.Link
	err := c.ShouldBindJSON(&newLink)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if newLink.Short == "" {
		newLink.Short = shortcode.New()
	}
	
	newLink.Created = time.Now().Unix()

	err = data.NewLink(newLink)
	if err != nil {
		if errors.Is(err, data.AlreadyExistsErr) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, newLink)
}
