# Contents Delivery Network for Map Tiles

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

### Preparing S3 buckets

1. Deploy S3 buckets.

    ```
    aws cloudformation deploy --template-file api/buckets-template.yaml --stack-name imaginary-map-buckets
    ```

   You have to specify an appropriate credential.

2. Remember the S3 bucket name for Lambda function code.

    ```
    CODE_REPOSITORY=`aws --query "Stacks[0].Outputs[?OutputKey=='CodeRepositoryName']|[0].OutputValue" cloudformation describe-stacks --stack-name imaginary-map-buckets | sed 's/^"//; s/"$//'`
    ```

   You have to specify an appropriate credential.

3. Remember the S3 bucket name for GeoJSON files.

You have to redo from the step 1 when you modify the template.

### Uploading GeoJSON files

### Building Lambda functions

1. Build functions.

    ```
    sam build --template api/lambda-functions.yaml
    ```

   [`api/lambda-functions.yaml`](api/lambda-functions.yaml) is used only for building Lambda functions.
   Go cannot be built with the `--use-container` option.

2. Package functions.

    ```
    aws cloudformation package --template-file api/api-template.yaml --s3-bucket $CODE_REPOSITORY --output-template-file api/api-template-packaged.yaml
    ```

   You have to specify an appropriate credential.

3. Depoly functions.

    ```
    aws cloudformation deploy --template-file api/api-template-packaged.yaml --stack-name imaginary-map-api --capabilities CAPABILITY_IAM --parameter-overrides GeoJsonBucketName=$GEO_JSON_BUCKET_NAME
    ```

   You have to specify an appropriate credential.

You have to redo from the step 1 when you modify function code, from the step 2 when you modify the template.