pipeline {
    agent any

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Docker Compose Deploy') {
            steps {
                // Ensure the network exists so the build doesn't fail
                sh "docker network create paybridge_network || true"
                
                withCredentials([file(credentialsId: 'transaction-svc-env', variable: 'ENV_FILE')]) {
                    sh "cp \$ENV_FILE .env"
                    sh "docker compose up -d --build"
                }
            }
        }

        stage('Verification') {
            steps {
                echo "======== Verifying Services ========"
                sh "docker compose ps"
                // Optional: Check if the application is responding
                // Adjust the port (8083) to match your docker-compose.yml mapping
                //sh "sleep 5 && curl -f http://localhost:8083/health || echo 'Service is up (Health endpoint check bypassed)'
            }
        }
    }

    post {
        success {
            echo "Deployment Successful! Access your service at http://localhost:8082"
        }
        failure {
            echo "Pipeline Failed. Checking logs..."
            sh "docker compose logs --tail=20"
        }
    }
}