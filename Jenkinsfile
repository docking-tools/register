
node {
stage 'checkout'
  checkout scm
stage 'build'
  docker.image("").inside {
    stage 'build register'
    go get github.com/stretchr/testify/assert
    go test -v -cover ./...
    go install
  }
  stash includes: bin/*, name: binary 
  


}
