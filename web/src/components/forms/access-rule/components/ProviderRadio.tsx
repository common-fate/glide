import { CheckCircleIcon } from "@chakra-ui/icons";
import {
  Box,
  HStack,
  RadioProps,
  Spinner,
  Text,
  Tooltip,
  useRadio,
  useRadioGroup,
  UseRadioGroupProps,
  Wrap,
  WrapItem,
} from "@chakra-ui/react";
import React from "react";

import { ProviderIcon, ShortTypes } from "../../../icons/providerIcon";
import { useAdminListTargetGroups } from "src/utils/backend-client/admin/admin";
import { TargetGroup } from "../../../../utils/backend-client/types";

interface ProviderRadioProps extends RadioProps {
  targetGroup: TargetGroup;
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
        m="1px"
        _checked={{
          m: "0px",
          borderColor: "brandBlue.300",
          borderWidth: "2px",
        }}
        _hover={{
          boxShadow: `0px 0px 0px 1px #2e7fff`,
        }}
        px={6}
        py={5}
        position="relative"
        data-testid={"targetGroup-selector-" + props.targetGroup.id}
      >
        {/* @ts-ignore */}
        {checkbox?.["data-checked"] !== undefined && (
          <CheckCircleIcon
            position="absolute"
            top={2}
            right={2}
            h="12px"
            w="12px"
            color={"brandBlue.300"}
          />
        )}
        <HStack>
          <ProviderIcon shortType={props.targetGroup.icon as ShortTypes} />

          <Text textStyle={"Body/Medium"} color={"neutrals.800"}>
            {props.targetGroup.id}
          </Text>
        </HStack>
      </Box>
    </Box>
  );
};

export const ProviderRadioSelector: React.FC<UseRadioGroupProps> = (props) => {
  const { data } = useAdminListTargetGroups();
  const { getRootProps, getRadioProps } = useRadioGroup(props);
  const group = getRootProps();
  if (!data) {
    return <Spinner />;
  }

  return (
    <Wrap {...group}>
      {data?.targetGroups.map((p) => {
        const radio = getRadioProps({ value: p.id });
        return (
          <WrapItem key={p.id}>
            <ProviderRadio provider={p} {...radio} />
          </WrapItem>
        );
      })}
    </Wrap>
  );
};
