GO:

A small blog application with two parts:

Backend : a JSON API written in Go, storing data in PostgreSQL.

Frontend : a Next.js website that talks to the API.

Visitors can read published posts. Registered users can log in, write posts,
save them as drafts, publish or unpublish them, edit them, and delete them.

# How the two parts fit together

Browser  ->  Next.js frontend (localhost:3000)  ->  Go API (localhost:8080)  ->  PostgreSQL

The browser never calls the Go API directly. Every request goes through the Next.js server (`frontend/src/app/actions.ts`), which forwards it to the Go API and attaches the user's session cookie.

Logging in works like this:

1. The login form sends the email and password to the Next.js server.
2. The Next.js server forwards them to `POST /api/auth/login` on the Go API.
3. The Go API checks the password, creates a session, and returns a session token.
4. The Next.js server stores that token in an `HttpOnly` cookie called `go_blog_session`.
5. On later requests, the Next.js server sends the token back to the Go API,
   which looks up the session to find out who the user is.

# Requirements

- [Go](https://go.dev/doc/install)
- [Node.js](https://nodejs.org/)
- [PostgreSQL](https://www.postgresql.org/download/)

# Setup
1. Create the database tables

Create an empty PostgreSQL database, then run `backend/schema.sql` against it:

psql "postgres://USER:PASSWORD@localhost:5432/DATABASE" -f backend/schema.sql

This creates three tables: `users`, `sessions`, and `posts`.

# Configure the backend
2. Copy `backend/.env.example` to `backend/.env` and fill in your database URL:


   DATABASE_URL=postgres://USER:PASSWORD@localhost:5432/DATABASE - Connection string for PostgreSQL. Required.
   COOKIE_SECURE=false -  Set to `true` in production so cookies are only sent over HTTPS
   FRONTEND_ORIGIN=http://localhost:3000 - The only website allowed to call the API from a browser (CORS).

   Each variable must be on its own line, written as `NAME=value`.


3. Configure the frontend
Copy `frontend/.env.example` to `frontend/.env`:
API_URL=http://localhost:8080 - Address of the Go API. No trailing slash. 
COOKIE_SECURE=false -  Set to `true` in production so cookies are only sent over HTTPS.


The file must be in the `frontend/` folder itself. Next.js does not load
`.env` files from inside `src/`.

4. Install the frontend packages

cd frontend
npm install

# Running the project

Start the API in one terminal:

cd backend
go run .

It prints `API running at http://localhost:8080`.

Start the website in a second terminal:

cd frontend
npm run dev

Then open <http://localhost:3000>.

# Pages

 Address             What it shows                                           Login needed 

 `/`                 All published posts, newest first                         No           
 `/posts/{id}`       One published post                                        No           
 `/register`        Create an account                                          No           
 `/login`           Log in                                                     No           
 `/posts/new`       Write a new post (saved as a draft)                        Yes          
`/posts/mine`      Your posts, with Publish / Unpublish, Edit, Delete          Yes          
 `/posts/{id}/edit`  Edit one of your posts                                    Yes          


# API endpoints

Method    Path                      What it does                         Login needed 
 
 `GET`     `/health`                  Checks the API and database are up    No           
 `POST`    `/api/preview`             Echoes a title and body back          No           
 `POST`    `/api/auth/register`       Creates an account                    No           
 `POST`    `/api/auth/login`          Logs in and sets the session cookie   No          
 `GET`     `/api/auth/me`             Returns the logged-in user            Yes          
 `POST`    `/api/auth/logout`         Deletes the session                   No           
 `GET`     `/api/posts`               Lists published posts (newest 50)     No           
 `GET`     `/api/posts/{id}`          Returns one published post            No           
 `GET`     `/api/me/posts`            Lists your posts, including drafts    Yes          
 `POST`    `/api/posts`               Creates a post (as a draft)           Yes          
 `PUT`     `/api/posts/{id}`          Updates the title and body            Yes          
 `PATCH`   `/api/posts/{id}/publish`  Publishes or unpublishes a post       Yes          
 `DELETE`  `/api/posts/{id}`          Deletes a post                        Yes          

Errors are returned as JSON in the form `{ "error": "message" }`.

# Project structure

backend/ :
  main.go               Starts the server and registers every route.
  
  database.go           Opens the PostgreSQL connection.
  
  schema.sql            Creates the database tables.
  
  cors.go               Decides which website may call the API from a browser.
  
  health.go             GET /health.
  
  handlers.go           POST /api/preview.
  
  auth_json.go          Helpers for reading and writing JSON.
  
  auth_users.go         The User type and email validation.
  
  auth_register.go      POST /api/auth/register.
  
  auth_login.go         POST /api/auth/login.
  
  auth_logout.go        POST /api/auth/logout.
  
  auth_sessions.go      Creates sessions and session cookies.
  
  auth.middleware.go    requireAuth, meHandler, trustedBrowserRequest.
  
  posts.go              The Post type and post validation.
  
  posts_*.go            One file per post endpoint.

frontend/
  src/app/  :
  
    actions.ts          Every call to the Go API (server-only code).
    layout.tsx          Wraps every page (header and page width).
    page.tsx            Home page: published posts.
    login/              Login page.
    register/           Register page.
    posts/new/          New post page.
    posts/mine/         "My posts" page.
    posts/[id]/         One published post.
    posts/[id]/edit/    Edit post page.
    not-found.tsx       Shown when a page or post does not exist.
    error.tsx           Shown when a page fails to load.
    globals.css         Tailwind CSS and theme colors.
    
  src/lib/ :
  
    types.ts            The User, Post, and ActionState types.
    format.ts           Date formatting.
    
  components/ :
    ui/                 Basic building blocks (button, card, input, ...).
    
    *.tsx               Forms and the site header.
