pipeline {
   agent any
    stages {
        stage('Build'){
            steps{
                bat 'go build CAP2IJC.go'
            }
        }
        stage('Test') {
            steps {
                bat 'go test -v CAP2IJC_test.go'
            }
        }
    }
}