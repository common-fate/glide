import { SearchIcon } from "@chakra-ui/icons";
import { Input, InputGroup, InputLeftElement } from "@chakra-ui/react";
import React from "react";

interface Props {
  placeholderMessage: string;
  setSearchVal?: (string: string) => void;
}

export const CFSearchBar: React.FC<Props> = ({
  placeholderMessage,
  setSearchVal,
}) => {
  return (
    <InputGroup>
      <InputLeftElement pointerEvents="none">
        <SearchIcon color={"neutrals.600"} mx={4} ml={6} />
      </InputLeftElement>
      <Input
        bg="white"
        rounded="xl"
        _hover={
          {
            // borderWidth: "0",
          }
        }
        // _focusWithin={{
        //   boxShadow: "lg",
        // }}
        borderWidth="1px"
        onChange={(e) => setSearchVal?.(e.target.value)}
        placeholder={placeholderMessage}
        sx={{
          _focus: {
            // boxShadow: "outline",
          },
        }}
      />
    </InputGroup>
  );
};
