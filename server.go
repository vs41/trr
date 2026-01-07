// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

//go:build !js
// +build !js

package main

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/pion/logging"
)

var (
	indexTemplate = &template.Template{}
	log           = logging.NewDefaultLoggerFactory().NewLogger("sfu-fiber")
)

func main() {
	// Read index.html
	indexHTML, err := os.ReadFile("index.html")
	if err != nil {
		panic(err)
	}
	indexTemplate = template.Must(template.New("").Parse(string(indexHTML)))

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://astt.live",
	}))
	// Serve index.html
	app.Get("/", func(c *fiber.Ctx) error {
		unitID := c.Query("unitID", "1") // default = 1
		teamID := c.Query("teamID", "1") // default = 1
		username := c.Query("username", "user"+unitID)

		wsURL := fmt.Sprintf(
			"wss://astt.live/voipsocket?gameID=28&teamID=%s&unitID=%s&encyptionOn=0&username=%s&radioRange=40",
			teamID, unitID, username,
		)
		fmt.Println("WebSocket URL:", wsURL)

		var buf bytes.Buffer
		if err := indexTemplate.Execute(&buf, wsURL); err != nil {
			return err
		}

		c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
		return c.Send(buf.Bytes())
	})

	// WebSocket handler
	// app.Get("/websocket", websocket.New(sever.WebsocketHandler))

	// // Periodically request keyframes
	// go func() {
	// 	for range time.NewTicker(time.Second * 3).C {
	// 		sever.DispatchAllKeyFrames()
	// 	}
	// }()

	// tlsCert := "cert.pem"
	// tlsKey := "key.pem"
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := app.Listen(":" + port); err != nil {
		log.Errorf("Failed to start Fiber server: %v", err)
	}
	// if err := app.ListenTLS(":"+port, tlsCert, tlsKey); err != nil {
	// 	log.Errorf("Failed to start Fiber server: %v", err)
	// }
}
