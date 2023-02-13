import React from "react";
import { CheckIcon, CopyIcon } from "@chakra-ui/icons";
import { CodeProps } from "react-markdown/lib/ast-to-react";
import {
  useClipboard,
  Stack,
  Code,
  CodeProps as CProps,
  Flex,
  Text,
  Spacer,
  IconButton,
} from "@chakra-ui/react";

export const CFCode: React.FC<CProps> = (props) => {
  return (
    <Code
      padding={0}
      bg="white"
      borderRadius="8px"
      borderColor="neutrals.300"
      borderWidth="1px"
    >
      <Text color="neutrals.700" px={1} whiteSpace="pre-wrap">
        {props.children}
      </Text>
    </Code>
  );
};

export interface CFCodeMultilineProps {
  text: string;
}

/**
 * A multiline preformatted code block with a copy button.
 */
export const CFCodeMultiline: React.FC<CFCodeMultilineProps> = ({ text }) => {
  const { hasCopied, onCopy } = useClipboard(text);

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
          {text}
        </Text>
      </Code>
    </Stack>
  );
};

/**
 * This component should only be used with react-markdown.
 *
 * Use CFMultilineCode if you're working with regular React components
 * rather than markdown.
 */
export const CFReactMarkownCode: React.FC<CodeProps> = (props) => {
  const { children, node, ...rest } = props;
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

  // if the code is inline should show in same line.
  if (props?.inline) {
    return <CFCode children={value} />;
  }

  return (
    <Stack>
      <Code
        padding={0}
        bg="white"
        borderRadius="8px"
        borderColor="neutrals.300"
        borderWidth="1px"
        {...rest}
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
