package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"

	// "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type ClipfyStackProps struct {
	awscdk.StackProps
}

type CognitoOutput struct {
	UserPoolId       *string
	UserPoolClientId *string
}

func createBucket(stack awscdk.Stack) *awss3.Bucket {
	bucket := awss3.NewBucket(stack, jsii.String("clipfy-videos-bucket"), &awss3.BucketProps{
		BucketName: jsii.String("clipfy-videos"),
	})

	return &bucket
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
		UserPoolId:       userPool.UserPoolId(),
		UserPoolClientId: userPoolClient.UserPoolClientId(),
	}
}

func NewClipfyStack(scope constructs.Construct, id string, props *ClipfyStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	createBucket(stack)
	createCognito(stack)

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
