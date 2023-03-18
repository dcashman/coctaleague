<h1> MarketDraft </h1>
Attempt at putting the fantasy football bidding war draft into a real application.

<h2>Dev Instructions</h2>
<h3> Prerequisites </h3>
<li> Install <a href = "https://go.dev/doc/install"> Go </a> </li>
<li> Install <a href = "https://www.postgresql.org/"> postgresql </a> </li>
<li> Install <a href = "https://github.com/golang-migrate/migrate"> migrate </a></li>
<h4> Setup MarketDraft DB </h4>
Connect to the database as the postgres superuser for your system, e.g. postgres, or dcashman:<br/>
<code>sudo -u postgres psql</code><br/>
<br/>
Create the marketdraft DB and connect to it. <br/>
<code>CREATE DATABASE marketdraft;</code><br/>
<code>\c marketdraft;</code><br/>
<br/>
Create a ROLE for the application to connect to the DB, replaceing pa55word with your custom password.<br/>
<code>CREATE ROLE greenlight WITH LOGIN PASSWORD 'pa55word';</code><br/>
<br/>
Add the citext extension and exit.<br/>
<code>CREATE EXTENSION IF NOT EXISTS citext;</code><br/>
<code>exit;</code><br/>
<br/>
Create an environment variable containing the password you created to allow for
local development on your postgresql db, without sharing credentials in source
control.<br/>
<code> export MARKETDRAFT_DB_DSN='postgres://marketdraft:pa55word@localhost/marketdraft?sslmode=disable'</code><br/>
<br/>
Install the go module for interfacing with postgresql.<br/>
<code> go get github.com/lib/pq@v1.10.0 </code><br/>
<br/>
<h4> Setup TLS </h4>
Create dev key and cert in special (hard-coded and ignored by git) tls dir.<br/>
<code>mkdir tls</code><br/>
<code>cd  tls</code><br/>
<code>go run /usr/local/go/src/crypto/tls/generate_cert.go --ecdsa-curve="P256" --host=localhost</code><br/>
<code>cd ..</code><br/>
<br/>
<h4> Misc setup </h4>
Install the go module for creating composable middleware routines.<br/>
<code> go get github.com/justinas/alice@v1 </code><br/>
<br/>
Install the routing package for our endpoints.<br/>
<code> go get github.com/julienschmidt/httprouter@v1.3.0</code><br/>
<br/>
Install the sessions package to keep state across requests.<br/>
<code>go get github.com/golangcollege/sessions@v1</code><br/>
<br/>
Install the latest bcrypt package for password hashing.<br/>
<code>go get golang.org/x/crypto/bcrypt@latest</code><br/>
<br/>
Install the nosurf package to help prevent CSRF<br/>
<code>go get github.com/justinas/nosurf@v1</code><br/>
<br/>
Create a new env variable to store your session secret key.<br/>
<code>openssl rand 32 -base64</code><br/>
<code>export MARKETDRAFT_SESSION_KEY='32 bytes of base64 output from above cmd'</code><br/>
<h3> Running </h3>
<li> 1) At root of repository, run "go run ./cmd/web" </li>
<li> 2) Visit https://localhost:4000 on a browser on the dev machine. </li>
