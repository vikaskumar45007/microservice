pipeline {
    agent any

    environment {
        DOCKER_REGISTRY = 'developer45007'
        IMAGE_NAME = 'user-service'
    }

    stages {
        stage('Checkout') {
            steps {
                git branch: 'main', url: 'https://github.com/vikaskumar45007/microservice.git'
            }
        }

        stage('Build') {
            steps {
                sh 'go version'
                sh 'go build ./...'
            }
        }

        stage('Run Tests') {
            steps {
                sh 'go test ./...'
            }
        }

        stage('Build Binary') {
            steps {
                sh 'go build -o user-service ./cmd/user-service/main.go'
            }
        }

        stage('Build & Push Docker Image') {
            steps {
                script {
                    
                    // Compute hash of the binary
                    def newHash = sh(script: "shasum -a 256 ./user-service | awk '{print \$1}'", returnStdout: true).trim()
                    
                    // Load previous hash from file, if it exists
                    def prevHash = ''
                    if (fileExists('last_build.hash')) {
                        prevHash = readFile('last_build.hash').trim()
                    }
                    
                    if (newHash != prevHash) {
                        echo "Binary changed, building and pushing Docker image..."
                        
                        withCredentials([usernamePassword(credentialsId: 'dockerhub-cred', usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS')]) {
                            sh '''
                                echo "$DOCKER_PASS" | docker login -u "$DOCKER_USER" --password-stdin
                                docker build -t "$DOCKER_USER/$IMAGE_NAME:$BUILD_NUMBER" .
                                docker push "$DOCKER_USER/$IMAGE_NAME:$BUILD_NUMBER"
                            '''
                        }
                        
                        // Save new hash for next run
                        writeFile(file: 'last_build.hash', text: newHash)
                    } else {
                        echo "Binary unchanged. Skipping Docker build."
                    }
                }
            }
        }

        stage('Deploy to Kubernetes') {
            steps {
                withCredentials([file(credentialsId: 'kubeconfig', variable: 'KUBECONFIG')]) {
                    sh '''
                        kubectl --kubeconfig="$KUBECONFIG" set image deployment/user-service user-service="$DOCKER_USER/$IMAGE_NAME:$BUILD_NUMBER"
                        kubectl --kubeconfig="$KUBECONFIG" rollout status deployment/user-service
                    '''
                }
            }
        }
    }

    post {
        success {
            echo '✅ Pipeline completed successfully!'
        }
        failure {
            echo '❌ Pipeline failed.'
        }
    }
}
