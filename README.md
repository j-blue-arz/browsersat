Solving SAT in the browser.

[Gophersat](https://github.com/crillab/gophersat) compiled to WebAssembly, with a React client.

## Run
    docker-compose up

Open `localhost` in your browser.

See [Gophersat documentation](https://pkg.go.dev/github.com/crillab/gophersat@v1.3.1/bf#Parse) for syntax of clauses.


## Roadmap
- [x] User can add clauses
- [x] User can compute SAT
- [x] Docker build
- [x] Display SAT assignment
- [x] User can flip literals
- [x] reduce binary size with TinyGo
- [ ] Ship it

### Other ideas
* Validate input
* Support pseudo-boolean logic
* SVG display to increase readability
* Show evaluation of subformulas
* Create UI for a specific use-case, e.g. some planning, scheduling or data science problem - or just simply Sudoku.
* ...
