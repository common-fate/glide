import {
  FormControl,
  FormErrorMessage,
  FormLabel,
  Text,
} from "@chakra-ui/react";
import React from "react";
import { useFormContext } from "react-hook-form";
import { useAdminListTargetGroups } from "../../../../utils/backend-client/admin/admin";

import { AccessRuleFormData } from "../CreateForm";
import { FormStep } from "./FormStep";

export const ProviderStep: React.FC = () => {
  const methods = useFormContext<AccessRuleFormData>();
  const target = methods.watch("targets");
  const targetGroups = useAdminListTargetGroups();
  const isFieldLoading = false; //(!provider && ivp) || (!providerArgs && ivpa);

  return (
    <FormStep
      heading="Provider"
      subHeading="The permissions that the rule gives access to"
      fields={["target", "target.providerId"]}
      // preview={<Preview target={target} provider={provider} />}
      isFieldLoading={isFieldLoading}
    >
      <>
        <FormControl isInvalid={false}>
          <FormLabel htmlFor="target.providerId">
            <Text textStyle={"Body/Medium"}>Provider</Text>
          </FormLabel>
          {/* <ProviderSetupNotice /> */}
          {/* <Controller
            control={methods.control}
            rules={{ required: true }}
            name={"target.providerId"}
            render={({ field: { ref, onChange, ...rest } }) => (
              <ProviderRadioSelector
                onChange={async (t) => {
                  onChange(t);
                  await methods.trigger("target.providerId");
                }}
                {...rest}
              />
            )}
          /> */}

          <FormErrorMessage>Provider is required</FormErrorMessage>
        </FormControl>
      </>
    </FormStep>
  );
};

// type RefreshButtonProps = {
//   providerId: string;
//   argId: string;
// } & Omit<IconButtonProps, "aria-label">;

// export const RefreshButton: React.FC<RefreshButtonProps> = ({
//   argId,
//   providerId,
//   ...props
// }) => {
//   const [loading, setLoading] = useState(false);
//   const { data, mutate, isValidating } = useAdminListTargetGroupArgOptions(
//     providerId,
//     argId
//   );
//   const onClick = async () => {
//     setLoading(true);
//     await mutate(
//       adminListTargetGroupArgOptions(providerId, argId, {
//         refresh: true,
//       })
//     );
//     setLoading(false);
//   };

//   return (
//     <Tooltip>
//       <IconButton
//         {...props}
//         onClick={onClick}
//         isLoading={(!data && isValidating) || loading}
//         icon={<RefreshIcon boxSize="24px" />}
//         aria-label="Refresh"
//         variant={"ghost"}
//       />
//     </Tooltip>
//   );
// };
