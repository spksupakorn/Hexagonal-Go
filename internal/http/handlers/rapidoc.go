package handlers

import "github.com/labstack/echo/v4"

func RapiDoc(c echo.Context) error {
	const html = `<!doctype html>
<html>
  <head>
    <meta charset="utf-8"/>
    <title>D&D API – RapiDoc</title>
    <meta name="viewport" content="width=device-width, initial-scale=1"/>
    <script type="module" src="https://unpkg.com/rapidoc/dist/rapidoc-min.js"></script>
    <style>
      html, body { height: 100%; margin: 0; }
      rapi-doc { height: 100vh; }
    </style>
  </head>
  <body>
    <rapi-doc
      id="rd"
      spec-url="/openapi.json"
	  server-url="http://localhost:8080/api/v1"
      render-style="read"
      show-header="true"
      theme="light"
      nav-item-spacing="compact"
      allow-try="true"
      show-curl-before-try="true"
      schema-style="table"
      sort-endpoints-by="path"
      sort-tags="true"
      use-path-in-nav-bar="true"
    >
      <div slot="nav-logo" style="font-weight:700;padding:8px 12px">D&D API</div>
      <div slot="header" style="padding:8px 12px">D&D Characters & Quests</div>
    </rapi-doc>
  </body>
</html>`
	return c.HTML(200, html)
}
