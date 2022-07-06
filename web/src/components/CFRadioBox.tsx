import { CheckCircleIcon } from "@chakra-ui/icons";
import { Box, HStack } from "@chakra-ui/layout";
import { RadioProps, useRadio } from "@chakra-ui/react";

export const CFRadioBox: React.FC<RadioProps> = (props) => {
  const {
    getInputProps,
    getCheckboxProps,
    state: { isChecked },
  } = useRadio(props);

  const input = getInputProps();
  const checkbox = getCheckboxProps();

  return (
    <Box as="label">
      <input {...input} />
      <Box
        bg="white"
        {...checkbox}
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
        py={3}
        position="relative"
      >
        {isChecked && (
          <CheckCircleIcon
            position="absolute"
            top={2}
            right={2}
            h="12px"
            w="12px"
            color={"brandGreen.300"}
          />
        )}
        <HStack>{props.children}</HStack>
      </Box>
    </Box>
  );
};
