package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfront"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfrontorigins"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambdaeventsources"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssns"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	awslambdago "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type ClipfyStackProps struct {
	awscdk.StackProps
}

type CognitoOutput struct {
	UserPool         awscognito.UserPool
	UserPoolClientId *string
}

type BrokerOutput struct {
	Queue awssqs.Queue
	Topic awssns.Topic
}

type StorageOutput struct {
	Bucket awss3.Bucket
	CDN    awscloudfront.Distribution
}

func createBroker(stack awscdk.Stack) *BrokerOutput {
	dlq := awssqs.NewQueue(stack, jsii.String("ClipfyDLQ"), &awssqs.QueueProps{
		QueueName: jsii.String("clipfy-dlq.fifo"),
		Fifo:      jsii.Bool(true),
	})
	queue := awssqs.NewQueue(stack, jsii.String("ClipfyQueue"), &awssqs.QueueProps{
		QueueName:                 jsii.String("clipfy-queue.fifo"),
		Fifo:                      jsii.Bool(true),
		ContentBasedDeduplication: jsii.Bool(true),
		DeadLetterQueue: &awssqs.DeadLetterQueue{
			Queue:           dlq,
			MaxReceiveCount: jsii.Number(3),
		},
	})

	topic := awssns.NewTopic(stack, jsii.String("ClipfyTopic"), &awssns.TopicProps{
		TopicName: jsii.String("clipfy-topic.fifo"),
		Fifo:      jsii.Bool(true),
	})

	// queue policy
	queue.AddToResourcePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions: jsii.Strings("sqs:SendMessage"),
		Effect:  awsiam.Effect_ALLOW,
		Principals: &[]awsiam.IPrincipal{
			awsiam.NewServicePrincipal(jsii.String("sns.amazonaws.com"), nil),
		},
		Resources: jsii.Strings(*queue.QueueArn()),
		Conditions: &map[string]interface{}{
			"StringEquals": map[string]interface{}{
				"aws:SourceArn": topic.TopicArn(),
			},
		},
	}))

	awssns.NewSubscription(stack, jsii.String("ClipfySubscription"), &awssns.SubscriptionProps{
		RawMessageDelivery: jsii.Bool(true),
		Topic:              topic,
		Endpoint:           queue.QueueArn(),
		Protocol:           awssns.SubscriptionProtocol_SQS,
	})

	return &BrokerOutput{
		Queue: queue,
		Topic: topic,
	}
}

func createAPI(stack awscdk.Stack) awslambdago.GoFunction {
	lambda := awslambdago.NewGoFunction(stack, jsii.String("API"), &awslambdago.GoFunctionProps{
		Entry:        jsii.String("cmd/api"),
		FunctionName: jsii.String("clipfy-api"),
		MemorySize:   jsii.Number(256),
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
	})

	functionURL := lambda.AddFunctionUrl(&awslambda.FunctionUrlOptions{
		Cors: &awslambda.FunctionUrlCorsOptions{
			AllowedHeaders: jsii.Strings("*"),
			AllowedOrigins: jsii.Strings("*"),
		},
		AuthType: awslambda.FunctionUrlAuthType_NONE,
	})

	awscdk.NewCfnOutput(stack, jsii.String("APIURL"), &awscdk.CfnOutputProps{
		Value: functionURL.Url(),
	})
	return lambda
}

func createFileProcessingLambda(stack awscdk.Stack) awslambdago.GoFunction {
	lambda := awslambdago.NewGoFunction(stack, jsii.String("FileProcessingLambda"), &awslambdago.GoFunctionProps{
		Entry:        jsii.String("cmd/file_processing"),
		FunctionName: jsii.String("clipfy-file-processing"),
		MemorySize:   jsii.Number(1024),
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Layers: &[]awslambda.ILayerVersion{
			awslambda.LayerVersion_FromLayerVersionArn(stack, jsii.String("ffmpegLayer"), jsii.String("arn:aws:lambda:us-east-1:800352480120:layer:ffmpeg:1")),
		},
	})

	return lambda
}

func createCoginitoLambda(stack awscdk.Stack) awslambdago.GoFunction {
	lambda := awslambdago.NewGoFunction(stack, jsii.String("CognitoLambda"), &awslambdago.GoFunctionProps{
		Entry:        jsii.String("cmd/cognito"),
		FunctionName: jsii.String("clipfy-cognito"),
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
	})

	return lambda
}

func createStorage(stack awscdk.Stack) *StorageOutput {
	bucket := awss3.NewBucket(stack, jsii.String("Bucket"), &awss3.BucketProps{
		BucketName: jsii.String("clipfy-videos"),
		Cors: &[]*awss3.CorsRule{
			{
				AllowedOrigins: jsii.Strings("*"),
				AllowedHeaders: jsii.Strings("*"),
				AllowedMethods: &[]awss3.HttpMethods{
					awss3.HttpMethods_GET,
					awss3.HttpMethods_POST,
					awss3.HttpMethods_PUT,
					awss3.HttpMethods_DELETE,
					awss3.HttpMethods_HEAD,
				},
			},
		},
	})
	oai := awscloudfront.NewOriginAccessIdentity(stack, jsii.String("OAI"), &awscloudfront.OriginAccessIdentityProps{})
	bucket.GrantRead(oai, nil)

	distribution := awscloudfront.NewDistribution(stack, jsii.String("CDN"), &awscloudfront.DistributionProps{
		DefaultBehavior: &awscloudfront.BehaviorOptions{
			Origin: awscloudfrontorigins.NewS3Origin(bucket, &awscloudfrontorigins.S3OriginProps{
				OriginAccessIdentity: oai,
			}),
		},
	})

	return &StorageOutput{
		Bucket: bucket,
		CDN:    distribution,
	}
}

