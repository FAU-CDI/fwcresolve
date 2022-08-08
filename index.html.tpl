<!DOCTYPE html>
<html lang="en">

<head>
    <title>wdresolve</title>
    <style>
        html {
            font-family: Roboto;
            color: #1a1a1a;
            background-color: #fdfdfd;
        }

        body {
            hyphens: auto;
            text-rendering: optimizeLegibility;
            font-kerning: normal;
        }

        p {
            margin: 1em 0;
            text-align: justify;
        }

        a,
        a:visited {
            color: #0000ff;
        }

        code {
            font-family: monospace;
            color: #ff0000;
        }
    </style>
</head>

<body>
    <h1>wdresolve</h1>
    <p>
        This page implements a resolver for the WissKI Distillery.
        You can either resolve a URI, or check out the list of known prefixes.
    </p>
    <h2>Resolve URI</h2>
    <form action="." method="GET">
        <label for="uri">Enter A URI To Resolve:</label>
        <input type="text" id="uri" name="uri" value="">
        <input type="submit" value="Resolve" />
    </form>
    <h2>Known Prefixes</h2>
    <ul>
        {{range .Prefixes}}
            <li>
                <code>
                    {{ (index . 0) }}
                </code>
                
                &#8658;

                <a
                    href="{{ (index . 1) }}"
                    target="_blank"
                    rel="noopener noreferrer"
                >
                    {{ (index . 1) }}
                </a>
            </li>
        {{end}}
    </ul>
</body>

</html>