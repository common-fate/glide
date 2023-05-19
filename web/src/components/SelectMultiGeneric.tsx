import { CloseIcon } from "@chakra-ui/icons";
import {
  Box,
  BoxProps,
  Flex,
  IconButton,
  Input,
  Wrap,
  useBoolean,
  useOutsideClick,
} from "@chakra-ui/react";
import React, { useState } from "react";
import { ShortTypes } from "./icons/providerIcon";

type Props<T, K extends keyof T> = {
  /** The key from `inputArray` item used in lookup */
  keyUsedForFilter: K;
  inputArray: T[];
  selectedItems: T[];
  setSelectedItems: React.Dispatch<React.SetStateAction<T[]>>;
  boxProps: BoxProps;
  /** renderFn takes T and passes it to the react child */
  renderFnTag: (item: T) => React.ReactNode | React.ReactNode[];
  renderFnMenuSelect: (item: T) => React.ReactNode | React.ReactNode[];
  /** onlyOne - limits selection to only one option */
  onlyOne?: boolean;
};

const SelectMultiGeneric = <T, K extends keyof T>({
  inputArray,
  keyUsedForFilter,
  selectedItems,
  setSelectedItems,
  boxProps,
  renderFnMenuSelect,
  renderFnTag,
  onlyOne,
}: Props<T, K>) => {
  const [input, setInput] = useState("");

  const [open, setOpen] = useBoolean(false);

  const inputRef = React.useRef<HTMLInputElement>(null);

  // useMemo filtered entries using Object.entries
  const filteredEntriesForMenu = React.useMemo(() => {
    return (
      inputArray
        // filter out the selected providers based on key match
        .filter((item) => {
          return (
            selectedItems.findIndex(
              (s) => s[keyUsedForFilter] === item[keyUsedForFilter]
            ) === -1
          );
        })
        // dont filter if input is empty; filter input otherwise
        .filter((item) =>
          input === ""
            ? true
            : (item[keyUsedForFilter] as string)
                .toString()
                .toLowerCase()
                .includes(input.toLowerCase())
        )
    );
  }, [input, inputArray, keyUsedForFilter, selectedItems]);

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
    if (filteredEntriesForMenu.length > 0 && input !== "" && !open) {
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
      if (filteredEntriesForMenu.length === 0) return;
      setSelectedItems((prev) =>
        prev.concat(filteredEntriesForMenu[selectedIndex])
      );
      // clear input
      setInput("");
      return;
    }
    // if up key, move selected index up
    if (e.key === "ArrowUp") {
      setSelectedIndex((prev) => {
        if (prev === 0) return filteredEntriesForMenu.length - 1;
        return prev - 1;
      });
      return;
    }
    // if down key, move selected index down
    if (e.key === "ArrowDown") {
      setSelectedIndex((prev) => {
        if (prev === filteredEntriesForMenu.length - 1) return 0;
        return prev + 1;
      });
      return;
    }
    // backspace remove last item in list
    if (e.key === "Backspace") {
      if (input === "") {
        setSelectedItems((prev) => prev.slice(0, prev.length - 1));
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
        borderColor="neutrals.300"
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
          {selectedItems.map((selectedItem) => (
            <Flex
              rounded="full"
              textStyle="Body/Small"
              bg="neutrals.100"
              p={1}
              px={2}
              align="center"
            >
              {renderFnTag(selectedItem)}
              <IconButton
                variant="ghost"
                size="xs"
                h="12px"
                w="12px"
                p={1}
                aria-label="remove item"
                isRound
                icon={<CloseIcon boxSize="8px" h="8px" w="8px" />}
                onClick={() => {
                  setSelectedItems((prev) => {
                    return prev.filter(
                      (s) =>
                        s[keyUsedForFilter] !== selectedItem[keyUsedForFilter]
                    );
                  });
                  setSelectedItems([]);
                }}
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
            placeholder="Search for a provider..."
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
            {filteredEntriesForMenu.map((entryItem) => (
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
                  // redundancy check, ensure not already in selected array
                  setOpen.off();

                  if (
                    selectedItems.find(
                      (item) =>
                        entryItem[keyUsedForFilter] === item[keyUsedForFilter]
                    )
                  )
                    return;
                  setSelectedItems((curr) => [...curr, entryItem]);
                }}
              >
                {renderFnMenuSelect(entryItem)}
              </Flex>
            ))}
          </Box>
        </Box>
      </Box>
    </Box>
  );
};

export default SelectMultiGeneric;
