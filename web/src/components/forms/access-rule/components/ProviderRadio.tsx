import { CheckCircleIcon } from "@chakra-ui/icons";
import {
  Box,
  HStack,
  RadioProps,
  Spinner,
  Text,
  useRadio,
  useRadioGroup,
  UseRadioGroupProps,
} from "@chakra-ui/react";
import React from "react";
import { useListProviders } from "../../../../utils/backend-client/admin/admin";
import { Provider } from "../../../../utils/backend-client/types";
import { ProviderIcon } from "../../../icons/providerIcon";

interface ProviderRadioProps extends RadioProps {
  provider: Provider;
}
const ProviderRadio: React.FC<ProviderRadioProps> = (props) => {
  const { getInputProps, getCheckboxProps } = useRadio(props);

  const input = getInputProps();
  const checkbox = getCheckboxProps();

  return (
    <Box as="label">
      <input {...input} />
      <Box
        {...checkbox}
        bg="white"
        cursor="pointer"
        borderWidth="1px"
        borderRadius="md"
        _checked={{
          borderColor: "brandGreen.300",
          borderWidth: "2px",
        }}
        _focus={{
          boxShadow: "outline",
        }}
        px={6}
        py={5}
        position="relative"
        id="provider-selector"
      >
        {/* @ts-ignore */}
        {checkbox?.["data-checked"] !== undefined && (
          <CheckCircleIcon
            position="absolute"
            top={2}
            right={2}
            h="12px"
            w="12px"
            color={"brandGreen.300"}
          />
        )}
        <HStack>
          <ProviderIcon provider={props.provider.type} />

          <Text textStyle={"Body/Medium"} color={"neutrals.800"}>
            {props.provider.id}
          </Text>
        </HStack>
      </Box>
    </Box>
  );
};

export const ProviderRadioSelector: React.FC<UseRadioGroupProps> = (props) => {
  const { data } = useListProviders();
  const { getRootProps, getRadioProps } = useRadioGroup(props);
  const group = getRootProps();
  if (!data) {
    return <Spinner />;
  }

  return (
    <HStack {...group}>
      {data?.map((p) => {
        const radio = getRadioProps({ value: p.id });
        return <ProviderRadio key={p.id} provider={p} {...radio} />;
      })}
    </HStack>
  );
};
