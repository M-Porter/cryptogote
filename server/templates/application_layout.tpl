<!DOCTYPE html>
<html>
    <head>
        <title>CryptoGote</title>
        <link rel="stylesheet" type="text/css" href="https://unpkg.com/normalize.css">
        <link rel="stylesheet" type="text/css" href="/assets/application.css">
        <script type="text/javascript" src="/assets/application.js"></script>
    </head>
    <body>
        <div class="container">
            <div class="header">
                <div class="header-left">
                    <a href="/" class="header-brand">Crypto<span class="header-brand-accent">Go</span>te</a>
                </div>
                <div class="header-right">
                    <div>Encrypted, single-view notes.</div>
                </div>
            </div>
            {{ yield }}
        </div>
    </body>
</html>
