import {
  Box,
  Button,
  Center,
  CenterProps,
  chakra,
  Flex,
  HStack,
  Input,
  Stack,
  Text,
  useBoolean,
  useEventListener,
  useToast,
} from "@chakra-ui/react";
import debounce from "lodash.debounce";
import { useMemo, useRef, useState } from "react";
import { FixedSizeList as ReactWindowList } from "react-window";
import {
  CreatePreflightRequestBody,
  TargetField,
} from "../utils/backend-client/types";
import { Command as CommandNew } from "../utils/cmdk";
import { ProviderIcon, ShortTypes } from "./icons/providerIcon";
// @ts-ignore
import axios from "axios";
import { useNavigate } from "react-location";
import {
  userRequestPreflight,
  useUserListEntitlements,
} from "../utils/backend-client/default/default";
import { useTargets } from "../utils/context/targetsContext";
import { TargetDetail } from "./Target";
const IS_MAC = /(Mac|iPhone|iPod|iPad)/i.test(
  navigator.userAgent || navigator.platform || "unknown"
);
const SELECTED_KEY = "__selected_key__";
const TARGET_HEIGHT = 64;
const TARGETS = 5;
const ACTION_KEY_DEFAULT = ["Ctrl", "Control"];
const ACTION_KEY_APPLE = ["⌘", "Command"];
const StyledCommandList = chakra(CommandNew.List);

interface EntitlementCheckoutProps {
  checked: Set<string>;
  setChecked: React.Dispatch<React.SetStateAction<Set<string>>>;
}

export const EntitlementCheckout: React.FC<EntitlementCheckoutProps> = ({
  checked,
  setChecked,
}) => {
  return <Search checked={checked} setChecked={setChecked} />;
};

