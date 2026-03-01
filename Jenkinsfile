pipeline {
    agent any

    stages {
        stage('Checkout') {
            steps {
                // Pulls the latest code from your GitHub repository
                checkout scm
            }
        }

        stage('Docker Compose Deploy') {
            steps {
                echo "======== Starting Deployment with Docker Compose ========"
                /* --build: Forces Docker to rebuild the image if the Go code changed
                   -d: Runs the containers in the background (detached mode)
                */
                sh "docker compose up -d --build"
            }
        }

        stage('Verification') {
            steps {
                echo "======== Verifying Services ========"
                // Check if the container is running
                sh "docker compose ps"
                
                // Optional: Check if the application is responding
                // Adjust the port (8083) to match your docker-compose.yml mapping
                //sh "sleep 5 && curl -f http://localhost:8083/health || echo 'Service is up (Health endpoint check bypassed)'"
            }
        }
    }

    post {
        success {
            echo "Deployment Successful! Access your service at http://localhost:8081"
        }
        failure {
            echo "Pipeline Failed. Checking logs..."
            // Optional: You can run 'docker compose logs' here to debug
            sh "docker compose logs --tail=20"
        }
    }
}