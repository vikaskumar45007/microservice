pipeline {
    agent any

    environment {
        DOCKER_REGISTRY = 'developer45007'
    }

    stages {

        // 1️⃣ Checkout code
        stage('Checkout') {
            steps {
                git branch: 'main', url: 'https://github.com/vikaskumar45007/microservice.git'
            }
        }

        // 2️⃣ Detect all services (folders under cmd/)
        stage('Detect Services') {
            steps {
                script {
                    SERVICES = sh(
                        script: "ls -1 cmd",
                        returnStdout: true
                    ).trim().split("\n")
                    echo "Detected services: ${SERVICES}"
                }
            }
        }

        // 3️⃣ Loop through each service
        stage('Build, Test, Push & Deploy') {
            steps {
                script {
                    for (service in SERVICES) {
                        echo "🔧 Processing service: ${service}"

                        dir("cmd/${service}") {
                            def SERVICE_NAME = service

                            // 3.1 Build binary
                            sh "go build -o ${SERVICE_NAME} ./main.go"

                            // 3.2 Run tests
                            sh "go test ./..."

                            // 3.3 Compute hash of binary
                            def newHash = sh(script: "shasum -a 256 ./${SERVICE_NAME} | awk '{print \$1}'", returnStdout: true).trim()
                            def imageTag = newHash.substring(0, 8)

                            def prevHash = fileExists('last_build.hash') ? readFile('last_build.hash').trim() : ''
                            
                            def dockerfileName = "../../Dockerfile.${SERVICE_NAME}"

                            if (newHash != prevHash) {
                                echo "🟢 Change detected in ${SERVICE_NAME}, building new image: ${imageTag}"

                                // 3.4 Build & push Docker image
                                withCredentials([usernamePassword(credentialsId: 'dockerhub-cred', usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS')]) {
                                    sh """
                                        echo "$DOCKER_PASS" | docker login -u "$DOCKER_USER" --password-stdin
                                        docker build -t "$DOCKER_USER/$SERVICE_NAME:$imageTag" -f $dockerfileName .
                                        docker push "$DOCKER_USER/$SERVICE_NAME:$imageTag"
                                    """
                                }

                                // 3.5 Save new hash
                                writeFile(file: 'last_build.hash', text: newHash)

                                // 3.6 Deploy to Kubernetes
                                withCredentials([file(credentialsId: 'kubeconfig', variable: 'KUBECONFIG')]) {
                                    sh """
                                        kubectl --kubeconfig="$KUBECONFIG" apply -f ../../../k8s/${SERVICE_NAME}/
                                        kubectl --kubeconfig="$KUBECONFIG" set image deployment/${SERVICE_NAME} ${SERVICE_NAME}=$DOCKER_USER/$SERVICE_NAME:$imageTag --record
                                        kubectl --kubeconfig="$KUBECONFIG" rollout status deployment/${SERVICE_NAME}
                                    """
                                }

                            } else {
                                echo "⚪ No change detected for ${SERVICE_NAME}, skipping Docker build."
                                echo "Using previous image tag: ${prevHash.take(8)}"

                                // Optionally redeploy last image tag
                                withCredentials([file(credentialsId: 'kubeconfig', variable: 'KUBECONFIG')]) {
                                    sh """
                                        kubectl --kubeconfig="$KUBECONFIG" apply -f ../../../k8s/${SERVICE_NAME}/
                                        kubectl --kubeconfig="$KUBECONFIG" set image deployment/${SERVICE_NAME} ${SERVICE_NAME}=$DOCKER_USER/$SERVICE_NAME:${prevHash.take(8)} --record
                                        kubectl --kubeconfig="$KUBECONFIG" rollout status deployment/${SERVICE_NAME}
                                    """
                                }
                            }
                        }
                    }
                }
            }
        }
    }

    post {
        success {
            echo '✅ All microservices processed successfully.'
        }
        failure {
            echo '❌ Pipeline failed.'
        }
    }
}