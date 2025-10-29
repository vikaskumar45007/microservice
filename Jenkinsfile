pipeline {
    agent any
    environment {
        DOCKER_REGISTRY = 'developer45007'
    }

    stages {
        stage('Checkout') {
            steps {
                git branch: 'main', url: 'https://github.com/vikaskumar45007/microservice.git'
            }
        }

        stage('Detect Services') {
            steps {
                script {
                    // Assign to pipeline-scoped variable (no def)
                    SERVICES = sh(
                        script: "ls -1 cmd | grep -v '@tmp'",
                        returnStdout: true
                    ).trim().split("\n")
                    echo "Detected services: ${SERVICES}"
                }
            }
        }


        stage('Build & Test') {
            steps {
                script {
                    for (service in SERVICES) {
                        dir("${env.WORKSPACE}/cmd/${service}") {
                            echo "🔧 Building & testing ${service}"
                            sh "go test ./..."
                        }
                    }
                }
            }
        }

        stage('Build & Push Docker Images') {
            steps {
                script {
                    for (service in SERVICES) {
                        def SERVICE_NAME = service
                        def dockerfileName = "Dockerfile.${SERVICE_NAME}"

                        dir("${env.WORKSPACE}") {  // docker build from repo root
                            // compute hash
                            def prevHashFile = "${env.WORKSPACE}/last_build_${SERVICE_NAME}.hash"
                            def newHash = sh(
                                script: """
                                    find cmd/${SERVICE_NAME} internal -type f -name '*.go' -print0 \
                                    | xargs -0 sha256sum \
                                    | sha256sum \
                                    | awk '{print \$1}'
                                """, returnStdout: true
                            ).trim()
                            def imageTag = newHash[0..7]
                            def prevHash = fileExists(prevHashFile) ? readFile(prevHashFile).trim() : ''

                            if (newHash != prevHash) {
                                echo "🟢 ${SERVICE_NAME} changed. Building Docker image ${imageTag}"
                                withCredentials([usernamePassword(credentialsId: 'dockerhub-cred', usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS')]) {
                                    env.DOCKER_USER = DOCKER_USER
                                    sh """
                                        echo "$DOCKER_PASS" | docker login -u "$DOCKER_USER" --password-stdin
                                        docker build -t "$DOCKER_USER/$SERVICE_NAME:$imageTag" -f $dockerfileName .
                                        docker push "$DOCKER_USER/$SERVICE_NAME:$imageTag"
                                    """
                                }
                                writeFile(file: prevHashFile, text: newHash)
                            } else {
                                echo "⚪ ${SERVICE_NAME} unchanged. Skipping Docker build."
                            }
                        }
                    }
                }
            }
        }

        stage('Deploy to Kubernetes') {
            steps {
                script {
                    for (service in SERVICES) {
                        def SERVICE_NAME = service
                        def prevHashFile = "${env.WORKSPACE}/last_build_${SERVICE_NAME}.hash"
                        def prevHash = fileExists(prevHashFile) ? readFile(prevHashFile).trim() : ''
                        def imageTag = prevHash.take(8)

                        withCredentials([file(credentialsId: 'kubeconfig', variable: 'KUBECONFIG')]) {
                            sh """
                                kubectl --kubeconfig="$KUBECONFIG" apply -f k8s/${SERVICE_NAME}/
                                kubectl --kubeconfig="$KUBECONFIG" set image deployment/${SERVICE_NAME} ${SERVICE_NAME}=$env.DOCKER_USER/$SERVICE_NAME:$imageTag --record
                                kubectl --kubeconfig="$KUBECONFIG" rollout status deployment/${SERVICE_NAME}
                            """
                        }
                    }
                }
            }
        }
    }

    post {
        success { echo '✅ Pipeline completed successfully!' }
        failure { echo '❌ Pipeline failed.' }
    }
}