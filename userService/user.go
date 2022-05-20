package userService

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
	"wxprojectApiGateway/middleware"
)

type User struct {
	userID   int64
	username string
}

func UserRouter() http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())
	e.POST("/api/login", func(context *gin.Context) {
		code := context.PostForm("code")
		userName := context.PostForm("userName")
		body := strings.NewReader("userName=" + userName)
		response, err := http.Post("http://localhost:8000/user/login/"+code, "application/x-www-form-urlencoded", body)
		if err != nil || response.StatusCode != http.StatusOK {
			context.Status(http.StatusServiceUnavailable)
			return
		}
		if res, err := parseResponse(response); err != nil {
			context.JSON(
				http.StatusServiceUnavailable,
				gin.H{
					"error":  res["error"],
					"token":  "",
					"openid": "",
				},
			)
		} else {
			fmt.Println("token: ", res["token"])
			context.JSON(
				http.StatusOK,
				gin.H{
					"error":  res["error"],
					"token":  res["token"],
					"openid": res["openid"],
				},
			)
		}
	})

	e.Use(middleware.Authorize())

	v1 := e.Group("/api/user")
	{
		v1.POST("/updateUserName", func(context *gin.Context) {
			newUserName := context.PostForm("userName")
			body := strings.NewReader("userName=" + newUserName)
			openid := middleware.GetOpenid(context.GetHeader("Authorization"))
			response, err := http.Post("http://localhost:8000/user/userName/"+openid, "application/x-www-form-urlencoded", body)
			if err != nil || response.StatusCode != http.StatusOK {
				context.Status(http.StatusServiceUnavailable)
			} else {
				fmt.Println("名字更改成功！")
				context.JSON(
					http.StatusOK,
					gin.H{
						"msg": "ok",
					},
				)
			}
		})

		v1.GET("/selectTeamsInfo", func(context *gin.Context) {
			openid := middleware.GetOpenid(context.GetHeader("Authorization"))
			response, err := http.Get("http://localhost:8000/user/teamsInfo/" + openid)
			if err != nil || response.StatusCode != http.StatusOK {
				context.Status(http.StatusServiceUnavailable)
			} else {
				context.JSON(
					http.StatusOK,
					gin.H{
						"msg":   "ok",
						"teams": "",
					},
				)
			}
		})

		v1.GET("/getUserName", func(context *gin.Context) {
			openid := context.Query("openid")
			response, err := http.Get("http://localhost:8000/user/userName/" + openid)
			if err != nil || response.StatusCode != http.StatusOK {
				context.Status(http.StatusServiceUnavailable)
			} else {
				res, _ := parseResponse(response)
				context.JSON(
					http.StatusOK,
					gin.H{
						"msg":      "ok",
						"userName": res["userName"],
					},
				)
			}
		})
	}

	v2 := e.Group("/api/team")
	{
		v2.POST("/createTeam", func(context *gin.Context) {
			teamName := context.PostForm("teamName")
			body := strings.NewReader("teamName=" + teamName)
			openid := middleware.GetOpenid(context.GetHeader("Authorization"))
			response, err := http.Post("http://localhost:8000/team/createTeam/"+openid, "application/x-www-form-urlencoded", body)
			if err != nil || response.StatusCode != http.StatusOK {
				context.Status(http.StatusServiceUnavailable)
			} else {
				res, _ := parseResponse(response)
				context.JSON(
					http.StatusOK,
					gin.H{
						"msg":    "ok",
						"teamID": res["teamID"],
					},
				)
			}
		})

		v2.POST("/updateTeamName", func(context *gin.Context) {
			teamName := context.PostForm("teamName")
			teamID := context.PostForm("teamID")

			body := strings.NewReader("teamName=" + teamName)
			response, err := http.Post("http://localhost:8000/team/updateTeam/"+teamID, "application/x-www-form-urlencoded", body)
			if err != nil || response.StatusCode != http.StatusOK {
				context.Status(http.StatusServiceUnavailable)
			} else {
				context.JSON(
					http.StatusOK,
					gin.H{
						"msg": "ok",
					},
				)
			}
		})

		v2.GET("/selectTeam", func(context *gin.Context) {
			teamID := context.Query("teamID")
			url := "http://localhost:8000/team/selectTeam/" + teamID
			response, err := http.Get(url)
			if err != nil || response.StatusCode != http.StatusOK {
				context.Status(response.StatusCode)
			} else {
				reader := response.Body
				contentLength := response.ContentLength
				contentType := response.Header.Get("Content-Type")
				extraHeaders := map[string]string{}
				context.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
			}
		})

		v2.GET("/selectMembers", func(context *gin.Context) {
			teamID := context.Query("teamID")
			url := "http://localhost:8000/team/selectMembers/" + teamID
			response, err := http.Get(url)
			if err != nil || response.StatusCode != http.StatusOK {
				context.Status(response.StatusCode)
			} else {
				reader := response.Body
				contentLength := response.ContentLength
				contentType := response.Header.Get("Content-Type")
				extraHeaders := map[string]string{}
				context.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
			}
		})

		v2.GET("/selectTeamCode", func(context *gin.Context) {
			teamID := context.Query("teamID")
			openid := middleware.GetOpenid(context.GetHeader("Authorization"))
			url := "http://localhost:8000/team/selectTeamCode/" + teamID + "?userID=" + openid
			response, err := http.Get(url)
			if err != nil || response.StatusCode != http.StatusOK {
				context.Status(response.StatusCode)
			} else {
				reader := response.Body
				contentLength := response.ContentLength
				contentType := response.Header.Get("Content-Type")
				extraHeaders := map[string]string{}
				context.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
			}
		})

		v2.POST("/updateTeamCode", func(context *gin.Context) {
			teamID := context.PostForm("teamID")
			openid := middleware.GetOpenid(context.GetHeader("Authorization"))
			url := "http://localhost:8000/team/updateTeamCode/" + teamID + "?userID=" + openid
			response, err := http.Get(url)
			if err != nil || response.StatusCode != http.StatusOK {
				context.Status(response.StatusCode)
			} else {
				reader := response.Body
				contentLength := response.ContentLength
				contentType := response.Header.Get("Content-Type")
				extraHeaders := map[string]string{}
				context.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
			}
		})
	}
	return e
}

func parseResponse(response *http.Response) (map[string]interface{}, error) {
	var result map[string]interface{}
	body, err := ioutil.ReadAll(response.Body)
	if err == nil {
		err = json.Unmarshal(body, &result)
	}

	return result, err
}
