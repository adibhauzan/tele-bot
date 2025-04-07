pipeline {
    agent any

    environment {
        IMAGE_NAME = credentials('tele-bot-docker-image')
        WEBHOOK_URL = credentials('tele-bot-webhook-url')
    }


    stages {
        stage('Say Hello') {
            steps {
                echo 'Hello World'
            }
        }

        // stage('Set Environment Variables') {
        //     steps {
        //         script {
        //             echo "GIT_BRANCH: ${env.GIT_BRANCH}"
        //             if (env.GIT_BRANCH == 'origin/main' || env.GIT_BRANCH == 'main' || env.GIT_BRANCH == 'refs/heads/main') {
        //                 env.PORTAINER_WEBHOOK_URL = env.WEBHOOK_PROD_URL
        //                 env.IMAGE_NAME = env.IMAGE_PROD_NAME
        //                 env.STACK_NAME = 'production-stack'
        //             } else if (env.GIT_BRANCH == 'origin/dev' || env.GIT_BRANCH == 'dev' || env.GIT_BRANCH == 'refs/heads/dev') {
        //                 env.PORTAINER_WEBHOOK_URL = env.WEBHOOK_URL
        //                 env.IMAGE_NAME = env.IMAGE_NAME
        //                 env.STACK_NAME = 'dev-stack'
        //             } 

        //             echo "Deploying to stack: ${env.STACK_NAME}"
        //         }
        //     }
        // }

        stage('Docker Build and Push') {
            steps {
                // echo "Build to stack: ${env.STACK_NAME}"
                withDockerRegistry([credentialsId: "docker_pat", url: ""]) {
                    retry(3) {
                        timeout(time: 25, unit: 'MINUTES') {
                            script {
                                sh 'printenv'
                                sh 'docker build -t ${IMAGE_NAME} .'
                                sh 'docker push ${IMAGE_NAME}'
                                sh 'docker image rm ${IMAGE_NAME}'
                            }
                        }
                    }
                }
            }
        }

        // stage('Deploy to Portainer') {
        //     steps {
        //         // echo "Deploying to stack: ${env.STACK_NAME}"
        //         script {
        //             sh """curl -X POST "${env.PORTAINER_WEBHOOK_URL}" """
        //             echo "Deployed to Portainer"
        //         }
        //     }
        // }
    }
}
