package cachesync

import (
	"context"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/service/cachesvc"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

type CacheSyncer struct {
	DB                  ddb.Storage
	AccessHandlerClient types.ClientWithResponsesInterface
	Cache               cachesvc.Service
}

// Sync will attempt to sync all argument options for all providers
// if a particular argument fails to sync, the error is logged and it continues to try syncing the other arguments/providers
func (s *CacheSyncer) Sync(ctx context.Context) error {
	log := logger.Get(ctx)
	q := storage.ListTargetGroups{}
	err := s.DB.All(ctx, &q)
	if err != nil {
		return err
	}
	for _, tg := range q.Result {
		log.Infow("started syncing target group resources cache", "targetgroup", tg)
		err = s.Cache.RefreshCachedTargetGroupResources(ctx, tg)
		if err != nil {
			log.Errorw("failed to refresh resources for targetgroup", "targetgroup", tg, "error", err)
			continue
		}
		log.Infow("completed syncing target group resources cache", "targetgroup", tg)
	}

	// Finally, update targets for requesting access
	return s.Cache.RefreshCachedTargets(ctx)
}

// "lambda execution error: Unhandled: e(**kwargs)\n  File \"/var/task/botocore/client.py\", line 530, in _api_call\n
// return self._make_api_call(operation_name, kwargs)\n  File \"/var/task/botocore/client.py\", line 960, in _make_api_call\n    raise error_class(parsed_response, operation_name)\nbotocore.exceptions.ClientError: An error occurred (Ac
// 	cessDenied) when calling the AssumeRole operation: User: arn:aws:sts::616777145260:assumed-role/cf-handler-common-fate-aws/cf-handler-common-fate-aws is not authorized to perform: sts:AssumeRole on resource: arn:aws:iam::632700053629:role/Granted-Access-Handler-SS-GrantedAccessHandlerSSOR-45YIIRJQGESR[ERROR] ClientError: An error occurred (AccessDenied) when calling the AssumeRole operation: User: arn:aws:sts::616777145260:assumed-role/cf-handler-common-fate-aws/cf-handler-common-fate-aws is not authorized to perform: sts:AssumeRole on resource: arn:aws:iam::632700053629:role/Granted-Access-Handler-SS-GrantedAccessHandlerSSOR-45YIIRJQGESR\nTraceback (most recent call last):\n\u00a0\u00a0File \"/var/task/commonfate_provider/runtime/aws_lambda_entrypoint.py\", line 67, in lambda_handler\n\u00a0\u00a0\u00a0\u00a0return runtime.handle(event, context)\n\u00a0\u00a0File \"/var/task/commonfate_provider/runtime/aws_lambda.py\", line 23, in handle\n\u00a0\u00a0\u00a0\u00a0result = self._do_handle(event=event, context=context)\n\u00a0\u00a0File \"/var/task/commonfate_provider/runtime/aws_lambda.py\", line 91, in _do_handle\n\u00a0\u00a0\u00a0\u00a0tasks._execute(\n\u00a0\u00a0File \"/var/task/commonfate_provider/tasks.py\", line 65, in _execute\n\u00a0\u00a0\u00a0\u00a0return task.run(provider)\n\u00a0\u00a0File \"/var/task/commonfate_provider_dist/provider.py\", line 509, in run\n\u00a0\u00a0\u00a0\u00a0res = p.idstore_client.list_users(\n\u00a0\u00a0File \"/var/task/botocore/client.py\", line 530, in _api_call\n\u00a0\u00a0\u00a0\u00a0return self._make_api_call(operation_name, kwargs)\n\u00a0\u00a0File \"/var/task/botocore/client.py\", line 943, in _make_api_call\n\u00a0\u00a0\u00a0\u00a0http, parsed_response = self._make_request(\n\u00a0\u00a0File \"/var/task/botocore/client.py\", line 966, in _make_request\n\u00a0\u00a0\u00a0\u00a0return self._endpoint.make_request(operation_model, request_dict)\n\u00a0\u00a0File \"/var/task/botocore/endpoint.py\", line 119, in make_request\n\u00a0\u00a0\u00a0\u00a0return self._send_request(request_dict, operation_model)\n\u00a0\u00a0File \"/var/task/botocore/endpoint.py\", line 198, in _send_request\n\u00a0\u00a0\u00a0\u00a0request = self.create_request(request_dict, operation_model)\n\u00a0\u00a0File \"/var/task/botocore/endpoint.py\", line 134, in create_request\n\u00a0\u00a0\u00a0\u00a0self._event_emitter.emit(\n\u00a0\u00a0File \"/var/task/botocore/hooks.py\", line 412, in emit\n\u00a0\u00a0\u00a0\u00a0return self._emitter.emit(aliased_event_name, **kwargs)\n\u00a0\u00a0File \"/var/task/botocore/hooks.py\", line 256, in emit\n\u00a0\u00a0\u00a0\u00a0return self._emit(event_name, kwargs)\n\u00a0\u00a0File \"/var/task/botocore/hooks.py\", line 239, in _emit\n\u00a0\u00a0\u00a0\u00a0response = handler(**kwargs)\n\u00a0\u00a0File \"/var/task/botocore/signers.py\", line 105, in handler\n\u00a0\u00a0\u00a0\u00a0return self.sign(operation_name, request)\n\u00a0\u00a0File \"/var/task/botocore/signers.py\", line 180, in sign\n\u00a0\u00a0\u00a0\u00a0auth = self.get_auth_instance(**kwargs)\n\u00a0\u00a0File \"/var/task/botocore/signers.py\", line 284, in get_auth_instance\n\u00a0\u00a0\u00a0\u00a0frozen_credentials = self._credentials.get_frozen_credentials()\n\u00a0\u00a0File \"/var/task/botocore/credentials.py\", line 610, in get_frozen_credentials\n\u00a0\u00a0\u00a0\u00a0self._refresh()\n\u00a0\u00a0File \"/var/task/botocore/credentials.py\", line 498, in _refresh\n\u00a0\u00a0\u00a0\u00a0self._protected_refresh(is_mandatory=is_mandatory_refresh)\n\u00a0\u00a0File \"/var/task/botocore/credentials.py\", line 514, in _protected_refresh\n\u00a0\u00a0\u00a0\u00a0metadata =
// self._refresh_using()\n\u00a0\u00a0File \"/var/task/botocore/credentials.py\", line 661, in fetch_credentials\n\u00a0\u00a0\u00a0\u00a0return self._get_cached_credentials()\n\u00a0\u00a0File \"/var/task/botocore/credentials.py\", line 671, in _get_cached_credentials\n\u00a0\u00a0\u00a0\u00a0response = self._get_credentials()\n\u00a0\u00a0File \"/var/task/botocore/credentials.py\", line 818, in _get_credentials\n\u00a0\u00a0\u00a0\u00a0return client.assume_role(**kwargs)\n\u00a0\u00a0File \"/var/task/botocore/client.py\", line 530, in _api_call\n\u00a0\u00a0\u00a0\u00a0return self._make_api_call(operation_name, kwargs)\n\u00a0\u00a0File \"/var/task/botocore/client.py\", line 960, in _make_api_call\n\u00a0\u00a0\u00a0\u00a0raise error_class(parsed_response, operation_name)END RequestId...+185 more"
