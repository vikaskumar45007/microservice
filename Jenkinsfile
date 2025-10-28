pipeline {
    agent {
        docker {
        image 'golang:1.22'
        args '-v /var/run/docker.sock:/var/run/docker.sock'
        }
    }

    environment {
        DOCKER_REGISTRY = 'developer45007'
        IMAGE_NAME = 'user-service'
        KUBE_CONFIG = credentials('kubeconfig-id') // Jenkins secret containing kubeconfig
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
                sh 'go build -o user-service main.go'
            }
        }

        stage('Build & Push Docker Image') {
            steps {
                sh """
                docker build -t $DOCKER_REGISTRY/$IMAGE_NAME:\$BUILD_NUMBER .
                docker push $DOCKER_REGISTRY/$IMAGE_NAME:\$BUILD_NUMBER
                """
            }
        }

        stage('Deploy to Kubernetes') {
            steps {
                withKubeConfig([credentialsId: 'kubeconfig-id']) {
                    sh """
                    kubectl set image deployment/user-service user-service=$DOCKER_REGISTRY/$IMAGE_NAME:\$BUILD_NUMBER
                    kubectl rollout status deployment/user-service
                    """
                }
            }
        }
    }

    post {
        success {
            echo 'Pipeline completed successfully!'
        }
        failure {
            echo 'Pipeline failed.'
        }
    }
}
