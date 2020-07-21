# Continuous Delivery

Continuous delivery (CD) is triggered when any commit in this repository is pushed to the master branch.
This is supposed to happen when a pull request is merged to the master branch.

## Creating an S3 Bucket for GeoJSON Files

Please refer to [README.md](README.md#creating-an-s3-bucket-for-geojson-files).

## Uploading GeoJSON Files

Please refer to [README.md](README.md#uploading-geojson-files).

## Creating a GitHub Personal Access Token

[Create one](https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token) and save it in [AWS Secrets Manager](https://aws.amazon.com/secrets-manager/).

Suppose you have the following variable configured,
- `GITHUB_ACCESS_TOKEN`: your generated GitHub personal access token
- `GITHUB_ACCESS_TOKEN_KEY`; e.g., 'GITHUB_ACCESS_TOKEN_KEY=imaginary-map-github-access-token'

1. Create a secret.

    ```
    aws secretsmanager create-secret --name $GITHUB_ACCESS_TOKEN_KEY --secret-string $GITHUB_ACCESS_TOKEN --description 'GitHub personal access token to monitor repositories'
    ```

   You have to provide an appropriate credential.

## Creating a CodePipeline for the API

Suppose you have the following variables configured,
- `GITHUB_REPO_OWNER_NAME`; e.g., `GITHUB_REPO_OWNER_NAME=kikuomax`
- `GITHUB_REPO_NAME`; e.g., `GITHUB_REPO_NAME=kikuomax/imaginary-map`

1. Deploy [`pipeline/pipeline-template.yaml`](api/pipeline-template.yaml).

    ```
    aws cloudformation deploy --template-file pipeline/pipeline-template.yaml --stack-name imaginary-map-codepipeline --capabilities CAPABILITY_IAM --parameter-overrides GeoJsonBucketName=$GEO_JSON_BUCKET GitHubRepoOwnerName=$GITHUB_REPO_OWNER_NAME GitHubRepoName=$GITHUB_REPO_NAME GitHubRepoAccessTokenKey=$GITHUB_ACCESS_TOKEN_KEY
    ```

   You have to provide an appropriate credential.