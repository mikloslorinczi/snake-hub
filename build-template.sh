#!/bin/bash

set -euo pipefail

css=`cat www/src/style.css`

js=`cat www/src/index.js`

cat > template/home.go <<- EOM
package template

import "html/template"

// Home is the index page
var Home = template.Must(template.New("").Parse(\`
<!DOCTYPE html>
<html lang="en">
<head>
<title>Snake-hub</title>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<meta http-equiv="X-UA-Compatible" content="ie=edge">
<link href="data:image/x-icon;base64,AAABAAEAEBAQAAEABAAoAQAAFgAAACgAAAAQAAAAIAAAAAEABAAAAAAAgAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAgv+AAARcxwA9Wv8ATGb8AAKoWgACJhUAByCwACE3tQCaqPUAOE7PAAgs/AAAAAAAAAAAAAAAAAAAAAAAAFEVAAAAAAAAURUAAAdwAABRFVUAqocAAFERFQAkqAAGUREVAJs6AAZVVVUACSAAAAZmAGAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABVVQAAAAAAAFEVAAAAAAAAURUAAAAAAABRFQAAAAAAAFEVAAAAAAAAURUACAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAAQAA" rel="icon" type="image/x-icon" />
</head>
<body>
<div class="container"></div>
<script type="text/javascript">
$js
</script>
<style>
$css     
</style>
</body>
</html>
\`))
EOM