func createCognito(stack awscdk.Stack) *CognitoOutput {
	userPool := awscognito.NewUserPool(stack, jsii.String("UserPool"), &awscognito.UserPoolProps{
		UserPoolName:      jsii.String("clipfy-user-pool"),
		SelfSignUpEnabled: jsii.Bool(true),
		SignInAliases: &awscognito.SignInAliases{
			Email: jsii.Bool(true),
		},
	})

	userPoolClient := awscognito.NewUserPoolClient(stack, jsii.String("UserPoolClient"), &awscognito.UserPoolClientProps{
		UserPool:           userPool,
		UserPoolClientName: jsii.String("clipfy-user-pool-client"),
	})

	awscognito.NewCfnIdentityPool(stack, jsii.String("IdentityPool"), &awscognito.CfnIdentityPoolProps{
		AllowUnauthenticatedIdentities: jsii.Bool(false),
		CognitoStreams:                 &awscognito.CfnIdentityPool_CognitoStreamsProperty{},
		CognitoIdentityProviders: &[]*awscognito.CfnIdentityPool_CognitoIdentityProviderProperty{
			{
				ClientId:     userPoolClient.UserPoolClientId(),
				ProviderName: userPool.UserPoolProviderName(),
			},
		},
	})

	return &CognitoOutput{
		UserPool:         userPool,
		UserPoolClientId: userPoolClient.UserPoolClientId(),
	}
}

func createTable(stack awscdk.Stack) awsdynamodb.TableV2 {
	return awsdynamodb.NewTableV2(stack, jsii.String("Table"), &awsdynamodb.TablePropsV2{
		TableName: jsii.String("clipfy"),
		Billing: awsdynamodb.Billing_OnDemand(&awsdynamodb.MaxThroughputProps{
			MaxReadRequestUnits:  jsii.Number(100),
			MaxWriteRequestUnits: jsii.Number(115),
		}),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("pk"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		SortKey: &awsdynamodb.Attribute{
			Name: jsii.String("sk"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		GlobalSecondaryIndexes: &[]*awsdynamodb.GlobalSecondaryIndexPropsV2{
			{
				IndexName: jsii.String("gsi1"),
				PartitionKey: &awsdynamodb.Attribute{
					Name: jsii.String("gsi1pk"),
					Type: awsdynamodb.AttributeType_STRING,
				},
				SortKey: &awsdynamodb.Attribute{
					Name: jsii.String("gsi1sk"),
					Type: awsdynamodb.AttributeType_STRING,
				},
			},
		},
		TimeToLiveAttribute: jsii.String("ttl"),
	})
}

func NewClipfyStack(scope constructs.Construct, id string, props *ClipfyStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	storage := createStorage(stack)
	cognito := createCognito(stack)
	broker := createBroker(stack)
	table := createTable(stack)
	api := createAPI(stack)
	fileProcessing := createFileProcessingLambda(stack)
	cognitoLambda := createCoginitoLambda(stack)

	cognito.UserPool.AddTrigger(awscognito.UserPoolOperation_PRE_SIGN_UP(), cognitoLambda, awscognito.LambdaVersion_V1_0)
	fileProcessing.AddEventSource(awslambdaeventsources.NewSqsEventSource(broker.Queue, &awslambdaeventsources.SqsEventSourceProps{}))

	broker.Topic.GrantPublish(api)
	broker.Queue.GrantSendMessages(api)
	broker.Queue.GrantConsumeMessages(fileProcessing)
	storage.Bucket.GrantReadWrite(fileProcessing, nil)
	storage.Bucket.GrantWrite(api, nil, nil)
	table.GrantReadWriteData(api)

	api.AddEnvironment(jsii.String("USER_POOL_ID"), cognito.UserPool.UserPoolId(), nil)
	api.AddEnvironment(jsii.String("USER_POOL_CLIENT_ID"), cognito.UserPoolClientId, nil)
	api.AddEnvironment(jsii.String("QUEUE_URL"), broker.Queue.QueueUrl(), nil)
	api.AddEnvironment(jsii.String("TOPIC_ARN"), broker.Topic.TopicArn(), nil)
	api.AddEnvironment(jsii.String("BUCKET_NAME"), storage.Bucket.BucketName(), nil)
	api.AddEnvironment(jsii.String("CDN_URL"), storage.CDN.DomainName(), nil)
	api.AddEnvironment(jsii.String("ENV"), jsii.String("PRODUCTION"), nil)

	fileProcessing.AddEnvironment(jsii.String("QUEUE_URL"), broker.Queue.QueueUrl(), nil)
	fileProcessing.AddEnvironment(jsii.String("TOPIC_ARN"), broker.Topic.TopicArn(), nil)
	fileProcessing.AddEnvironment(jsii.String("BUCKET_NAME"), storage.Bucket.BucketName(), nil)

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewClipfyStack(app, "ClipfyStack", &ClipfyStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return nil
}
