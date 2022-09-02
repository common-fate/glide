import React from "react";
import { CheckIcon, CopyIcon } from "@chakra-ui/icons";
import {
  useClipboard,
  Stack,
  Code,
  Flex,
  Text,
  Spacer,
  IconButton,
} from "@chakra-ui/react";
import { CodeProps } from "react-markdown/lib/ast-to-react";
export const CodeInstruction: React.FC<CodeProps> = (props) => {
  const { children, node } = props;
  let value = "";
  if (
    node &&
    node.children &&
    node.children.length == 1 &&
    node.children[0].type == "text"
  ) {
    value = node.children[0].value;
  }

  const { hasCopied, onCopy } = useClipboard(value);
  return (
    <Stack>
      <Code
        padding={0}
        bg="white"
        borderRadius="8px"
        borderColor="neutrals.300"
        borderWidth="1px"
      >
        <Flex
          borderColor="neutrals.300"
          borderBottomWidth="1px"
          py="8px"
          px="16px"
          minH="36px"
        >
          <Spacer />
          <IconButton
            variant="ghost"
            h="20px"
            icon={hasCopied ? <CheckIcon /> : <CopyIcon />}
            onClick={onCopy}
            aria-label={"Copy"}
          />
        </Flex>
        <Text
          overflowX="auto"
          color="neutrals.700"
          padding={4}
          whiteSpace="pre-wrap"
        >
          {children}
        </Text>
      </Code>
    </Stack>
  );
};
