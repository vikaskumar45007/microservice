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

        // 2️⃣ Detect services (folders under cmd/)
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
                        def SERVICE_NAME = service
                        def dockerfileName = "Dockerfile.${SERVICE_NAME}"

                        // Build binary in a temporary directory inside workspace
                        dir("${env.WORKSPACE}/cmd/${SERVICE_NAME}") {
                            
                            // 3.1 Build binary
                            sh "go build -o ${SERVICE_NAME} ./main.go"

                            // 3.2 Run tests
                            sh "go test ./..."

                            // 3.3 Compute binary hash
                            def newHash = sh(
                                script: "shasum -a 256 ./${SERVICE_NAME} | awk '{print \$1}'",
                                returnStdout: true
                            ).trim()
                            def imageTag = newHash[0..7]

                            // 3.4 Load previous hash
                            def prevHashFile = "${env.WORKSPACE}/last_build_${SERVICE_NAME}.hash"
                            def prevHash = fileExists(prevHashFile) ? readFile(prevHashFile).trim() : ''

                            if (newHash != prevHash) {
                                echo "🟢 ${SERVICE_NAME} changed, building image ${imageTag}"

                                // 3.5 Build & push Docker image (from repo root)
                                dir("${env.WORKSPACE}") {
                                    withCredentials([usernamePassword(credentialsId: 'dockerhub-cred', usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS')]) {
                                        sh """
                                            echo "$DOCKER_PASS" | docker login -u "$DOCKER_USER" --password-stdin
                                            docker build -t "$DOCKER_USER/$SERVICE_NAME:$imageTag" -f $dockerfileName .
                                            docker push "$DOCKER_USER/$SERVICE_NAME:$imageTag"
                                        """
                                    }
                                }

                                // 3.6 Save new hash
                                writeFile(file: prevHashFile, text: newHash)

                                // 3.7 Deploy to Kubernetes
                                withCredentials([file(credentialsId: 'kubeconfig', variable: 'KUBECONFIG')]) {
                                    sh """
                                        kubectl --kubeconfig="$KUBECONFIG" apply -f k8s/${SERVICE_NAME}/
                                        kubectl --kubeconfig="$KUBECONFIG" set image deployment/${SERVICE_NAME} ${SERVICE_NAME}=$DOCKER_USER/$SERVICE_NAME:$imageTag --record
                                        kubectl --kubeconfig="$KUBECONFIG" rollout status deployment/${SERVICE_NAME}
                                    """
                                }

                            } else {
                                echo "⚪ ${SERVICE_NAME} unchanged. Using previous image tag: ${prevHash.take(8)}"

                                // Optional: redeploy last image tag
                                withCredentials([file(credentialsId: 'kubeconfig', variable: 'KUBECONFIG')]) {
                                    sh """
                                        kubectl --kubeconfig="$KUBECONFIG" apply -f k8s/${SERVICE_NAME}/
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
            echo '✅ Pipeline completed successfully!'
        }
        failure {
            echo '❌ Pipeline failed.'
        }
    }
}