// eslint-disable-next-line @typescript-eslint/no-empty-interface
interface SearchProps {
  checked: Set<string>;
  setChecked: React.Dispatch<React.SetStateAction<Set<string>>>;
}
const Search: React.FC<SearchProps> = ({ checked, setChecked }) => {
  // const [checked, setChecked] = useState<Set<string>>(new Set());
  const [searchValue, setSearchValue] = useState<string>("");
  const searchInputRef = useRef<HTMLInputElement>(null);
  const [actionKey] = useState<string[]>(
    IS_MAC ? ACTION_KEY_APPLE : ACTION_KEY_DEFAULT
  );
  const [submitLoading, submitLoadingToggle] = useBoolean();

  const { targets, validating: targetsValidating } = useTargets();
  const navigate = useNavigate();
  const toast = useToast();

  const handleSubmit = async () => {
    try {
      const preflightRequest: CreatePreflightRequestBody = {
        targets: targets
          .filter((t) => checked.has(t.id.toLowerCase()))
          .map((t) => t.id),
      };
      submitLoadingToggle.on();
      const preflightResponse = await userRequestPreflight(preflightRequest);
      navigate({ to: `/preflight/${preflightResponse.id}` });
    } catch (err) {
      let description: string | undefined;
      if (axios.isAxiosError(err)) {
        // @ts-ignore
        description = err?.response?.data.error;
      }
      toast({
        title: "Error submitting request",
        description,
        status: "error",
        variant: "subtle",
        duration: 2200,
        isClosable: true,
      });
    } finally {
      submitLoadingToggle.off();
    }
  };
  // Watch keys for cmd Enter submit
  useEventListener("keydown", (event) => {
    const hotkey = IS_MAC ? "metaKey" : "ctrlKey";
    if (event?.key?.toLowerCase() === "enter" && event[hotkey]) {
      event.preventDefault();
      checked.size > 0 && handleSubmit();
    }
  });

  const targetFieldsToString = (targetFields: TargetField[]): string => {
    return targetFields
      .map(
        (targetField) =>
          targetField.valueLabel + ";" + targetField.valueDescription
      )
      .join(";")
      .toLowerCase();
  };

  const onDisplay = useMemo(() => {
    if (targets.length === 0) return targets;
    if (searchValue === "") return targets;
    if (searchValue === SELECTED_KEY)
      return targets.filter((t) => checked.has(t.id.toLowerCase()));
    return targets.filter((target) => {
      // you can use a space between elements of your search, and it will filter by matching on those elements
      const key = target.id.toLowerCase() + targetFieldsToString(target.fields);
      const searchValues = searchValue.toLowerCase().split(" ");
      for (let i = 0; i < searchValues.length; i++) {
        if (!key.includes(searchValues[i])) {
          return false;
        }
      }
      return true;
    });
  }, [targets, searchValue, checked]);

  const debouceSearchInput = debounce((value) => {
    if (value === SELECTED_KEY) {
      setSearchValue(SELECTED_KEY);
    } else {
      setSearchValue(value);
    }
  }, 500);

  const onShowSelected = () => {
    if (searchInputRef.current?.value !== undefined) {
      searchInputRef.current.value = "";
    }
    debouceSearchInput(SELECTED_KEY);
  };
  // using a ref to set the value here to avoid a react rerender when the input is updated
  const onSetSearch = (value: string) => {
    if (searchInputRef.current?.value !== undefined) {
      searchInputRef.current.value = value;
    }
    debouceSearchInput(value);
  };

  return (
    <CommandNew
      style={{ width: "100%" }}
      shouldFilter={false}
      label="Global Command Menu"
      checked={checked}
      check={(key) =>
        setChecked((old) => {
          const newSet = new Set(old);
          newSet.add(key);
          return newSet;
        })
      }
      uncheck={(key) =>
        setChecked((old) => {
          const newSet = new Set(old);
          newSet.delete(key);
          return newSet;
        })
      }
    >
      <Stack>
        <Input
          ref={searchInputRef}
          size="lg"
          type="text"
          placeholder="What do you want to access?"
          onValueChange={debouceSearchInput}
          autoFocus={true}
          as={CommandNew.Input}
        />
        <Entitlements
          checked={checked}
          onSetSearch={onSetSearch}
          onShowSelected={onShowSelected}
        />
        <StyledCommandList
          mt={2}
          border="1px solid"
          rounded="md"
          borderColor="neutrals.300"
          p={1}
          pt={2}
        >
          <ReactWindowList
            style={{}}
            height={TARGETS * TARGET_HEIGHT}
            itemCount={onDisplay.length}
            itemSize={TARGET_HEIGHT}
            width="100%"
          >
            {({ index, style }) => {
              const target = onDisplay[index];
              return (
                <TargetDetail
                  showIcon
                  key={target.id}
                  as={CommandNew.Item}
                  h={TARGET_HEIGHT}
                  target={target}
                  style={style}
                  _selected={{
                    bg: "neutrals.100",
                  }}
                  // this value is used by the command palette
                  // ts-ignored because the typing doesn't propagate perfectly with the 'as' property
                  // @ts-ignore
                  value={target.id}
                  isChecked={checked.has(target.id.toLowerCase())}
                />
              );
            }}
          </ReactWindowList>
        </StyledCommandList>
        <Flex w="100%" mt={4}>
          <Button
            isDisabled={checked.size == 0 || submitLoading || targetsValidating}
            ml="auto"
            onClick={handleSubmit}
            isLoading={submitLoading || targetsValidating}
            loadingText={
              targetsValidating ? "Refreshing targets" : "Processing request..."
            }
          >
            Next ({actionKey[0]}+Enter)
          </Button>
        </Flex>
      </Stack>
    </CommandNew>
  );
};
interface EntitlementsProps {
  checked: Set<string>;
  onShowSelected: () => void;
  onSetSearch: (value: string) => void;
}
const Entitlements: React.FC<EntitlementsProps> = ({
  checked,
  onSetSearch,
  onShowSelected,
}) => {
  const { data: entitlements } = useUserListEntitlements();
  const { targets } = useTargets();
  return (
    <HStack mt={2} overflowX="auto">
      <FilterBlock
        label="All Resources"
        total={targets.length}
        onClick={() => {
          onSetSearch("");
        }}
      />
      <FilterBlock
        label="Selected"
        selected={checked.size}
        onClick={onShowSelected}
      />
      {entitlements?.entitlements.map((kind) => {
        const key = (
          kind.publisher +
          "#" +
          kind.name +
          "#" +
          kind.kind +
          "#"
        ).toLowerCase();
        return (
          <FilterBlock
            key={key}
            label={kind.kind}
            icon={kind.icon as ShortTypes}
            onClick={() => {
              onSetSearch(key);
            }}
            selected={[...checked].filter((id) => id.startsWith(key)).length}
          />
        );
      })}
    </HStack>
  );
};
interface FilterBlockProps extends CenterProps {
  icon?: ShortTypes;
  total?: number;
  selected?: number;
  label: string;
}
const FilterBlock: React.FC<FilterBlockProps> = ({
  label,
  total,
  selected,
  icon,
  ...rest
}) => {
  return (
    <Center
      rounded="md"
      h="84px"
      borderColor="neutrals.300"
      bg="white"
      borderWidth="1px"
      px={2}
      flexDirection="column"
      as={"button"}
      {...rest}
    >
      {icon !== undefined ? (
        <ProviderIcon shortType={icon} />
      ) : (
        <Box boxSize="22px" />
      )}
      <Text textStyle="Body/Small" noOfLines={1} textAlign="center">
        {label}
      </Text>
      {total === undefined ? (
        selected === undefined ? (
          <Box boxSize="22px" />
        ) : (
          <Text
            textStyle="Body/Small"
            noOfLines={1}
            textAlign="center"
            color="neutrals.500"
          >
            {`${selected} selected`}
          </Text>
        )
      ) : (
        <Text
          textStyle="Body/Small"
          noOfLines={1}
          textAlign="center"
          color="neutrals.500"
        >
          {`${total} total`}
        </Text>
      )}
    </Center>
  );
};
