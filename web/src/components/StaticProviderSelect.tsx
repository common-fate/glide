import { CheckIcon, ChevronDownIcon, CloseIcon } from "@chakra-ui/icons";
import {
  Button,
  useBoolean,
  Input,
  Box,
  Stack,
  Flex,
  Text,
  IconButton,
  Wrap,
} from "@chakra-ui/react";
import React, { useState } from "react";
import {
  ProviderIcon,
  ShortTypes,
  shortTypeValues,
} from "./icons/providerIcon";

type Props = {};

/**
 * This is just a static version,
 * In the future we will want to use a version that
 * dyanmicallly fetches the users available providers i.e.
 * `useGetProviders()`
 */
const StaticProviderSelect = (props: Props) => {
  // controlled input for menu

  //   input state
  const [input, setInput] = useState("");

  const [selectedProviders, setSelectedProviders] = useState<ShortTypes[]>([]);

  const [open, setOpen] = useBoolean(false);

  return (
    <Box>
      <Input
        as={Flex}
        align="center"
        bg="white !important"
        rounded="md"
        px={0}
        outline="none"
        minH="46px"
        h="min-content"
      >
        {/* result preview box */}
        <Wrap px={2}>
          {selectedProviders.map((shortType) => (
            <Flex
              rounded="full"
              textStyle="Body/Small"
              bg="neutrals.100"
              p={1}
              px={2}
              align="center"
            >
              <ProviderIcon mr={1} shortType={shortType} />
              {shortType}
              {/* <Text>{shortType}</Text> */}
              <IconButton
                variant="ghost"
                size="xs"
                h="12px"
                w="12px"
                p={1}
                aria-label="remove item"
                isRound
                icon={<CloseIcon boxSize="8px" h="8px" w="8px" />}
              />
            </Flex>
          ))}
          <Input
            h="30px"
            variant="unstyled"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onClick={setOpen.on}
            type="text"
            flex="1"
            minW="120px"
            bg="white !important"
            p={0}
          />
        </Wrap>
      </Input>
      <Box position="relative">
        <Box
          display={open ? "block" : "none"}
          pos="absolute"
          top={2}
          left={0}
          rounded="md"
          borderColor="neutrals.300"
          borderWidth="1px"
          py={2}
          bg="white"
          w="100%"
          zIndex={2}
        >
          <Box>
            {Object.entries(shortTypeValues)
              // filter out the selected providers
              .filter(([shortType, value]) => {
                return !selectedProviders.includes(
                  shortType as ShortTypes
                ) as ShortTypes;
              })
              // dont filter if input is empty; filter input otherwise
              .filter(([shortType, value]) =>
                input === "" ? true : value.toLowerCase().includes(input)
              )
              .map(([shortType, value]) => (
                <Flex
                  display="flex"
                  // variant="unstyled"
                  px={2}
                  w="100%"
                  minH="48px"
                  className="group"
                  role="group"
                  value={shortType}
                  align="center"
                  bg="white"
                  checked={true}
                  data-checkbox={true}
                  aria-checked={true}
                  _hover={{
                    bg: "neutrals.100",
                  }}
                  _selected={{
                    bg: "neutrals.100",
                  }}
                  onClick={() => {
                    // if selected, remove
                    if (selectedProviders.includes(shortType as ShortTypes)) {
                      setSelectedProviders(
                        selectedProviders.filter((s) => s !== shortType)
                      );
                      return;
                    }
                    setSelectedProviders([
                      ...selectedProviders,
                      shortType as ShortTypes,
                    ]);
                  }}
                >
                  <ProviderIcon
                    shortType={shortType as ShortTypes}
                    key={shortType}
                    mr={2}
                  />
                  <span>{value}</span>
                  <CheckIcon
                    display="none"
                    _groupChecked={{
                      display: "block",
                    }}
                  />
                </Flex>
              ))}
          </Box>
        </Box>
      </Box>
    </Box>
  );
};

export default StaticProviderSelect;
