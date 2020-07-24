# Contents Delivery Network for Map Tiles

Here is a schematic diagram of the service.

![Tile Delivery Service](tile-delivery-service.png)

The following instructions explain how to configure a map vector tile API step by step.
These steps are integrated into an AWS CodePipeline that is described in [`continuous-delivery.md`](continuous-delivery.md).

## Map Tile Generator

A map tile generator is a Lambda function on AWS.

A Lambda function accepts the following parameters,
- `zoom`: zoom level of a tile
- `x`: x-position of a tile
- `y`: y-position of a tile 

## Building API

### Prerequisites

You need the following software installed,
- [AWS CLI](https://aws.amazon.com/cli/?nc1=h_ls)
- [AWS SAM](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/what-is-sam.html)

### Creating an S3 bucket for code

1. Deploy [`api/code-bucket-template.yaml`](api/code-bucket-template.yaml).

    ```
    aws cloudformation deploy --template-file api/code-bucket-template.yaml --stack-name imaginary-map-code-bucket
    ```

   You have to specify an appropriate credential.

2. Remember the S3 bucket name as `CODE_REPOSITORY`.

    ```
    CODE_REPOSITORY=`aws --query "Stacks[0].Outputs[?OutputKey=='CodeRepositoryName']|[0].OutputValue" cloudformation describe-stacks --stack-name imaginary-map-code-bucket | sed 's/^"//; s/"$//'`
    ```

   You have to specify an appropriate credential.

### Creating an S3 bucket for GeoJSON files

1. Deploy [`api/geo-json-bucket-template.yaml`](api/geo-json-bucket-template.yaml).

    ```
    aws cloudformation deploy --template-file api/geo-json-bucket-template.yaml --stack-name imaginary-map-geo-json-bucket
    ```

   You have to specify an appropriate credential.

2. Remember the S3 bucket name for GeoJSON files as `GEO_JSON_BUCKET`.

    ```
    GEO_JSON_BUCKET=`aws --query "Stacks[0].Outputs[?OutputKey=='GeoJsonBucketName']|[0].OutputValue" cloudformation describe-stacks --stack-name imaginary-map-geo-json-bucket | sed 's/^"//; s/"$//'`
    ```

   You have to specify an appropriate credential.

### Uploading GeoJSON files

Suppose the following variable is defined,
- `ISLANDS_GEO_JSON_VERSION`: version of the GeoJSON file for islands
- `PAPERS_GEO_JSON_VERSION`: version of the GeoJSON file for papers

1. Upload a GeoJSON file of islands.

    ```
    aws s3 cp islands.json s3://$GEO_JSON_BUCKET/$ISLANDS_GEO_JSON_VERSION/islands.json
    ```

   You have to specify an appropriate credential.

2. Upload a GeoJSON file of papers.

    ```
    aws s3 cp papers.json s3://$GEO_JSON_BUCKET/$PAPERS_GEO_JSON_VERSION/papers.json
    ```

   You have to specify an appropriate credential.

### Deploying API stack

1. Build Lambda functions.

    ```
    sam build --template api/api-template.yaml
    ```

   Go cannot be built with the `--use-container` option.

2. Package and deploy functions and API.

    ```
    sam deploy --stack-name imaginary-map-api --capabilities CAPABILITY_IAM --s3-bucket $CODE_REPOSITORY --parameter-overrides GeoJsonBucketName=$GEO_JSON_BUCKET
    ```

   You have to specify an appropriate credential.

**NOTE**: when you modify the template, you have to start over from the step 1 even if you have not modified Lambda functions.

### Deploying API Gateway

The CloudFront asks a `develop` stage of the API for new map vector tiles.
You have to manually deploy a `develop` stage of the API to enable the CDN.