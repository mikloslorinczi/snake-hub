package template

import "html/template"

// Home is the index page
var Home = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Snake-hub</title>
    <style>
        body {
            background-color: rgb(10, 51, 56);
            color :rgb(203, 231, 222);
            font-family: fantasy, sans-serif;
        }
    </style>
</head>
<body>
    <div id="box">
        <h2>Snake-Hub</h2>
        <p>
            A simple server-cient snake game writen in Golang
        </p>
        <h3>How to play</h3>
        <p>
            Visit <a href="https://github.com/mikloslorinczi/snake-hub">the Git page</a>
            and download the binary for Linux, Mac or Windows, or build it yourself...
        </p>
    </div>
</body>
</html>
`))
