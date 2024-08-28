<div align="center">
  
  ![algo-archive](https://github.com/user-attachments/assets/2edf9870-1a5b-4cbc-b224-b0d3dc12642d)
  
</div>

## Table of Contents

- [Overview](#overview)
- [Project Structure](#project-structure)
- [Development](#development)
- [Commit Message Guidelines \& Conventions](#commit-message-guidelines--conventions)
  - [File Naming Conventions](#file-naming-conventions)
- [API Endpoints](#api-endpoints)

## Overview
**AlgoArchive** is a Chrome extension that automatically grabs the LeetCode or HackerRank problem and its solution from the web page and pushes it to GitHub. The extension is designed to help users save and organize their submissions on GitHub, making it easier to track their progress and share their solutions with others.

This repository contains the REST API that the extension and web app use to perform the user-related tasks, like committing to users' repositories and by feeding it solution data.

## Project Structure
```
api.algoarchive.org/
│
├── config                          // Configuration files 
│
├── controllers/                    // Structs designed to house the logic that the routes use
│   └── solutions.go                // Controller for the /solutions route
│
├── forms/                          // Structs templates to bind with request bodies for type checking
│   └── solutions.go                // Forms for the /solutions route
│
├── infra/                          // Infrastructure-related utilities
│   ├── logger                      // Logging class for logging messages to the server. Wrapper around Logrus
│   └── utils/                      // Utility methods for all around the application
│       ├── forms                   // Utils for form binding. Has the method to generate missing properties
│       ├── github                  // Shorthands for interactions with GitHub
│       └── responses               // Shorthands for route responses
│
├── models/                         // Structs for defining Database and form-types
│   └── solutions/                  
│       ├── solution.go             // Form-type for LeetCode submissions
│       └── user.go                 // Form-type for users. **MIGHT BE DEPRECATED SOON**
│
├── routers/                        // Anything to do with API routes
│   ├── middleware/                 // Middleware
│       └── cors.go                 // Utility to write CORS headers
│   ├── routes/                     // Holds functions for routing endpoints under groups
│       └── solutions.go            // Solutions endpoints
│   ├── index.go                    // Main file for routing endpoints
│   └── router.go                   // Router creation
│
└── main.go                         // Server entry point
```

## Development
1. **Setup:** 
    - Clone the repository: 
      ```bash
      git clone https://github.com/AlgoArchiveExt/api.algoarchive.org.git
      cd api.algoarchive.org
      ```
    - Install dependencies:
      ```bash
      go install 
      ```
2. **Run the development server:**
    - Using Air:
      ```bash
      # The air utility lets us to hot-reload the server on filesave for faster and easier development.
      air
      ```
    - Or if you're traditional, use Make:
      ```bash
      # For first-time running
      make
  
      # Use after your first build
      make rebuild

      # Feel free to change the rebuild command to use 'clean' or 'clean-windows', just don't push it please.
      # Any pull requests with it modified off Linux clean will be asked to fix it before merging.
      rebuild: clean-windows build run
      ```
3. **Test the server:**
    - First, see if your terminal gives you any errors
    - Then, go to http://localhost:8080/api/health on your browser to see if your server is running correctly

## Commit Message Guidelines & Conventions

To maintain consistency and clarity in our project’s commit history, please follow these guidelines for commit messages:

- **Type**: Specifies the type of commit being made. Common types include:
  - `feat`: New feature
  - `fix`: Bug fix or change (not as big in scope as a feature)
  - `docs`: Documentation changes
  - `style`: Code style updates (formatting, missing semicolons, etc.)
  - `refactor`: Code changes that neither fix a bug nor add a feature
  - `test`: Adding or modifying tests
  - `chore`: Other changes that do not modify routes or test files (e.g., updates to build scripts)

- **Scope**: Indicates the area or module affected by the commit. For example:
  - `endpoints`
  - `routes`
  - `controllers`
  - `utils`
  
- **Subject**: A concise description of the changes introduced by the commit.

**Format:**

```sh
<type>(<scope>): <subject>
```

**Descriptions**: You may also add a longer description under your commit as long as you separate it with two newlines.

```
feat(endpoints): add routing function for users route

I added a function for routing the users endpoints under /api/v1/users
```

**Examples:**

- `feat(controllers): add users controller`
- `fix(models): create problem schema model`
- `test(utils): add unit tests for route responses`
- `docs: update README with development instructions`
- `chore: update dependencies in package.json`
- `build: update webpack configuration for production build`
- `ci: add GitHub Actions workflow for linting`
- `perf: optimize websocket connection uptime`
- `revert: revert previous commit due to incorrect implementation`
- `merge: merge branch 'feature' into 'main'`
- `deploy: deploy to production server`

For more detailed information on commit message conventions, please refer to [Conventional Commits](https://www.conventionalcommits.org).

**Optional:** You may use emojis to visually represent commit types (e.g., 🔥 for `feat`, 🐛 for `fix`, 📝 for `docs`, etc.).

#### File Naming Conventions

- Use `snake_case` for everything (e.g., `solutions.go`, `form-utils.go`).

## API Endpoints

All of the endpoints are under the ```/api``` directory.

#### List: 
- General:
  - ```GET``` Health ```/health``` Checks if the server is up and routing correctly. Just returns a "live": "ok" response if it's running properly.
- Solutions:
  - ```POST``` Commit Solution ```/v1/solutions``` - Commit a problem solution to a user's repository. As of right now, it doesn't matter if the user doesn't have the AlgoArchive App installed on their repository, we are going to implement this soon.
  - ```GET```  Get all Solutions ```/v1/solutions/:owner/:repo``` - Get all committed solutions from a user's repository.
 




