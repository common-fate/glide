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
  useOutsideClick,
  BoxProps,
} from "@chakra-ui/react";
import React, { useState } from "react";
import {
  ProviderIcon,
  ShortTypes,
  shortTypeValues,
} from "./icons/providerIcon";

type Props = {
  selectedProviders: ShortTypes[];
  setSelectedProviders: React.Dispatch<React.SetStateAction<ShortTypes[]>>;
  boxProps?: BoxProps;
};

/**
 * This is just a static version,
 * In the future we will want to use a version that
 * dyanmicallly fetches the users available providers i.e.
 * `useGetProviders()`
 */
const StaticProviderSelect = ({
  selectedProviders,
  setSelectedProviders,
  boxProps,
}: Props) => {
  //   input state
  const [input, setInput] = useState("");

  const [open, setOpen] = useBoolean(false);

  const inputRef = React.useRef<HTMLInputElement>(null);

  // useMemo filtered entries using Object.entries
  const filteredEntries = React.useMemo(() => {
    return (
      Object.entries(shortTypeValues)
        // filter out the selected providers
        .filter(([shortType, value]) => {
          return !selectedProviders.includes(shortType as ShortTypes);
        })
        // dont filter if input is empty; filter input otherwise
        .filter(([shortType, value]) =>
          input === ""
            ? true
            : value.toLowerCase().includes(input.toLowerCase())
        )
    );
  }, [input, selectedProviders]);

  // menuRef for outside click
  const menuRef = React.useRef<HTMLDivElement>(null);

  useOutsideClick({
    ref: menuRef,
    handler: () => setOpen.off(),
  });

  // selected index
  const [selectedIndex, setSelectedIndex] = React.useState(0);

  // hanlde for selected index:
  // up/down keys moving selected index
  // note: we will add hover suport too
  // enter key press will add the selected index to the selected providers
  // escape key press will close the menu
  const handleSearchInputKeyDown = (e: React.KeyboardEvent) => {
    // if there are results and input and the modal isnt open, open it
    if (filteredEntries.length > 0 && input !== "" && !open) {
      setOpen.on();
      return;
    }

    // if escape key, close menu
    if (e.key === "Escape") {
      setOpen.off();
      return;
    }
    if (e.key === "Enter") {
      // if enter key, add selected index to selected providers
      if (filteredEntries.length === 0) return;
      setSelectedProviders((prev) => [
        ...prev,
        filteredEntries[selectedIndex][0] as ShortTypes,
      ]);
      // clear input
      setInput("");
      return;
    }
    // if up key, move selected index up
    if (e.key === "ArrowUp") {
      setSelectedIndex((prev) => {
        if (prev === 0) return filteredEntries.length - 1;
        return prev - 1;
      });
      return;
    }
    // if down key, move selected index down
    if (e.key === "ArrowDown") {
      setSelectedIndex((prev) => {
        if (prev === filteredEntries.length - 1) return 0;
        return prev + 1;
      });
      return;
    }
    // backspace remove last item in list
    if (e.key === "Backspace") {
      if (input === "") {
        setSelectedProviders((prev) => prev.slice(0, prev.length - 1));
      }
    }
  };

  return (
    <Box {...boxProps}>
      <Input
        as={Flex}
        placeholder="Search for a provider..."
        align="center"
        bg="white !important"
        rounded="md"
        px={0}
        outline="none"
        minH="46px"
        h="min-content"
        cursor="text"
        onClick={(e) => {
          // @ts-ignore
          inputRef.current.focus();
          !open && setOpen.on();
        }}
        // onBlur={setOpen.off}
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
              <IconButton
                variant="ghost"
                size="xs"
                h="12px"
                w="12px"
                p={1}
                aria-label="remove item"
                isRound
                onClick={() => {
                  // if selected, remove
                  if (selectedProviders.includes(shortType as ShortTypes)) {
                    setSelectedProviders(
                      selectedProviders.filter((s) => s !== shortType)
                    );
                    return;
                  }
                }}
                icon={<CloseIcon boxSize="8px" h="8px" w="8px" />}
              />
            </Flex>
          ))}
          <Input
            ref={inputRef}
            h="30px"
            variant="unstyled"
            value={input}
            onKeyDown={handleSearchInputKeyDown}
            onChange={(e) => setInput(e.target.value)}
            onClick={setOpen.on}
            type="text"
            flex="1"
            minW="120px"
            bg="white !important"
            _focusWithin={{
              outline: "none",
            }}
            _focus={{
              outline: "none",
            }}
            _focusVisible={{
              outline: "none",
            }}
            _hover={{
              outline: "none",
            }}
            p={0}
            // onFocus={setOpen.on}
          />
        </Wrap>
      </Input>
      <Box position="relative">
        <Box
          ref={menuRef}
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
            {filteredEntries.map(([shortType, value]) => (
              <Flex
                display="flex"
                // variant="unstyled"
                px={2}
                w="100%"
                minH="48px"
                className="group"
                role="group"
                align="center"
                bg="white"
                // checked={true}
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
