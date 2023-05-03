import { ChevronDownIcon } from "@chakra-ui/icons";
import {
  Menu,
  MenuButton,
  Button,
  MenuList,
  MenuItem,
  useBoolean,
  Input,
  MenuOptionGroup,
  MenuItemOption,
  Box,
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
        value={input}
        onChange={(e) => setInput(e.target.value)}
        onClick={setOpen.on}
        type="text"
      />
      <Menu isOpen={open}>
        {/* <MenuButton as={Button} rightIcon={<ChevronDownIcon />}>
            Your Cats
        </MenuButton> */}
        <MenuList pos="relative" top={0} left={0}>
          <MenuOptionGroup
            position="relative"
            value={selectedProviders}
            type="checkbox"
          >
            {/* Add filtering */}
            {Object.entries(shortTypeValues)
              // dont filter if input is empty
              .filter(([shortType, value]) =>
                input === "" ? true : value.toLowerCase().includes(input)
              )
              .map(([shortType, value]) => (
                <MenuItemOption
                  minH="48px"
                  value={shortType}
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
                </MenuItemOption>
              ))}
          </MenuOptionGroup>
        </MenuList>
      </Menu>
    </Box>
  );
};

export default StaticProviderSelect;
