package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	"github.com/common-fate/granted-approvals/accesshandler/cmd/cli/commands/grants"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/aws-iam-authenticator/pkg/token"
)

func main() {
	flags := []cli.Flag{
		&cli.StringFlag{Name: "api-url", Value: "http://localhost:9092", EnvVars: []string{"ACCESS_HANDLER_URL"}, Hidden: true},
	}

	app := &cli.App{
		Flags:       flags,
		Name:        "granted",
		Usage:       "https://granted.dev",
		UsageText:   "granted [global options] command [command options] [arguments...]",
		HideVersion: false,
		Commands: []*cli.Command{&grants.Command, {Name: "test", Action: func(c *cli.Context) error {
			ctx := c.Context
			cfg, err := config.LoadDefaultConfig(ctx)
			if err != nil {
				return err
			}
			creds, err := cfg.Credentials.Retrieve(ctx)
			if err != nil {
				return err
			}
			if creds.Expired() {
				return errors.New("AWS credentials are expired")
			}
			eksClient := eks.NewFromConfig(cfg)
			r, err := eksClient.DescribeCluster(ctx, &eks.DescribeClusterInput{Name: aws.String("provider-eks-test")})
			if err != nil {
				return err
			}
			pem, err := base64.StdEncoding.DecodeString(aws.ToString(r.Cluster.CertificateAuthority.Data))
			if err != nil {
				return err
			}
			cf := clientcmdapi.Config{
				APIVersion: "v1",
				Kind:       "Config",
				Clusters: map[string]*clientcmdapi.Cluster{
					"cluster": {

						Server:                   aws.ToString(r.Cluster.Endpoint),
						CertificateAuthorityData: pem,
					},
				},
				Contexts: map[string]*clientcmdapi.Context{
					"cluster": {
						Cluster: "cluster",
					},
				},
				CurrentContext: "cluster",
			}
			cc := clientcmd.NewDefaultClientConfig(cf, nil)
			kubeConfig, err := cc.ClientConfig()
			if err != nil {
				return err
			}

			g, err := token.NewGenerator(false, false)
			if err != nil {
				return err
			}
			token, err := g.Get("provider-eks-test")
			if err != nil {
				return err
			}
			kubeConfig.BearerToken = token.Token
			client, err := kubernetes.NewForConfig(kubeConfig)
			if err != nil {
				return err
			}
			res, err := client.RbacV1().Roles("default").List(ctx, v1.ListOptions{})
			if err != nil {
				return err
			}
			_ = res

			awsAuth, err := client.CoreV1().ConfigMaps("kube-system").Get(ctx, "aws-auth", v1.GetOptions{})
			if err != nil {
				return err
			}
			_ = awsAuth
			return nil
		}}},
		EnableBashCompletion: true,
	}

	logCfg := zap.NewDevelopmentConfig()
	logCfg.DisableStacktrace = true

	log, err := logCfg.Build()
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	zap.ReplaceGlobals(log)

	err = app.Run(os.Args)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
}
