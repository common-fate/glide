## Providers

The access handler uses a provider framework which allows new services to be integrated easily by implementing a specific interface.

### Adding a provider

Providers are defined in accesshandler/pkg/providers if this provider is for a service which may have multiple possible providers, create a parent folder and then create a specific folder for your new provider. E.g aws/sso or azure/ad otherwise, a single folder is ok e.g okta

One of the simplest ways to get started is to copy an existing provider and update the API calls.
A provider should have the following files:

```
accesshandler
 ┗ pkg
   ┗ providers
     ┗ okta
       ┣ 📜 access.go
       ┣ 📜 errors.go
       ┣ 📜 okta_test.go
       ┣ 📜 okta.go
       ┣ 📜 options.go
       ┗ 📜 validate.go
```

### access.go

This file should contain an implementation for the accessor interface.

```go
type Accessor interface {
	// Grant the access.
	Grant(ctx context.Context, subject string, args []byte) error

	// Revoke the access.
	Revoke(ctx context.Context, subject string, args []byte) error
}
```

It should also contain a struct `Args` always call this `Args` by convention. Args defines what is being requested when granting access. For example, in Okta, a `groupId` is requested. Args should contain all the parameters required by your provider implementation.

The `Args` struct must define `json` tags.

```
type Args struct {
	GroupID string `json:"groupId"`
}
```

The `Grant` and `Revoke` interface provides `args []byte` which will be a serialised json string. It can be unmarshalled into the `Args` type.

```go
// Grant the access by calling the AWS SSO API.
func (p *Provider) Grant(ctx context.Context, subject string, args []byte) error {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}
    ...
}
```

### errors.go

This file should contain named error declarations. For example:

```go
var ErrTargetNotExist error = errors.New("the traget does not exist")
```

### okta_test.go

This is a test file containing testing that works with the testing framework. See [testing](testing.md) for more information.

You may implement any other unit tests that as you see fit.

### okta.go

This file is named after the provider and should contain a struct definition `Provider`. By convention, always name this `Provider`.

The `Provider` struct should contain unexported configuration params required to make API calls etc.

As required, your provider may implement the `Configer`, `Initer` and `ArgSchemarer` interfaces. Full descriptions are available in the [source code](../../accesshandler/pkg/providers/providers.go). These optional interfaces are run to initialise your provider.

If implemented, `Config` is called first, followed by `Init`. Find out more about gconfig and how to use it [here](../backend/gconfig.md).

`ArgSchema` should return `providers.ArgSchema`. This schema is used to render a dynamic form element in the frontend.

```go
type Configer interface {
	Config() gconfig.Config
}
type Initer interface {
	Init(ctx context.Context) error
}
type ArgSchemarer interface {
	ArgSchema() providers.ArgSchema
}
```

```go
type Provider struct {
	client   *okta.Client
	orgURL   string
	apiToken string
}

func (o *Provider) Config() gconfig.Config {
	return gconfig.Config{
		gconfig.String("orgUrl", &o.orgURL, "the Okta organization URL"),
		gconfig.SecretString("apiToken", &o.apiToken, "the Okta API token"),
	}
}

// Init the Okta provider.
func (o *Provider) Init(ctx context.Context) error {
	zap.S().Infow("configuring okta client", "orgUrl", o.orgURL)

	_, client, err := okta.NewClient(ctx, okta.WithOrgUrl(o.orgURL), okta.WithToken(o.apiToken), okta.WithCache(false))
	if err != nil {
		return err
	}
	zap.S().Info("okta client configured")

	o.client = client
	return nil
}

func (p *Provider) ArgSchema() providers.ArgSchema {
	arg := providers.ArgSchema{
		"groupId": {
			Id:          "groupId",
			Title:       "Group",
			FormElement: types.MULTISELECT,
		},
	}

	return arg
}

```

### options.go

The options file should contain an implementation of the `ArgOptioner` interface. This interface accepts an arg key and the response should be the options for that arg. The args are defined in the `Args` struct in access.go

The `Options` type is a label value pair which is rendered in the frontend based on the arg schema.

```go
type ArgOptioner interface {
	Options(ctx context.Context, arg string) ([]types.Option, error)
}

```

### validate.go

The validate file contains an implementation of the `Validator` interface. This interface provides args, which is serialised json object.
As above, unmarshal this into the `Args` struct and validate the parameters. For example, if args contains a groupId, check that the groupId exists in the target. Also check that the user exists in the target.

The subject will be an email address of the user to grant access to.

```go
type Validator interface {
	// Validate arguments and a subject for access without actually granting it.
	Validate(ctx context.Context, subject string, args []byte) error
}

```
